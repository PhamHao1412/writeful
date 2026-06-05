package service

import (
	"content-service/internal/entity"
	"content-service/internal/model"
	"content-service/internal/repository"
	"context"
	"net/http"

	"gorm.io/gorm"
)

type IMediaService interface {
	UploadMedia(ctx context.Context, media []entity.Media) (error, int)
	GetMedia(ctx context.Context, id string) (*model.MediaResponse, error)
	DeleteMedia(ctx context.Context, id string) error
}

type MediaService struct {
	baseRepo  repository.IBaseRepository
	mediaRepo repository.IMediaRepository
}

func NewMediaService(baseRepo repository.IBaseRepository, mediaRepo repository.IMediaRepository) *MediaService {
	return &MediaService{baseRepo, mediaRepo}
}

func (s *MediaService) UploadMedia(ctx context.Context, medias []entity.Media) (error, int) {
	if len(medias) == 0 {
		return nil, http.StatusBadRequest
	}

	err := s.mediaRepo.Transaction(ctx, func(tx *gorm.DB) error {
		for _, media := range medias {
			// Check if media already exists in content-service to avoid duplicate key errors
			var exists entity.Media
			if err := tx.Where("id = ?", media.ID).First(&exists).Error; err == nil {
				// Media already exists (due to deduplication), safe to skip
				continue
			} else if err != gorm.ErrRecordNotFound {
				return err
			}

			var base entity.BaseEntity
			base.ID = media.ID
			base.CreatedAt = media.CreatedAt
			base.UpdatedAt = media.UpdatedAt
			if err := s.baseRepo.CreateTx(tx, &entity.Media{
				BaseEntity: base,
				URL:        media.URL,
				Type:       media.Type,
				Provider:   media.Provider,
				PublicID:   media.PublicID,
				Format:     media.Format,
				MimeType:   media.MimeType,
				FileSize:   media.FileSize,
				Width:      media.Width,
				Height:     media.Height,
				UploadedBy: media.UploadedBy,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusCreated

}

func (s *MediaService) GetMedia(ctx context.Context, id string) (*model.MediaResponse, error) {
	media, err := s.mediaRepo.First(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.MediaResponse{
		ID:         media.ID,
		URL:        media.URL,
		Type:       media.Type,
		Provider:   media.Provider,
		PublicID:   media.PublicID,
		MimeType:   media.MimeType,
		FileSize:   media.FileSize,
		Width:      media.Width,
		Height:     media.Height,
		UploadedBy: media.UploadedBy,
		CreatedAt:  media.CreatedAt,
	}, nil
}

func (s *MediaService) DeleteMedia(ctx context.Context, id string) error {
	return s.mediaRepo.Delete(ctx, id)
}
