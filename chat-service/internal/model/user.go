package model

import (
	"github.com/google/uuid"
	"time"
)

type GetUserRequest struct {
	ID     string
	Status string
	IDs    []string
}

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Status       string     `json:"status"`
	DisplayName  string     `json:"display_name"`
	AvtarURL     string     `json:"avatar_url"`
	Roles        []Role     `json:"roles"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
}

type Role struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	AvatarURL string    `json:"avatar_url"`
	Roles     []string  `json:"roles"`
}
