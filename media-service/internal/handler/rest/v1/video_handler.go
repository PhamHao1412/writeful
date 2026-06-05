package v1

import (
	"errors"
	"media_service/internal/handler/rest/dto"
	"media_service/internal/model"
	"media_service/internal/repository"
	"media_service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoSvc service.IVideoService
}

func NewVideoHandler(videoSvc service.IVideoService) *VideoHandler {
	return &VideoHandler{
		videoSvc: videoSvc,
	}
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("missing file")))
		return
	}
	ctx := c.Request.Context()
	meta, err := h.videoSvc.Upload(ctx, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(meta))
}

func (h *VideoHandler) GetVideo(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	video, err := h.videoSvc.GetMetadata(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ResponseNotFound(errors.New("video not found")))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(video))
}

func (h *VideoHandler) ListVideos(c *gin.Context) {
	var req model.GetListRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	videos, err := h.videoSvc.List(ctx, repository.GetVideoFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   req.Sort,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(videos).TotalItem(int64(len(videos))))
}

func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.videoSvc.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("video deleted successfully"))
}
