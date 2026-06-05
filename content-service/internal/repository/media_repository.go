package repository

import (
	"content-service/internal/entity"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	First(ctx context.Context, id string) (*entity.Media, error)
	FindByPostID(ctx context.Context, postID uuid.UUID) ([]entity.Media, error)
	Delete(ctx context.Context, id string) error
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(ctx context.Context, media *entity.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *MediaRepository) First(ctx context.Context, id string) (*entity.Media, error) {
	var media entity.Media
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&media).Error
	return &media, err
}

func (r *MediaRepository) FindByPostID(ctx context.Context, postID uuid.UUID) ([]entity.Media, error) {
	var medias []entity.Media
	err := r.db.WithContext(ctx).
		Joins("JOIN post_medias ON medias.id = post_medias.media_id").
		Where("post_medias.post_id = ?", postID).
		Order("post_medias.display_order").
		Find(&medias).Error
	return medias, err
}

func (r *MediaRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Media{}).Error
}
func (r *MediaRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
