package entity

import (
	"github.com/google/uuid"
)

type PostVersion struct {
	BaseEntity

	PostID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Version int       `gorm:"not null"`

	Content string `gorm:"type:text;not null"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
}

func (PostVersion) TableName() string {
	return SchemaName() + "post_versions"

}
