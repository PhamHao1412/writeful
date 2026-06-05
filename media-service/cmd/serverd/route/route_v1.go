package route

import (
	v1 "media_service/internal/handler/rest/v1"

	"github.com/gin-gonic/gin"
)

func V1Router(r *gin.Engine, healthHandler *v1.HealthHandler, imageHandler *v1.ImageHandler, videoHandler *v1.VideoHandler) {
	apiV1 := r.Group("/api/v1")
	{
		// Health check
		apiV1.GET("/health", healthHandler.Health)

		// Image endpoints
		imageRoute := apiV1.Group("/images")
		imageRoute.GET("", imageHandler.ListImages)
		imageRoute.GET("/:id", imageHandler.GetImage)
		imageRoute.POST("/upload", imageHandler.Upload)
		imageRoute.POST("/resize", imageHandler.Resize)
		imageRoute.POST("/convert", imageHandler.Convert)
		imageRoute.POST("/filter", imageHandler.Filter)
		imageRoute.POST("/crop", imageHandler.Crop)
		imageRoute.POST("/rotate", imageHandler.Rotate)
		imageRoute.POST("/flip", imageHandler.Flip)
		imageRoute.POST("/watermark", imageHandler.Watermark)
		imageRoute.POST("/compress", imageHandler.Compress)

		// Video endpoints
		videoRoute := apiV1.Group("/videos")
		videoRoute.POST("/upload", videoHandler.UploadVideo)
		videoRoute.GET("/:id", videoHandler.GetVideo)
		videoRoute.GET("", videoHandler.ListVideos)
		videoRoute.DELETE("/:id", videoHandler.DeleteVideo)
	}
}
