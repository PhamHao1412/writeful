package service

import (
	"context"
	"fmt"
	"media_service/internal/entity"
	"media_service/internal/model"
	"media_service/internal/pkg/cloudinary"
	"media_service/internal/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

const (
	MediaTypeVideo = "video"
	VideoFolder    = "image-service/originals/videos"
)

type IVideoService interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (*entity.Video, error)
	GetMetadata(ctx context.Context, id string) (*entity.Video, error)
	List(ctx context.Context, filter repository.GetVideoFilter) ([]entity.Video, error)
	Delete(ctx context.Context, id string) error
}

type VideoService struct {
	cldUtil cloudinary.ICloudinary
	repo    repository.IVideoRepository
}

func NewVideoService(cldUtil cloudinary.ICloudinary, repo repository.IVideoRepository) *VideoService {
	return &VideoService{
		cldUtil: cldUtil,
		repo:    repo,
	}
}

func (s *VideoService) Upload(ctx context.Context, file *multipart.FileHeader) (*entity.Video, error) {
	id := uuid.NewString()

	rs, err := s.cldUtil.UploadFile(file, model.UploadParams{
		Id:     id,
		Folder: VideoFolder,
		Type:   MediaTypeVideo,
	})
	if err != nil {
		return nil, err
	}

	video := &entity.Video{
		ID:     id,
		URL:    rs.URL,
		Format: rs.Format,
		Type:   MediaTypeVideo,
		Size:   file.Size,
	}

	if err := s.repo.Create(ctx, video); err != nil {
		return nil, err
	}
	return video, nil
}

func (s *VideoService) GetMetadata(ctx context.Context, id string) (*entity.Video, error) {
	video, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("video not found")
	}
	return video, nil
}

func (s *VideoService) List(ctx context.Context, filter repository.GetVideoFilter) ([]entity.Video, error) {
	return s.repo.List(ctx, filter)
}

func (s *VideoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
