package v1

import (
	"chat-service/internal/handler/rest"
	"chat-service/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking in production
		return true
	},
}

type WebSocketHandler struct {
	hub *service.Hub
}

func NewWebSocketHandler(hub *service.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleWebSocket godoc
// @Summary WebSocket connection
// @Description Establish WebSocket connection for real-time messaging
// @Tags websocket
// @Security BearerAuth
// @Router /ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Try to get userID from header first (for compatibility)
	userID := c.GetHeader(rest.HeaderUserID)

	// If not in header, try query params (browsers can't send custom headers in WebSocket)
	if userID == "" {
		userID = c.Query("user_id")
	}

	if userID == "" {
		log.Printf("WebSocket connection rejected: missing user_id")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user_id required"})
		return
	}

	// Optionally validate token from query params
	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
	}
	// TODO: Validate token if needed

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &service.Client{
		ID:     uuid.New().String(),
		UserID: userID,
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	client.Hub.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket connection established for user: %s", userID)
}
