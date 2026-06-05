package v1

import (
	"content-service/internal/handler/rest"
	"content-service/internal/handler/rest/dto"
	"content-service/internal/model"
	"content-service/internal/repository"
	"content-service/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type PostHandler struct {
	postSvc service.IPostService
}

func NewPostHandler(postSvc service.IPostService) *PostHandler {
	return &PostHandler{postSvc}
}

func (h *PostHandler) Create(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(err))
		return
	}

	userId, _ := uuid.Parse(userIdStr)

	post, err, status := h.postSvc.CreatePost(c.Request.Context(), userId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseOK(post).WithMessage("post created"))
}

func (h *PostHandler) Update(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid payload")))
		return
	}

	userID, _ := uuid.Parse(userIdStr)
	postID := c.Param("id")
	if err := h.postSvc.UpdatePost(c.Request.Context(), userID, postID, req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("post updated"))
}

func (h *PostHandler) Get(c *gin.Context) {
	postID := c.Param("id")

	post, err, status := h.postSvc.GetOne(c.Request.Context(), postID)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postSvc.GetPostBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) List(c *gin.Context) {
	//userIdStr := c.GetHeader(rest.HeaderUserID)
	//if userIdStr == "" {
	//	c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
	//	return
	//}

	var req model.GetPostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(err))
		return
	}

	filter := repository.GetPostFilter{
		Status: req.Status,
		Page:   req.Page,
		Size:   req.Size,
		Sort:   req.Sort,
		UserID: req.UserID,
	}

	result, total, err := h.postSvc.GetList(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(result).TotalItem(total))
}

func (h *PostHandler) Publish(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	postID := c.Param("id")

	if err, status := h.postSvc.Publish(c.Request.Context(), postID, userID); err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseOK(gin.H{}).WithMessage("post published"))

}

func (h *PostHandler) Unpublish(c *gin.Context) {
	userID := c.GetHeader(rest.HeaderUserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	postID := c.Param("id")
	if err, status := h.postSvc.Unpublish(c.Request.Context(), postID, userID); err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseOK(gin.H{}).WithMessage("post unpublished"))

}

func (h *PostHandler) Delete(c *gin.Context) {
	postID := c.Param("id")

	if err := h.postSvc.DeletePost(c.Request.Context(), postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

// Health godoc
// @Summary Health check
// @Description Check service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *PostHandler) Health(c *gin.Context) {
	cpuPercent, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "chat-service",
		"system": gin.H{
			"cpu_usage":    cpuPercent,
			"memory_usage": memInfo.UsedPercent,
		},
	})
}
