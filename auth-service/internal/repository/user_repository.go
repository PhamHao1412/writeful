package repository

import (
	"auth-service/internal/entity"
	"auth-service/pkg/util"
	"context"

	"gorm.io/gorm"
)

type GetUsersFilter struct {
	ID       string
	Name     string
	Username string
	Role     string
	Email    string
	Status   string
	IDs      []string
	Limit    int
	Offset   int
	Token    string
}

type IUserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	First(ctx context.Context, req GetUsersFilter) (*entity.User, error)
	Find(ctx context.Context, filter GetUsersFilter) ([]entity.User, int64, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error
	GetDB() *gorm.DB
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return u.DB.WithContext(ctx).Create(user).Error
}

func (u *UserRepository) Update(ctx context.Context, user *entity.User) error {
	return u.DB.WithContext(ctx).Save(user).Error

}

func (u *UserRepository) Delete(ctx context.Context, user *entity.User) error {
	return u.DB.WithContext(ctx).Delete(user).Error
}

func (u *UserRepository) First(ctx context.Context, filter GetUsersFilter) (*entity.User, error) {
	var user entity.User
	queryBuilder := u.DB.WithContext(ctx).Table(user.TableName()).Preload("Roles")

	if filter.ID != "" {
		queryBuilder = queryBuilder.Where("id = ?", filter.ID)
	}

	if filter.Name != "" {
		queryBuilder = queryBuilder.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Username != "" {
		queryBuilder = queryBuilder.Where("username = ?", filter.Username)
	}

	if filter.Email != "" {
		queryBuilder = queryBuilder.Where("email = ?", filter.Email)
	}

	if filter.Role != "" {
		queryBuilder = queryBuilder.Where("role = ?", filter.Role)
	}

	if filter.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", filter.Status)
	}

	if len(filter.IDs) > 0 {
		queryBuilder = queryBuilder.Where("id IN (?)", filter.IDs)
	}

	err := queryBuilder.First(&user).Error

	return &user, err
}

func (u *UserRepository) Find(ctx context.Context, filter GetUsersFilter) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	queryBuilder := u.DB.WithContext(ctx).Model(entity.User{}).Preload("Roles")

	// Apply filters
	if filter.ID != "" {
		queryBuilder = queryBuilder.Where("id = ?", filter.ID)
	}

	if filter.Name != "" {
		queryBuilder = queryBuilder.Where("display_name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Username != "" {
		queryBuilder = queryBuilder.Where("username LIKE ?", "%"+filter.Username+"%")
	}

	if filter.Email != "" {
		queryBuilder = queryBuilder.Where("email LIKE ?", "%"+filter.Email+"%")
	}

	if filter.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", filter.Status)
	}

	if len(filter.IDs) > 0 {
		queryBuilder = queryBuilder.Where("id IN ?", filter.IDs)
	}

	// Count total records
	countQuery := queryBuilder
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := util.ToOffsetLimit(filter.Offset, filter.Limit)
	if err := queryBuilder.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *UserRepository) GetDB() *gorm.DB {
	return u.DB
}

func (u *UserRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return u.DB.WithContext(ctx).Transaction(fn)
}
