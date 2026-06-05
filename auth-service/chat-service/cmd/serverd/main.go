package main

import (
	"chat-service/cmd/serverd/route"
	"chat-service/internal/app"
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
	config, err := env.ReadAppConfig[app.Config]()
	if err != nil {
		log.Fatal("failed to read app:", err)
	}

	application := app.NewApp(&config)

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
	port := config.Port
	appName := config.AppName

	log.Printf("%v running at :%s", appName, port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
