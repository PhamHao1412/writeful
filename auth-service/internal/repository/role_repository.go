package repository

import (
	"auth-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	First(ctx context.Context, code string) (*entity.Role, error)
	Create(ctx context.Context, userRole *entity.UserRole) error
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) First(ctx context.Context, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.DB.WithContext(ctx).Table(role.TableName()).Where("code = ?", code).First(&role).Error
	return &role, err
}

func (r *RoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	return r.DB.WithContext(ctx).Create(userRole).Error
}
