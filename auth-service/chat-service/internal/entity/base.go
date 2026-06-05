package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var schemaName = "chat_service"

func SetSchemaName(schema string) {
	schemaName = schema
}

func SchemaName() string {
	return schemaName + "."
}

type BaseEntity struct {
	ID        string         `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
