package entity

import (
	"fmt"
	"time"
)

type User struct {
	BaseEntity

	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`

	Status     string `gorm:"type:varchar(20);default:'active'" json:"status"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`

	DisplayName  string     `json:"display_name"`
	Bio          string     `json:"bio"`
	AvatarURL    string     `json:"avatar_url"`
	LastActiveAt *time.Time `gorm:"column:last_active_at" json:"last_active_at,omitempty"`

	Roles []Role `gorm:"many2many:auth_service.user_roles" json:"roles,omitempty"`
}

func (User) TableName() string {
	return fmt.Sprintf("%susers", SchemaName())

}
