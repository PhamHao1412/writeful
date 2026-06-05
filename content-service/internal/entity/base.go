package entity

import (
	"content-service/internal/app"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var config *app.Config

type BaseEntity struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
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
