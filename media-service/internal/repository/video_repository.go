package repository

import (
	"context"
	"media_service/internal/entity"
	"media_service/internal/pkg/util"

	"gorm.io/gorm"
)

type GetVideoFilter struct {
	Limit  int
	Offset int
	Sort   string
	Type   string
}

type IVideoRepository interface {
	Create(ctx context.Context, video *entity.Video) error
	FindByID(ctx context.Context, id string) (*entity.Video, error)
	Update(ctx context.Context, video *entity.Video) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter GetVideoFilter) ([]entity.Video, error)
}

type VideoRepository struct {
	DB *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{DB: db}
}

func (r *VideoRepository) Create(ctx context.Context, video *entity.Video) error {
	return r.DB.WithContext(ctx).Create(video).Error
}

func (r *VideoRepository) FindByID(ctx context.Context, id string) (*entity.Video, error) {
	var video entity.Video
	if err := r.DB.WithContext(ctx).First(&video, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *VideoRepository) Update(ctx context.Context, video *entity.Video) error {
	return r.DB.WithContext(ctx).Save(video).Error
}

func (r *VideoRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.Video{}, "id = ?", id).Error
}

func (r *VideoRepository) List(ctx context.Context, filter GetVideoFilter) ([]entity.Video, error) {
	var videos []entity.Video
	queryBuilder := r.DB.WithContext(ctx)
	if filter.Sort == "asc" {
		queryBuilder = queryBuilder.Order("created_at ASC")
	} else {
		queryBuilder = queryBuilder.Order("created_at DESC")
	}
	offset, limit := util.ToOffsetLimit(filter.Offset, filter.Limit)
	if err := queryBuilder.Limit(limit).Offset(offset).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
