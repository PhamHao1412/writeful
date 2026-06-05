package repository

import (
	"auth-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type IFollowerRepository interface {
	Create(ctx context.Context, follower *entity.Follower) error
	Delete(ctx context.Context, followerID, followingID string) error
	GetFollowers(ctx context.Context, userID string, limit, offset int) ([]entity.User, int64, error)
	GetFollowing(ctx context.Context, userID string, limit, offset int) ([]entity.User, int64, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
}

type FollowerRepository struct {
	DB *gorm.DB
}

func NewFollowerRepository(db *gorm.DB) *FollowerRepository {
	return &FollowerRepository{DB: db}
}

func (r *FollowerRepository) Create(ctx context.Context, follower *entity.Follower) error {
	return r.DB.WithContext(ctx).Create(follower).Error
}

func (r *FollowerRepository) Delete(ctx context.Context, followerID, followingID string) error {
	return r.DB.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&entity.Follower{}).Error
}

func (r *FollowerRepository) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	baseQuery := r.DB.WithContext(ctx).
		Table(entity.User{}.TableName()).
		Joins("JOIN "+entity.Follower{}.TableName()+" ON "+entity.User{}.TableName()+".id::text = "+entity.Follower{}.TableName()+".follower_id").
		Where(entity.Follower{}.TableName()+".following_id = ?", userID)

	if err := baseQuery.Session(&gorm.Session{}).Distinct(entity.User{}.TableName() + ".id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.Select("DISTINCT " + entity.User{}.TableName() + ".*").Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *FollowerRepository) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	baseQuery := r.DB.WithContext(ctx).
		Table(entity.User{}.TableName()).
		Joins("JOIN "+entity.Follower{}.TableName()+" ON "+entity.User{}.TableName()+".id::text = "+entity.Follower{}.TableName()+".following_id").
		Where(entity.Follower{}.TableName()+".follower_id = ?", userID)

	if err := baseQuery.Session(&gorm.Session{}).Distinct(entity.User{}.TableName() + ".id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.Select("DISTINCT " + entity.User{}.TableName() + ".*").Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *FollowerRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&entity.Follower{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
