package service

import (
	"chat-service/internal/gateway/auth"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	ID     string
	UserID string
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients mapped by user ID
	clients map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan *BroadcastMessage

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	authClient auth.Client
	convRepo   repository.IConversationRepository

	mu sync.RWMutex
}

type BroadcastMessage struct {
	UserIDs []string
	Message []byte
}

func NewHub(authClient auth.Client, convRepo repository.IConversationRepository) *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
		authClient: authClient,
		convRepo:   convRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; !ok {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			isFirstConn := len(h.clients[client.UserID]) == 1
			h.mu.Unlock()

			log.Printf("Client registered: UserID=%s, ClientID=%s, Total users=%d",
				client.UserID, client.ID, len(h.clients))

			// Send currently online users to the newly connected user
			onlineUsers := h.GetOnlineUsers()
			wsMessage := model.WSMessage{
				Type:    "online_users",
				Payload: onlineUsers,
			}
			if data, err := json.Marshal(wsMessage); err == nil {
				select {
				case client.Send <- data:
				default:
				}
			}

			// Broadcast online status to others
			if isFirstConn {
				go h.handleUserOnline(client.UserID)
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			_, isStillOnline := h.clients[client.UserID]
			h.mu.Unlock()

			log.Printf("Client unregistered: UserID=%s, ClientID=%s, Total users=%d",
				client.UserID, client.ID, len(h.clients))

			// Broadcast offline status to others
			if !isStillOnline {
				go h.handleUserOffline(client.UserID)
			}

		case message := <-h.broadcast:
			h.mu.RLock()
			log.Printf("[WebSocket Hub] Processing broadcast queue. Target UserIDs: %v", message.UserIDs)
			for _, userID := range message.UserIDs {
				if clients, ok := h.clients[userID]; ok {
					log.Printf("[WebSocket Hub] Found %d active client connection(s) for UserID=%s. Forwarding...", len(clients), userID)
					for client := range clients {
						select {
						case client.Send <- message.Message:
							log.Printf("[WebSocket Hub] Message successfully queued to send queue for UserID=%s, ClientID=%s", userID, client.ID)
						default:
							close(client.Send)
							delete(clients, client)
							if len(clients) == 0 {
								delete(h.clients, userID)
							}
							log.Printf("[WebSocket Hub WARNING] client Send channel blocked/full for UserID=%s. Disconnected client.", userID)
						}
					}
				} else {
					log.Printf("[WebSocket Hub WARNING] UserID=%s is NOT online (no active client connections). Message dropped.", userID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToUsers(userIDs []string, messageType model.WSMessageType, payload interface{}) {
	wsMessage := model.WSMessage{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	h.broadcast <- &BroadcastMessage{
		UserIDs: userIDs,
		Message: data,
	}
}

func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

func (h *Hub) getRelatedUserIDs(userID string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conversations, _, err := h.convRepo.GetByUserID(ctx, userID, 1, 500)
	if err != nil {
		log.Printf("[WebSocket Hub Error] Failed to get conversations for user %s: %v", userID, err)
		return nil
	}

	userMap := make(map[string]bool)
	for _, conv := range conversations {
		for _, part := range conv.Participants {
			if part.UserID != userID && part.LeftAt == nil {
				userMap[part.UserID] = true
			}
		}
	}

	userIDs := make([]string, 0, len(userMap))
	for id := range userMap {
		userIDs = append(userIDs, id)
	}
	return userIDs
}

func (h *Hub) handleUserOnline(userID string) {
	// 1. Update active status timestamp in DB
	err, status := h.authClient.UpdateActiveStatus(userID, time.Now())
	if err != nil {
		log.Printf("[WebSocket Hub Error] Failed to update active status for user %s: %v (status=%d)", userID, err, status)
	}

	// 2. Find conversation partners
	relatedIDs := h.getRelatedUserIDs(userID)
	if len(relatedIDs) > 0 {
		// 3. Broadcast user_online event
		h.BroadcastToUsers(relatedIDs, model.WSMessageTypeUserOnline, map[string]interface{}{
			"user_id": userID,
			"online":  true,
		})
	}
}

func (h *Hub) handleUserOffline(userID string) {
	now := time.Now()
	// 1. Update active status timestamp in DB
	err, status := h.authClient.UpdateActiveStatus(userID, now)
	if err != nil {
		log.Printf("[WebSocket Hub Error] Failed to update active status for user %s: %v (status=%d)", userID, err, status)
	}

	// 2. Find conversation partners
	relatedIDs := h.getRelatedUserIDs(userID)
	if len(relatedIDs) > 0 {
		// 3. Broadcast user_offline event with last active timestamp
		h.BroadcastToUsers(relatedIDs, model.WSMessageTypeUserOffline, map[string]interface{}{
			"user_id":        userID,
			"online":         false,
			"last_active_at": now.Format(time.RFC3339),
		})
	}
}

func (h *Hub) handleTyping(userID string, convID string, payload map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	participants, err := h.convRepo.GetParticipants(ctx, convID)
	if err != nil {
		log.Printf("[WebSocket Hub Error] Failed to get participants for typing in conv %s: %v", convID, err)
		return
	}

	targetUserIDs := make([]string, 0, len(participants))
	for _, part := range participants {
		if part.UserID != userID && part.LeftAt == nil {
			targetUserIDs = append(targetUserIDs, part.UserID)
		}
	}

	if len(targetUserIDs) > 0 {
		h.BroadcastToUsers(targetUserIDs, model.WSMessageTypeTyping, payload)
	}
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Log raw received message to diagnose routing or serialization issues
		log.Printf("[WebSocket Raw] Received frame from UserID=%s: %s", c.UserID, string(message))

		// Handle incoming messages (e.g., typing indicators)
		var wsMsg model.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Process message based on type
		switch wsMsg.Type {
		case model.WSMessageTypeTyping:
			log.Printf("[WebSocket] Received typing indicator from UserID=%s", c.UserID)
			// Broadcast typing indicator to other participants
			if payload, ok := wsMsg.Payload.(map[string]interface{}); ok {
				if convID, ok := payload["conversation_id"].(string); ok {
					go c.Hub.handleTyping(c.UserID, convID, payload)
				}
			}
		case model.WSMessageTypeCallInitiate,
			model.WSMessageTypeCallReceive,
			model.WSMessageTypeCallRinging,
			model.WSMessageTypeCallReject,
			model.WSMessageTypeCallCancel,
			model.WSMessageTypeCallHangup,
			model.WSMessageTypeWebRTCOffer,
			model.WSMessageTypeWebRTCAnswer,
			model.WSMessageTypeWebRTCICE:
			log.Printf("[Signaling] Received calling packet type='%s' from caller UserID=%s", wsMsg.Type, c.UserID)
			if payload, ok := wsMsg.Payload.(map[string]interface{}); ok {
				targetUserID, _ := payload["target_user_id"].(string)
				if targetUserID != "" {
					log.Printf("[Signaling] Relaying packet type='%s' from UserID=%s to target UserID=%s", wsMsg.Type, c.UserID, targetUserID)
					c.Hub.BroadcastToUsers([]string{targetUserID}, wsMsg.Type, payload)
				} else {
					log.Printf("[Signaling WARNING] Signaling message type='%s' missing target_user_id in payload: %+v", wsMsg.Type, payload)
				}
			} else {
				log.Printf("[Signaling ERROR] Failed to cast signaling message payload to map[string]interface{}. Type was: %T", wsMsg.Payload)
			}
		default:
			log.Printf("[WebSocket WARNING] Received unhandled or unrecognized message type='%s' from UserID=%s", wsMsg.Type, c.UserID)
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		// Add queued messages to the current websocket message
		n := len(c.Send)
		for i := 0; i < n; i++ {
			w.Write([]byte{'\n'})
			w.Write(<-c.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}
