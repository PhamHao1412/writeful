package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//func InitDB() *gorm.DB {
//	host := lit.GetEnv("DB_HOST", "localhost")
//	port := lit.GetEnv("DB_PORT", "5432")
//	user := lit.GetEnv("DB_USER", "postgres")
//	password := lit.GetEnv("DB_PASSWORD", "postgres")
//	dbname := lit.GetEnv("DB_NAME", "chat_db")
//	schema := lit.GetEnv("DB_SCHEMA", "public")
//	sslmode := lit.GetEnv("DB_SSL_MODE", "disable")
//
//	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
//		host, port, user, password, dbname, sslmode)
//
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//
//	if err != nil {
//		log.Fatalf("Failed to connect to database: %v", err)
//	}
//
//	// Set schema
//	entity.SetSchemaName(schema)
//
//	// Auto migrate
//	if err := db.AutoMigrate(
//		&entity.Conversation{},
//		&entity.Participant{},
//		&entity.Message{},
//	); err != nil {
//		log.Fatalf("Failed to migrate database: %v", err)
//	}
//
//	log.Println("Database connected and migrated successfully")
//	return db
//}

func Connect(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
