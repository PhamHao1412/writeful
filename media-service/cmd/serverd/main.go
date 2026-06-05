package main

import (
	"log"
	"media_service/cmd/serverd/route"
	"media_service/internal/app"
	"media_service/internal/db"
	"media_service/internal/entity"
	v1 "media_service/internal/handler/rest/v1"
	"media_service/internal/middleware"
	"media_service/internal/pkg/cloudinary"
	"media_service/internal/repository"
	"media_service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/viebiz/lit/env"
)

func main() {
	gin.ForceConsoleColor()
	r := gin.Default()
	//r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	config, err := env.ReadAppConfig[app.Config]()
	if err != nil {
		log.Fatal("failed to read config:", err)
	}
	entity.SetConfig(&config)

	// DB
	database, err := db.Connect(config.PG.URL)
	if err != nil {
		log.Fatal("failed to connect DB:", err)
	}

	// Cloudinary
	cloudStore, err := cloudinary.NewCloudinaryUtil(config.Cloudinary)
	if err != nil {
		log.Fatal("cloudinary init error:", err)
	}

	// Repositories
	imageRepo := repository.NewImageRepository(database)
	videoRepo := repository.NewVideoRepository(database)

	// Services
	imgService := service.NewImageService(cloudStore, imageRepo)
	videoService := service.NewVideoService(cloudStore, videoRepo)

	//handler
	imageHandler := v1.NewImageHandler(imgService)
	videoHandler := v1.NewVideoHandler(videoService)
	healthHandler := v1.NewHealthHandler(database)

	route.V1Router(r, healthHandler, imageHandler, videoHandler)
	port := config.Port
	log.Printf("%s service running at :%s", config.AppName, port)
	r.Run(":" + port)
}
