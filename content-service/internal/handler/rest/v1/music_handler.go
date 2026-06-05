package v1

import (
	"content-service/internal/handler/rest"
	"content-service/internal/handler/rest/dto"
	"content-service/internal/model"
	"content-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MusicHandler struct {
	musicSvc service.IMusicService
}

func NewMusicHandler(musicSvc service.IMusicService) *MusicHandler {
	return &MusicHandler{musicSvc: musicSvc}
}

func (h *MusicHandler) Create(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	var uploaderID *uuid.UUID
	if userIdStr != "" {
		uid, err := uuid.Parse(userIdStr)
		if err == nil {
			uploaderID = &uid
		}
	}

	var req model.AddMusicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(err))
		return
	}

	music, err := h.musicSvc.AddMusic(c.Request.Context(), uploaderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseOK(music).WithMessage("music track added to server library"))
}

func (h *MusicHandler) List(c *gin.Context) {
	genre := c.Query("genre")
	search := c.Query("search")

	list, err := h.musicSvc.ListMusic(c.Request.Context(), genre, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(list))
}
