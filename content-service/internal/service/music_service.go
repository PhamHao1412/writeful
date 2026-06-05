package service

import (
	"content-service/internal/entity"
	"content-service/internal/model"
	"content-service/internal/repository"
	"context"

	"github.com/google/uuid"
)

type IMusicService interface {
	AddMusic(ctx context.Context, uploaderID *uuid.UUID, req model.AddMusicRequest) (*entity.Music, error)
	ListMusic(ctx context.Context, genre string, search string) ([]entity.Music, error)
}

type MusicService struct {
	repo repository.IMusicRepository
}

func NewMusicService(repo repository.IMusicRepository) *MusicService {
	return &MusicService{repo: repo}
}

func (s *MusicService) AddMusic(ctx context.Context, uploaderID *uuid.UUID, req model.AddMusicRequest) (*entity.Music, error) {
	music := &entity.Music{
		ID:         uuid.New(),
		Title:      req.Title,
		Artist:     req.Artist,
		URL:        req.URL,
		CoverURL:   req.CoverURL,
		Genre:      req.Genre,
		UploadedBy: uploaderID,
	}
	if err := s.repo.Create(ctx, music); err != nil {
		return nil, err
	}
	return music, nil
}

func (s *MusicService) ListMusic(ctx context.Context, genre string, search string) ([]entity.Music, error) {
	return s.repo.List(ctx, genre, search)
}
