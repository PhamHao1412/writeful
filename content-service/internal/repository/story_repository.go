package repository

import (
	"content-service/internal/entity"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IStoryRepository interface {
	Create(ctx context.Context, story *entity.Story) error
	GetActiveStories(ctx context.Context) ([]entity.Story, error)
	MarkAsSeen(ctx context.Context, storyID, viewerID uuid.UUID) error
	GetSeenStoryIDs(ctx context.Context, viewerID uuid.UUID) (map[uuid.UUID]bool, error)
	Delete(ctx context.Context, storyID, userID uuid.UUID) error
}

type StoryRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) *StoryRepository {
	return &StoryRepository{db: db}
}

func (r *StoryRepository) Create(ctx context.Context, story *entity.Story) error {
	return r.db.WithContext(ctx).Create(story).Error
}

func (r *StoryRepository) Delete(ctx context.Context, storyID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", storyID, userID).
		Delete(&entity.Story{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *StoryRepository) GetActiveStories(ctx context.Context) ([]entity.Story, error) {
	var list []entity.Story
	err := r.db.WithContext(ctx).
		Where("expires_at > ? AND status = 'active'", time.Now()).
		Order("created_at asc").
		Find(&list).Error
	return list, err
}

func (r *StoryRepository) MarkAsSeen(ctx context.Context, storyID, viewerID uuid.UUID) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Story{}).Where("id = ?", storyID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	view := &entity.StoryView{
		StoryID:  storyID,
		ViewerID: viewerID,
	}
	// Use clause to ignore duplicate key violation if user views story multiple times
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(view).Error
}

func (r *StoryRepository) GetSeenStoryIDs(ctx context.Context, viewerID uuid.UUID) (map[uuid.UUID]bool, error) {
	var views []entity.StoryView
	err := r.db.WithContext(ctx).Where("viewer_id = ?", viewerID).Find(&views).Error
	if err != nil {
		return nil, err
	}
	m := make(map[uuid.UUID]bool)
	for _, v := range views {
		m[v.StoryID] = true
	}
	return m, nil
}
