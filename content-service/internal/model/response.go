package model

import (
	"time"

	"github.com/google/uuid"
)

type TResponse[T any] struct {
	Data T `json:"data"`
}

type PostResponse struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	Title         string         `json:"title"`
	Subtitle      string         `json:"subtitle"`
	Slug          string         `json:"slug"`
	Excerpt       string         `json:"excerpt"`
	Content       string         `json:"content"`
	Status        string         `json:"status"`
	PublishedAt   *time.Time     `json:"published_at"`
	CoverImageURL string         `json:"cover_image_url"`
	CoverMedia    *MediaResponse `json:"cover_media,omitempty"`
	User          *User          `json:"user"`
	Media         []Media        `json:"media,omitempty"`
	Tags          []TagResponse  `json:"tags,omitempty"`
	Version       int            `json:"version"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Media struct {
	PostID       string `json:"post_id"`
	MediaID      string `json:"media_id"`
	URL          string `json:"url"`
	DisplayOrder int    `json:"display_order"`
}

type MediaResponse struct {
	ID         uuid.UUID `json:"id"`
	URL        string    `json:"url"`
	Type       string    `json:"type"`
	Provider   string    `json:"provider"`
	PublicID   *string   `json:"public_id"`
	MimeType   *string   `json:"mime_type"`
	FileSize   *int64    `json:"file_size"`
	Width      *int      `json:"width"`
	Height     *int      `json:"height"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type ResponseList struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}
