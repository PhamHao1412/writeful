package app

import (
	"chat-service/internal/db"
	v1 "chat-service/internal/handler/rest/v1"
	"chat-service/internal/repository"
	"chat-service/internal/service"
	"log"

	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB

	// Repositories
	ConversationRepo repository.IConversationRepository
	MessageRepo      repository.IMessageRepository

	// Services
	ConversationSvc service.IConversationService
	MessageSvc      service.IMessageService
	Hub             *service.Hub

	// Handlers
	HealthHandler       *v1.HealthHandler
	ConversationHandler *v1.ConversationHandler
	MessageHandler      *v1.MessageHandler
	WebSocketHandler    *v1.WebSocketHandler
}

func NewApp(config *Config) *App {
	app := &App{}

	// Initialize database
	database, err := db.Connect(config.PG.URL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	app.DB = database
	// Initialize repositories
	app.ConversationRepo = repository.NewConversationRepository(app.DB)
	app.MessageRepo = repository.NewMessageRepository(app.DB)

	// Initialize WebSocket hub
	app.Hub = service.NewHub()
	go app.Hub.Run()
	log.Println("WebSocket hub started")

	// Initialize services
	app.ConversationSvc = service.NewConversationService(app.ConversationRepo, app.MessageRepo)
	app.MessageSvc = service.NewMessageService(app.MessageRepo, app.ConversationRepo)

	// Initialize handlers
	app.HealthHandler = v1.NewHealthHandler()
	app.ConversationHandler = v1.NewConversationHandler(app.ConversationSvc)
	app.MessageHandler = v1.NewMessageHandler(app.MessageSvc, app.ConversationRepo, app.Hub)
	app.WebSocketHandler = v1.NewWebSocketHandler(app.Hub)

	return app
}
