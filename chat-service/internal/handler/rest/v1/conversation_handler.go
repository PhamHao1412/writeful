package v1

import (
	"chat-service/internal/handler/rest"
	"chat-service/internal/handler/rest/v1/dto"
	"chat-service/internal/model"
	"chat-service/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	conversationSvc service.IConversationService
}

func NewConversationHandler(conversationSvc service.IConversationService) *ConversationHandler {
	return &ConversationHandler{conversationSvc: conversationSvc}
}

// CreateConversation godoc
// @Summary Create a new conversation
// @Description Create a new direct or group conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateConversationRequest true "Create conversation request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /conversations [post]
func (h *ConversationHandler) Create(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	response, err, status := h.conversationSvc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(status, dto.ResponseOK(response).WithMessage("conversation created successfully"))
}

// GetConversations godoc
// @Summary Get user's conversations
// @Description Get list of conversations for the authenticated user
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /conversations [get]
func (h *ConversationHandler) GetList(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.GetConversationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid query parameters")))
		return
	}

	response, total, err := h.conversationSvc.GetList(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(response).TotalItem(total))
}

// GetConversation godoc
// @Summary Get conversation details
// @Description Get details of a specific conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /conversations/{id} [get]
func (h *ConversationHandler) GetByID(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("conversation ID is required")))
		return
	}

	response, err, status := h.conversationSvc.GetByID(c.Request.Context(), userID, conversationID)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(response))
}

// AddParticipants godoc
// @Summary Add participants to conversation
// @Description Add new participants to a group conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Param request body model.AddParticipantsRequest true "Add participants request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /conversations/{id}/participants [post]
func (h *ConversationHandler) AddParticipants(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("conversation ID is required")))
		return
	}

	var req model.AddParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	err := h.conversationSvc.AddParticipants(c.Request.Context(), userID, conversationID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("participants added successfully"))
}

// RemoveParticipant godoc
// @Summary Remove participant from conversation
// @Description Remove a participant from a conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversation ID"
// @Param participant_id path string true "Participant ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /conversations/{id}/participants/{participant_id} [delete]
func (h *ConversationHandler) RemoveParticipant(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	conversationID := c.Param("id")
	participantID := c.Param("participant_id")

	if conversationID == "" || participantID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("conversation ID and participant ID are required")))
		return
	}

	err := h.conversationSvc.RemoveParticipant(c.Request.Context(), userID, conversationID, participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("participant removed successfully"))
}

// MarkAsRead godoc
// @Summary Mark conversation as read
// @Description Mark all messages in a conversation as read
// @Tags conversations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.MarkAsReadRequest true "Mark as read request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /conversations/read [post]
func (h *ConversationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	err := h.conversationSvc.MarkAsRead(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("marked as read"))
}
