package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"media_service/internal/handler/rest/dto"
	"media_service/internal/model"
	"media_service/internal/pkg/util"
	"media_service/internal/repository"
	"media_service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imgSvc service.IImageService
}

func NewImageHandler(imgSvc service.IImageService) *ImageHandler {
	return &ImageHandler{
		imgSvc: imgSvc,
	}
}

func (h *ImageHandler) Upload(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	slog.Info("user id:", userID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequestWithTitle(err, "invalid form data"))
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("missing file")))
		return
	}

	for _, file := range files {
		if !util.ValidateImageType(file.Filename, file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(
				fmt.Errorf("invalid image type for file: %s", file.Filename)))
			return
		}
	}

	ctx := c.Request.Context()
	images, err := h.imgSvc.Upload(ctx, userID, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(images))
}

func (h *ImageHandler) Resize(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Resize(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) Convert(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Convert(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) Filter(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Filter(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) GetImage(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	meta, ok := h.imgSvc.GetMetadata(ctx, id)
	if !ok {
		c.JSON(http.StatusNotFound, dto.ResponseNotFound(errors.New("image not found")))
		return
	}
	c.JSON(http.StatusOK, meta)
}

func (h *ImageHandler) ListImages(c *gin.Context) {
	var req model.GetListRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	images, err := h.imgSvc.List(ctx, repository.GetImageFilter{
		Limit:      req.Limit,
		Offset:     req.Offset,
		Sort:       req.Sort,
		Type:       req.Type,
		Provider:   req.Provider,
		UploadedBy: req.UploadedBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(images).TotalItem(int64(len(images))))
}

func (h *ImageHandler) Crop(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Crop(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) Rotate(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Rotate(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) Flip(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Flip(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) Watermark(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Watermark(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}

func (h *ImageHandler) Compress(c *gin.Context) {
	var req model.TransformRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request")))
		return
	}
	ctx := c.Request.Context()
	url, err := h.imgSvc.Compress(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseInternalServerError(err))
		return
	}
	c.JSON(http.StatusOK, dto.ResponseOK(url))
}
