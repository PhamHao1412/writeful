package entity

import (
	"fmt"

	"github.com/google/uuid"
)

type UserRole struct {
	BaseEntity

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (UserRole) TableName() string {
	return fmt.Sprintf("%suser_roles", SchemaName())

}
