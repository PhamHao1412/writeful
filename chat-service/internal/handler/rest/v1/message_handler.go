package v1

import (
	"chat-service/internal/handler/rest"
	"chat-service/internal/handler/rest/v1/dto"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"chat-service/internal/service"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageSvc       service.IMessageService
	conversationRepo repository.IConversationRepository
	hub              *service.Hub
}

func NewMessageHandler(
	messageSvc service.IMessageService,
	conversationRepo repository.IConversationRepository,
	hub *service.Hub,
) *MessageHandler {
	return &MessageHandler{
		messageSvc:       messageSvc,
		conversationRepo: conversationRepo,
		hub:              hub,
	}
}

// SendMessage godoc
// @Summary Send a message
// @Description Send a new message to a conversation
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.SendMessageRequest true "Send message request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /messages [post]
func (h *MessageHandler) Send(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	response, err, status := h.messageSvc.Send(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	// Broadcast message to conversation participants via WebSocket
	// Use context.Background() to avoid cancellation when request finishes
	go h.broadcastMessage(context.Background(), req.ConversationID, response)

	c.JSON(status, dto.ResponseOK(response).WithMessage("message sent successfully"))
}

// GetMessages godoc
// @Summary Get messages
// @Description Get messages from a conversation
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation_id query string true "Conversation ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param before query string false "Message ID for cursor pagination"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.GetMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid query parameters")))
		return
	}

	response, total, err, status := h.messageSvc.GetMessages(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(response).TotalItem(total))
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Delete a message (only sender can delete)
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /messages/{id} [delete]
func (h *MessageHandler) Delete(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("message ID is required")))
		return
	}

	err := h.messageSvc.Delete(c.Request.Context(), userID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("message deleted successfully"))
}

func (h *MessageHandler) broadcastMessage(ctx context.Context, conversationID string, message *model.MessageResponse) {
	// Get all participants
	participants, err := h.conversationRepo.GetParticipants(ctx, conversationID)
	if err != nil {
		return
	}

	// Extract user IDs - send to ALL participants including sender
	// This ensures sender also gets real-time update via WebSocket
	userIDs := make([]string, 0, len(participants))
	for _, p := range participants {
		userIDs = append(userIDs, p.UserID)
	}

	// Broadcast via WebSocket to all participants
	h.hub.BroadcastToUsers(userIDs, model.WSMessageTypeNewMessage, message)
}

// SendSignaling godoc
// @Summary Send WebRTC signaling packet
// @Description Send WebRTC signaling packet to a specific user via HTTP POST, to be relayed via WebSocket
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.WSMessage true "Signaling packet request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /messages/signaling [post]
func (h *MessageHandler) SendSignaling(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var wsMsg model.WSMessage
	if err := c.ShouldBindJSON(&wsMsg); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	payloadMap, ok := wsMsg.Payload.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("payload must be a json object")))
		return
	}

	targetUserID, _ := payloadMap["target_user_id"].(string)
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("target_user_id is required in payload")))
		return
	}

	// Relay signaling message down to target recipient via active WS connection
	h.hub.BroadcastToUsers([]string{targetUserID}, wsMsg.Type, wsMsg.Payload)

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("signaling message relayed successfully"))
}

// ToggleReaction godoc
// @Summary Toggle message reaction
// @Description Toggle a reaction on a message
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Param request body model.ToggleReactionRequest true "Toggle reaction request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /messages/{id}/reactions [post]
func (h *MessageHandler) ToggleReaction(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	messageID := c.Param("id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("message ID is required")))
		return
	}

	var req model.ToggleReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	response, action, err, status := h.messageSvc.ToggleReaction(c.Request.Context(), userID, messageID, req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	// Broadcast to conversation participants via WebSocket in background
	go h.broadcastReaction(context.Background(), messageID, userID, req.Emoji, action)

	c.JSON(status, dto.ResponseOK(response).WithMessage("reaction toggled successfully"))
}

func (h *MessageHandler) broadcastReaction(ctx context.Context, messageID, userID, emoji, action string) {
	msg, err := h.messageSvc.GetByID(ctx, messageID)
	if err != nil {
		return
	}

	// Get all participants
	participants, err := h.conversationRepo.GetParticipants(ctx, msg.ConversationID)
	if err != nil {
		return
	}

	userIDs := make([]string, 0, len(participants))
	for _, p := range participants {
		userIDs = append(userIDs, p.UserID)
	}

	payload := model.WSReactionPayload{
		MessageID:      messageID,
		ConversationID: msg.ConversationID,
		UserID:         userID,
		Emoji:          emoji,
		Action:         action,
	}

	h.hub.BroadcastToUsers(userIDs, model.WSMessageTypeReaction, payload)
}
