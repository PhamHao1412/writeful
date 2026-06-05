package main

import (
	"content-service/internal/db"
	"content-service/internal/entity"
	"log"
)

func main() {
	dsn := "postgres://postgres:hacker1412@127.0.0.1:5438/playground?sslmode=disable&search_path=content_service"
	log.Printf("Connecting directly to %s...", dsn)
	dbConn, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("Database connected. Starting GORM AutoMigrate...")
	err = dbConn.AutoMigrate(
		&entity.Music{},
		&entity.Story{},
		&entity.StoryView{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	log.Println("GORM AutoMigrate completed successfully.")
}
