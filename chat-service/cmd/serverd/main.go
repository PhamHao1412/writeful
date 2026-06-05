package main

import (
	"chat-service/cmd/serverd/route"
	"chat-service/internal/app"
	"chat-service/internal/config"
	"chat-service/pkg/logger"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/viebiz/lit/env"
)

func main() {
	// Load environment variables
	//if err := lit.LoadEnv(".env"); err != nil {
	//	log.Printf("Warning: .env file not found: %v", err)
	//}

	// Initialize application

	// Set Gin mode
	cfg, err := env.ReadAppConfig[config.Config]()
	if err != nil {
		log.Fatal("failed to read app:", err)
	}
	logger.Init(&cfg)

	application := app.NewApp(&cfg)

	r := gin.Default()

	// Setup routes
	route.V1Router(
		r,
		application.HealthHandler,
		application.ConversationHandler,
		application.MessageHandler,
		application.WebSocketHandler,
	)

	// Start server
	port := cfg.Port

	log.Printf("🚀🚀🚀 CHAT SERVICE v2.0-DIAGNOSTICS running at :%s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
