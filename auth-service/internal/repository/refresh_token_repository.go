package repository

import (
	"auth-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type GetTokensFilter struct {
	ID    string
	Token string
}

type IRefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *entity.RefreshToken) error
	First(ctx context.Context, filter GetTokensFilter) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAll(ctx context.Context, userId string) error
}

type RefreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (u *RefreshTokenRepository) Create(ctx context.Context, refreshToken *entity.RefreshToken) error {
	return u.DB.WithContext(ctx).Create(refreshToken).Error
}

func (u *RefreshTokenRepository) First(ctx context.Context, filter GetTokensFilter) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	queryBuilder := u.DB.WithContext(ctx).Table(refreshToken.TableName())

	if filter.ID != "" {
		queryBuilder = queryBuilder.Where("id = ?", filter.ID)
	}

	if filter.Token != "" {
		queryBuilder = queryBuilder.Where("token = ?", filter.Token)
	}

	err := queryBuilder.First(&refreshToken).Error

	return &refreshToken, err
}

func (u *RefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	return u.DB.WithContext(ctx).Model(&entity.RefreshToken{}).Where("id = ?", id).Update("is_revoked", true).Error
}

func (u *RefreshTokenRepository) RevokeAll(ctx context.Context, userId string) error {
	return u.DB.WithContext(ctx).Model(&entity.RefreshToken{}).Where("user_id = ?", userId).Update("is_revoked", true).Error
}
