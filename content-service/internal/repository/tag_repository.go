package repository

import (
	"content-service/internal/entity"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	First(ctx context.Context, id string) (*entity.Tag, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	FindByPostID(ctx context.Context, postID uuid.UUID) ([]entity.Tag, error)
	Find(ctx context.Context) ([]entity.Tag, error)
	Delete(ctx context.Context, id string) error
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *TagRepository) First(ctx context.Context, id string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	return &tag, err
}

func (r *TagRepository) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	return &tag, err
}

func (r *TagRepository) FindByPostID(ctx context.Context, postID uuid.UUID) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN post_tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Find(ctx context.Context) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Tag{}).Error
}
