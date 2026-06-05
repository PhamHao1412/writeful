package repository

import (
	"context"
	"media_service/internal/entity"
	"media_service/internal/pkg/util"

	"gorm.io/gorm"
)

type GetImageFilter struct {
	Limit      int
	Offset     int
	Id         string
	Type       string
	Provider   string
	UploadedBy string
	Sort       string
}

type IImageRepository interface {
	Create(ctx context.Context, img *entity.Image) error
	FindByID(ctx context.Context, id string) (*entity.Image, error)
	FindByHash(ctx context.Context, hash string) (*entity.Image, error)
	Update(ctx context.Context, img *entity.Image) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter GetImageFilter) ([]entity.Image, error)
}

type ImageRepository struct {
	DB *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{DB: db}
}

func (r *ImageRepository) FindByHash(ctx context.Context, hash string) (*entity.Image, error) {
	var img entity.Image
	if err := r.DB.WithContext(ctx).Where("file_hash = ?", hash).First(&img).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepository) Create(ctx context.Context, img *entity.Image) error {
	return r.DB.WithContext(ctx).Create(img).Error
}

func (r *ImageRepository) FindByID(ctx context.Context, id string) (*entity.Image, error) {
	var img entity.Image
	if err := r.DB.WithContext(ctx).First(&img, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepository) Update(ctx context.Context, img *entity.Image) error {
	return r.DB.WithContext(ctx).Save(img).Error
}

func (r *ImageRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Model(&entity.Image{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *ImageRepository) List(ctx context.Context, filter GetImageFilter) ([]entity.Image, error) {
	var images []entity.Image
	queryBuilder := r.DB.WithContext(ctx)

	if filter.Type != "" {
		queryBuilder = queryBuilder.Where("type = ?", filter.Type)
	}
	if filter.Provider != "" {
		queryBuilder = queryBuilder.Where("provider = ?", filter.Provider)
	}
	if filter.UploadedBy != "" {
		queryBuilder = queryBuilder.Where("uploaded_by = ?", filter.UploadedBy)
	}

	if filter.Sort == "asc" {
		queryBuilder = queryBuilder.Order("created_at ASC")
	} else {
		queryBuilder = queryBuilder.Order("created_at DESC")
	}

	offset, limit := util.ToOffsetLimit(filter.Offset, filter.Limit)
	if err := queryBuilder.Limit(limit).Offset(offset).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
