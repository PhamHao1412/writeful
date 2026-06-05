package v1

import (
	"content-service/internal/entity"
	"content-service/internal/handler/rest/dto"
	"content-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	mediaSvc service.IMediaService
}

func NewMediaHandler(mediaSvc service.IMediaService) *MediaHandler {
	return &MediaHandler{mediaSvc}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	var media []entity.Media
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, status := h.mediaSvc.UploadMedia(c.Request.Context(), media)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusCreated, media)
}

func (h *MediaHandler) Get(c *gin.Context) {
	id := c.Param("id")

	media, err := h.mediaSvc.GetMedia(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *MediaHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.mediaSvc.DeleteMedia(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "media deleted"})
}
