package entity

import (
	"fmt"

	"github.com/google/uuid"
)

type RefreshToken struct {
	BaseEntity

	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"not null;uniqueIndex" json:"-"`
	ExpiredAt int64     `gorm:"not null" json:"expired_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
}

func (RefreshToken) TableName() string {
	return fmt.Sprintf("%srefresh_tokens", SchemaName())
}
