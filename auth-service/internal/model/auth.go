package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	AvatarURL string    `json:"avatar_url"`
	Roles     []string  `json:"roles"`
}

type AuthTokenSignedResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"exp"`
}

type TokenClaims struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type Token struct {
	ID     string   `json:"id"`
	Email  string   `json:"email"`
	Status string   `json:"status"`
	Roles  []string `json:"roles"`
}

type JWTClaims struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Status    string   `json:"status"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"` // access | refresh
	jwt.RegisteredClaims
}

type UserResponse struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	Status       string     `json:"status"`
	AvatarURL    string     `json:"avatar_url"`
	DisplayName  string     `json:"display_name"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	Roles        []Role     `json:"roles"`
}

type Role struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
