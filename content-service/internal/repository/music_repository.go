package repository

import (
	"content-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type IMusicRepository interface {
	Create(ctx context.Context, music *entity.Music) error
	List(ctx context.Context, genre string, search string) ([]entity.Music, error)
}

type MusicRepository struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) *MusicRepository {
	return &MusicRepository{db: db}
}

func (r *MusicRepository) Create(ctx context.Context, music *entity.Music) error {
	return r.db.WithContext(ctx).Create(music).Error
}

func (r *MusicRepository) List(ctx context.Context, genre string, search string) ([]entity.Music, error) {
	var list []entity.Music
	query := r.db.WithContext(ctx)
	if genre != "" && genre != "all" {
		query = query.Where("genre = ?", genre)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR artist ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	err := query.Order("created_at desc").Find(&list).Error
	return list, err
}
