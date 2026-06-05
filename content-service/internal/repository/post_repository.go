package repository

import (
	"content-service/internal/entity"
	"content-service/pkg/util"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetPostFilter struct {
	Slug      string
	ExactSlug string
	Status    string
	ID        string
	Sort      string
	UserID    string
	TagID     string
	Page      int
	Size      int
}

type IPostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	Update(ctx context.Context, post *entity.Post) error
	UpdateStatus(ctx context.Context, id, userID string, status string, publishedAt *time.Time) error
	First(ctx context.Context, filter GetPostFilter) (*entity.Post, error)
	Find(ctx context.Context, filter GetPostFilter) ([]entity.Post, int64, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter GetPostFilter) (int64, error)
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	AddTags(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTags(ctx context.Context, postID uuid.UUID) error
	AddMedias(ctx context.Context, postID uuid.UUID, mediaIDs []uuid.UUID) error
	RemoveMedias(ctx context.Context, postID uuid.UUID) error
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *PostRepository) UpdateStatus(ctx context.Context, id, userID string, status string, publishedAt *time.Time) error {
	if err := r.db.WithContext(ctx).Model(&entity.Post{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"status":       status,
			"published_at": publishedAt,
			"updated_at":   time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) First(ctx context.Context, filter GetPostFilter) (*entity.Post, error) {
	var post entity.Post
	queryBuilder := r.db.WithContext(ctx).Model(&entity.Post{}).Preload("PostMedia.Media")
	if filter.ID != "" {
		queryBuilder = queryBuilder.Where("id = ?", filter.ID)
	}

	if filter.Slug != "" {
		queryBuilder = queryBuilder.Where("slug = ?", filter.Slug)
	}

	if filter.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", filter.Status)
	}

	if filter.UserID != "" {
		queryBuilder = queryBuilder.Where("user_id = ?", filter.UserID)
	}

	err := queryBuilder.First(&post).Error
	return &post, err
}

func (r *PostRepository) Find(ctx context.Context, filter GetPostFilter) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64

	queryBuilder := r.db.WithContext(ctx).Model(&entity.Post{}).Preload("PostMedia.Media")

	if filter.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", filter.Status)
	}

	if filter.UserID != "" {
		queryBuilder = queryBuilder.Where("user_id = ?", filter.UserID)
	}

	if filter.Sort != "" {
		queryBuilder = queryBuilder.Order(filter.Sort)
	}

	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := util.ToOffsetLimit(filter.Page, filter.Size)

	if err := queryBuilder.Offset(offset).Limit(limit).Find(&posts).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepository) Count(ctx context.Context, filter GetPostFilter) (int64, error) {
	queryBuilder := r.db.WithContext(ctx).Model(&entity.Post{}).Unscoped()

	if filter.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", filter.Status)
	}

	if filter.UserID != "" {
		queryBuilder = queryBuilder.Where("user_id = ?", filter.UserID)
	}

	if filter.TagID != "" {
		queryBuilder = queryBuilder.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ?", filter.TagID)
	}

	if filter.ExactSlug != "" {
		queryBuilder = queryBuilder.Where("slug = ?", filter.ExactSlug)
	} else if filter.Slug != "" {
		queryBuilder = queryBuilder.Where("slug LIKE ?", filter.Slug+"%")
	}
	var total int64
	if err := queryBuilder.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Post{}).Error
}

func (r *PostRepository) AddTags(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, tagID := range tagIDs {
			postTag := &entity.PostTag{
				PostID: postID,
				TagID:  tagID,
			}
			if err := tx.Create(postTag).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostRepository) RemoveTags(ctx context.Context, postID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&entity.PostTag{}).Error
}

func (r *PostRepository) AddMedias(ctx context.Context, postID uuid.UUID, mediaIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, mediaID := range mediaIDs {
			postMedia := &entity.PostMedia{
				PostID:       postID,
				MediaID:      mediaID,
				DisplayOrder: i,
			}
			if err := tx.Create(postMedia).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostRepository) RemoveMedias(ctx context.Context, postID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&entity.PostMedia{}).Error
}

func (r *PostRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
