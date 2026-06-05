package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	BaseEntity

	UserId uuid.UUID `gorm:"type:uuid;not null;index"`

	Title      string `gorm:"title;not null"`
	Subtitle   string `gorm:"subtitle;not null"`
	Slug       string `gorm:"uniqueIndex;not null"`
	Visibility string `gorm:"type:varchar(20);default:'public'"`
	Excerpt    string

	Status        string      `gorm:"type:varchar(20);default:'draft'"`
	PublishedAt   *time.Time  `gorm:"published_at"`
	CoverImageURL string      `gorm:"cover_image_url""`
	Media         []Media     `gorm:"many2many:post_media" json:"media,omitempty"`
	PostMedia     []PostMedia ` json:"post_media,omitempty"`
}

func (Post) TableName() string {
	return SchemaName() + "posts"
}
func (Post) TableNameAlias(alias string) string {
	return SchemaName() + "posts " + alias
}
