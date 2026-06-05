package entity

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name string    `gorm:"uniqueIndex;not null"`
	Slug string    `gorm:"uniqueIndex;not null"`
}

func (Tag) TableName() string {
	return SchemaName() + "tags"
}
