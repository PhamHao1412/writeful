package entity

import (
	"fmt"
	"media_service/internal/app"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var config *app.Config

type BaseEntity struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt *time.Time     `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

// SetConfig sets the application config for entity package
func SetConfig(cfg *app.Config) {
	config = cfg
}

func SchemaName() string {
	if config == nil || config.DBSchemaName == "" {
		return ""
	}
	return fmt.Sprintf("%s.", config.DBSchemaName)
}
