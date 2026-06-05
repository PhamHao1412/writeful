package model

import "time"

type CreateStoryRequest struct {
	MediaURL    string  `json:"media_url" binding:"required"`
	Caption     string  `json:"caption"`
	AudioURL    *string `json:"audio_url"`
	AudioTitle  *string `json:"audio_title"`
	AudioArtist *string `json:"audio_artist"`
	AudioOffset *int    `json:"audio_offset"`
}

type AddMusicRequest struct {
	Title    string `json:"title" binding:"required"`
	Artist   string `json:"artist" binding:"required"`
	URL      string `json:"url" binding:"required"`
	CoverURL string `json:"cover_url"`
	Genre    string `json:"genre" binding:"required"`
}

type StoryViewDTO struct {
	StoryID  string `json:"story_id"`
	ViewerID string `json:"viewer_id"`
}

type UserStoriesGroup struct {
	UserID    string            `json:"user_id"`
	Username  string            `json:"username"`
	AvatarURL string            `json:"avatar_url"`
	HasUnread bool              `json:"has_unread"`
	Stories   []StoryDisplayDTO `json:"stories"`
}

type StoryDisplayDTO struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	MediaURL    string    `json:"media_url"`
	Caption     string    `json:"caption,omitempty"`
	AudioURL    *string   `json:"audio_url,omitempty"`
	AudioTitle  *string   `json:"audio_title,omitempty"`
	AudioArtist *string   `json:"audio_artist,omitempty"`
	AudioOffset *int      `json:"audio_offset,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Seen        bool      `json:"seen"`
}
