package entity

import "github.com/google/uuid"

type PostTag struct {
	PostID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TagID  uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (PostTag) TableName() string {
	return SchemaName() + "post_tags"
}
