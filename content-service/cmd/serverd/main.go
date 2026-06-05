package main

import (
	"content-service/cmd/serverd/route"
	"content-service/internal/app"
	"content-service/internal/db"
	"content-service/internal/entity"
	"content-service/internal/gateway/auth"
	v1 "content-service/internal/handler/rest/v1"
	"content-service/internal/repository"
	"content-service/internal/service"
	"content-service/pkg/logger"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/viebiz/lit/env"
)

func main() {
	gin.ForceConsoleColor()
	r := gin.Default()
	//r.Use(middleware.CORS())
	//r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	cfg, _ := env.ReadAppConfig[app.Config]()
	entity.SetConfig(&cfg)
	logger.Init(&cfg)
	dbConn, _ := db.Connect(cfg.PG.URL)
	if err := db.AutoMigrateAndSeed(dbConn); err != nil {
		log.Fatalf("failed to auto migrate and seed database: %v", err)
	}

	// gateway
	authClient := auth.NewClient(&cfg)

	// repositories
	baseRepo := repository.NewBaseRepository(dbConn)
	postRepo := repository.NewPostRepository(dbConn)
	versionRepo := repository.NewPostVersionRepository(dbConn)
	mediaRepo := repository.NewMediaRepository(dbConn)
	tagRepo := repository.NewTagRepository(dbConn)
	musicRepo := repository.NewMusicRepository(dbConn)
	storyRepo := repository.NewStoryRepository(dbConn)

	// services
	postSvc := service.NewPostService(baseRepo, postRepo, versionRepo, mediaRepo, tagRepo, authClient)
	mediaSvc := service.NewMediaService(baseRepo, mediaRepo)
	tagSvc := service.NewTagService(tagRepo)
	musicSvc := service.NewMusicService(musicRepo)
	storySvc := service.NewStoryService(storyRepo, authClient)

	// Handlers
	postHandler := v1.NewPostHandler(postSvc)
	mediaHandler := v1.NewMediaHandler(mediaSvc)
	tagHandler := v1.NewTagHandler(tagSvc)
	healthHandler := v1.NewHealthHandler()
	musicHandler := v1.NewMusicHandler(musicSvc)
	storyHandler := v1.NewStoryHandler(storySvc)

	route.V1Router(r, postHandler, mediaHandler, tagHandler, healthHandler, storyHandler, musicHandler)
	log.Printf("%s service running at :%s", cfg.AppName, cfg.Port)

	r.Run(":" + cfg.Port)
}
