package entity

import (
	"auth-service/internal/app"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var config *app.Config

type BaseEntity struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func SetConfig(cfg *app.Config) {
	config = cfg
}

func SchemaName() string {
	if config == nil || config.DBSchemaName == "" {
		return ""
	}
	return fmt.Sprintf("%s.", config.DBSchemaName)
}
