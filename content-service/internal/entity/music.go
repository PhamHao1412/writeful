package entity

import (
	"time"

	"github.com/google/uuid"
)

type Music struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title      string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Artist     string     `gorm:"column:artist;type:varchar(255);not null" json:"artist"`
	URL        string     `gorm:"column:url;type:text;not null" json:"url"`
	CoverURL   string     `gorm:"column:cover_url;type:text" json:"cover_url,omitempty"`
	Genre      string     `gorm:"column:genre;type:varchar(50);not null;default:'vpop'" json:"genre"`
	UploadedBy *uuid.UUID `gorm:"column:uploaded_by;type:uuid" json:"uploaded_by,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Music) TableName() string {
	return SchemaName() + "musics"
}
