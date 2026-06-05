package model

import "github.com/google/uuid"

type CreatePostRequest struct {
	Content       string      `json:"content"`
	Title         string      `json:"title" binding:"required"`
	Subtitle      string      `json:"subtitle"`
	Slug          string      `json:"slug"`
	Visibility    string      `json:"visibility"`
	Excerpt       string      `json:"excerpt"`
	CoverImageURL string      `json:"cover_image_url"`
	MediaIDs      []uuid.UUID `json:"media_ids"`
	TagIDs        []uuid.UUID `json:"tag_ids"`
}

type UpdatePostRequest struct {
	Content       string      `json:"content"`
	Title         string      `json:"title"`
	Subtitle      string      `json:"subtitle"`
	Excerpt       string      `json:"excerpt"`
	CoverImageURL string      `json:"cover_image_url"`
	MediaIDs      []uuid.UUID `json:"media_ids"`
	TagIDs        []uuid.UUID `json:"tag_ids"`
}

type UploadMediaRequest struct {
	File     string `form:"file" binding:"required"`
	Type     string `form:"type" binding:"required,oneof=image video"`
	MimeType string `form:"mime_type"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type GetPostRequest struct {
	Page   int    `form:"page"`
	Size   int    `form:"size"`
	Status string `form:"status"`
	UserID string `form:"user_id"`
	Sort   string `form:"sort"`
}
