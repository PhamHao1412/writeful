package entity

import (
	"time"

	"github.com/google/uuid"
)

type Story struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	Type        string    `gorm:"column:type;type:varchar(50);not null;default:'image'" json:"type"`
	MediaURL    string    `gorm:"column:media_url;type:text;not null" json:"media_url"`
	Caption     string    `gorm:"column:caption;type:text" json:"caption,omitempty"`
	AudioURL    *string   `gorm:"column:audio_url;type:text" json:"audio_url,omitempty"`
	AudioTitle  *string   `gorm:"column:audio_title;type:varchar(255)" json:"audio_title,omitempty"`
	AudioArtist *string   `gorm:"column:audio_artist;type:varchar(255)" json:"audio_artist,omitempty"`
	AudioOffset *int      `gorm:"column:audio_offset;type:integer;default:0" json:"audio_offset,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt   time.Time `gorm:"column:expires_at;type:timestamp;not null" json:"expires_at"`
	Status      string    `gorm:"column:status;type:varchar(50);default:'active'" json:"status"`

	// Relations
	Views []StoryView `gorm:"foreignKey:StoryID" json:"views,omitempty"`
}

func (Story) TableName() string {
	return SchemaName() + "stories"
}

type StoryView struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoryID   uuid.UUID `gorm:"column:story_id;type:uuid;not null" json:"story_id"`
	ViewerID  uuid.UUID `gorm:"column:viewer_id;type:uuid;not null" json:"viewer_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (StoryView) TableName() string {
	return SchemaName() + "story_views"
}
