package route

import (
	v1 "content-service/internal/handler/rest/v1"

	"github.com/gin-gonic/gin"
)

func V1Router(
	r *gin.Engine,
	postHandler *v1.PostHandler,
	mediaHandler *v1.MediaHandler,
	tagHandler *v1.TagHandler,
	healthHandler *v1.HealthHandler,
	storyHandler *v1.StoryHandler,
	musicHandler *v1.MusicHandler,
) {
	apiV1 := r.Group("/api/v1")

	// Public routes
	apiV1.GET("/health", healthHandler.Health)

	apiV1.GET("/post", postHandler.List)
	apiV1.GET("/post/:id", postHandler.Get)
	apiV1.GET("/tag", tagHandler.List)
	apiV1.GET("/tag/:id", tagHandler.Get)

	// Stories & Musics
	apiV1.GET("/stories", storyHandler.List)
	apiV1.POST("/stories", storyHandler.Create)
	apiV1.POST("/stories/:id/seen", storyHandler.MarkSeen)
	apiV1.DELETE("/stories/:id", storyHandler.Delete)
	apiV1.GET("/musics", musicHandler.List)
	apiV1.POST("/musics", musicHandler.Create)

	// Posts
	content := apiV1.Group("/post")
	content.POST("", postHandler.Create)
	content.PUT("/:id", postHandler.Update)
	content.POST("/:id/publish", postHandler.Publish)
	content.POST("/:id/unpublish", postHandler.Unpublish)
	content.DELETE("/:id", postHandler.Delete)
	media := apiV1.Group("/media")

	// Media
	media.POST("", mediaHandler.Upload)
	media.GET("/:id", mediaHandler.Get)
	media.DELETE("/:id", mediaHandler.Delete)
	tag := apiV1.Group("/tag")

	// Tags
	tag.POST("", tagHandler.Create)
	tag.DELETE("/:id", tagHandler.Delete)
}
