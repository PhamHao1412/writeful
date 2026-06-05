package entity

import "github.com/google/uuid"

type PostMedia struct {
	BaseEntity

	PostID       uuid.UUID `json:"post_id" gorm:"type:uuid;index"`
	MediaID      uuid.UUID `json:"media_id" gorm:"type:uuid;index"`
	DisplayOrder int       `json:"display_display_order"`

	// Relationships
	Media Media `json:"media,omitempty" gorm:"foreignKey:MediaID"`
}

func (PostMedia) TableName() string {
	return SchemaName() + "post_media"
}
func (PostMedia) TableNameAlias(alias string) string {
	return SchemaName() + "post_media " + alias
}
