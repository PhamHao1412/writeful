package v1

import (
	"content-service/internal/model"
	"content-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagSvc service.ITagService
}

func NewTagHandler(tagSvc service.ITagService) *TagHandler {
	return &TagHandler{tagSvc}
}

func (h *TagHandler) Create(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tagSvc.CreateTag(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func (h *TagHandler) Get(c *gin.Context) {
	id := c.Param("id")

	tag, err := h.tagSvc.GetTag(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tagSvc.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.tagSvc.DeleteTag(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}
