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
	go h.broadcastMessage(c.Request.Context(), req.ConversationID, response)

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

	// Extract user IDs
	userIDs := make([]string, 0, len(participants))
	for _, p := range participants {
		// Don't send to the sender
		if p.UserID != message.SenderID {
			userIDs = append(userIDs, p.UserID)
		}
	}

	// Broadcast via WebSocket
	h.hub.BroadcastToUsers(userIDs, model.WSMessageTypeNewMessage, message)
}
