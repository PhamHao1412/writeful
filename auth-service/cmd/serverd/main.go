package main

import (
	"auth-service/cmd/serverd/route"
	"auth-service/internal/app"
	"auth-service/internal/db"
	"auth-service/internal/entity"
	apiv1 "auth-service/internal/handler/rest/v1"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/viebiz/lit/env"
)

func main() {
	gin.ForceConsoleColor()
	r := gin.Default()
	config, err := env.ReadAppConfig[app.Config]()
	if err != nil {
		log.Fatal("failed to read app:", err)
	}
	entity.SetConfig(&config)

	// DB
	database, err := db.Connect(config.PG.URL)
	if err != nil {
		log.Fatal("failed to connect DB:", err)
	}

	userRepo := repository.NewUserRepository(database)
	roleRepo := repository.NewRoleRepository(database)
	refreshTokenRepo := repository.NewRefreshTokenRepository(database)
	followRepo := repository.NewFollowerRepository(database)

	// Services
	jwkSvc := service.NewJWKService(&config)
	userSvc := service.NewUserService(&config, userRepo, roleRepo, refreshTokenRepo, jwkSvc)
	followSvc := service.NewFollowerService(followRepo, userRepo)

	// Handlers
	healthHandler := apiv1.NewHealthHandler(database)
	authHandler := apiv1.NewAuthHandler(userSvc)
	followHandler := apiv1.NewFollowerHandler(followSvc)

	// Middleware
	//r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	route.V1Router(r, healthHandler, authHandler, followHandler)
	port := config.Port
	appName := config.AppName
	log.Printf("%v running at :%s", appName, port)
	log.Println("🔐 Auth endpoints: /api/v1/auth/signup, /api/v1/auth/login, /api/v1/auth/refresh")
	log.Println("📝 Post endpoints: /api/v1/posts (protected with JWT)")
	r.Run(":" + port)
}
