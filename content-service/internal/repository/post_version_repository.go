package repository

import (
	"content-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type GetPostVersionFilter struct {
	PostID string
	Slug   string
	Status string
	Sort   string
}

type IPostVersionRepository interface {
	First(ctx context.Context, filter GetPostVersionFilter) (*entity.PostVersion, error)
	Create(ctx context.Context, v *entity.PostVersion) error
}

type PostVersionRepository struct {
	db *gorm.DB
}

func NewPostVersionRepository(db *gorm.DB) *PostVersionRepository {
	return &PostVersionRepository{db: db}
}

func (r *PostVersionRepository) First(ctx context.Context, filter GetPostVersionFilter) (*entity.PostVersion, error) {
	var v entity.PostVersion
	queryBuilder := r.db.WithContext(ctx).Model(&entity.PostVersion{})

	if filter.PostID != "" {
		queryBuilder = queryBuilder.Where("post_id = ?", filter.PostID)
	}

	if filter.Slug != "" {
		queryBuilder = queryBuilder.Where("slug = ?", filter.Slug)
	}

	if filter.Sort != "" {
		queryBuilder = queryBuilder.Order(filter.Sort)
	}

	err := queryBuilder.First(&v).Error
	return &v, err
}

func (r *PostVersionRepository) Create(ctx context.Context, v *entity.PostVersion) error {
	return r.db.WithContext(ctx).Create(v).Error
}
