package route

import (
	v1 "chat-service/internal/handler/rest/v1"
	"chat-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func V1Router(
	r *gin.Engine,
	healthHandler *v1.HealthHandler,
	conversationHandler *v1.ConversationHandler,
	messageHandler *v1.MessageHandler,
	wsHandler *v1.WebSocketHandler,
) {
	// CORS middleware
	//r.Use(middleware.CORSMiddleware())

	// Root health check
	r.GET("/health", healthHandler.Health)

	// WebSocket endpoint (requires auth)
	r.GET("/ws", middleware.AuthMiddleware(), wsHandler.HandleWebSocket)

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Apply auth middleware to all API routes
		api.Use(middleware.AuthMiddleware())

		// Conversation routes
		conversations := api.Group("/conversations")
		{
			conversations.POST("", conversationHandler.Create)
			conversations.GET("", conversationHandler.GetList)
			conversations.GET("/:id", conversationHandler.GetByID)
			conversations.POST("/:id/participants", conversationHandler.AddParticipants)
			conversations.DELETE("/:id/participants/:participant_id", conversationHandler.RemoveParticipant)
			conversations.POST("/read", conversationHandler.MarkAsRead)
		}

		// Message routes
		messages := api.Group("/messages")
		{
			messages.POST("", messageHandler.Send)
			messages.GET("", messageHandler.GetMessages)
			messages.DELETE("/:id", messageHandler.Delete)
		}
	}
}
