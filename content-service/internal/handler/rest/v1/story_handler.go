package v1

import (
	"content-service/internal/handler/rest"
	"content-service/internal/handler/rest/dto"
	"content-service/internal/model"
	"content-service/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoryHandler struct {
	storySvc service.IStoryService
}

func NewStoryHandler(storySvc service.IStoryService) *StoryHandler {
	return &StoryHandler{storySvc: storySvc}
}

func (h *StoryHandler) Create(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}

	var req model.CreateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(err))
		return
	}

	story, err := h.storySvc.CreateStory(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseOK(story).WithMessage("story created successfully"))
}

func (h *StoryHandler) List(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}

	groups, err := h.storySvc.GetFeedStories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(groups))
}

func (h *StoryHandler) MarkSeen(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}

	storyIdStr := c.Param("id")
	storyID, err := uuid.Parse(storyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid story id")))
		return
	}

	if err := h.storySvc.MarkStoryAsSeen(c.Request.Context(), storyID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ResponseError(errors.New("story not found"), http.StatusNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("story marked as seen"))
}

func (h *StoryHandler) Delete(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}

	storyIdStr := c.Param("id")
	storyID, err := uuid.Parse(storyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid story id")))
		return
	}

	if err := h.storySvc.DeleteStory(c.Request.Context(), storyID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ResponseError(errors.New("story not found or unauthorized"), http.StatusNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("story deleted successfully"))
}
