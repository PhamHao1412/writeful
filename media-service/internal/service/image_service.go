package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"media_service/internal/entity"
	"media_service/internal/model"
	"media_service/internal/pkg/cloudinary"
	"media_service/internal/repository"
	"mime/multipart"

	"gorm.io/gorm"
)

const (
	MediaTypeImage = "image"
	ImageFolder    = "image-service/originals/images"
)

type IImageService interface {
	Upload(ctx context.Context, uploadedBy string, files []*multipart.FileHeader) ([]*entity.Image, error)
	GetMetadata(ctx context.Context, id string) (*entity.Image, bool)
	Resize(ctx context.Context, req model.TransformRequest) (string, error)
	Convert(ctx context.Context, req model.TransformRequest) (string, error)
	Filter(ctx context.Context, req model.TransformRequest) (string, error)
	Crop(ctx context.Context, req model.TransformRequest) (string, error)
	Rotate(ctx context.Context, req model.TransformRequest) (string, error)
	Flip(ctx context.Context, req model.TransformRequest) (string, error)
	Watermark(ctx context.Context, req model.TransformRequest) (string, error)
	Compress(ctx context.Context, req model.TransformRequest) (string, error)
	List(ctx context.Context, filter repository.GetImageFilter) ([]entity.Image, error)
}

type ImageService struct {
	cldUtil cloudinary.ICloudinary
	repo    repository.IImageRepository
}

func NewImageService(cldUtil cloudinary.ICloudinary, repo repository.IImageRepository) *ImageService {
	return &ImageService{cldUtil: cldUtil, repo: repo}
}

func (s *ImageService) Upload(ctx context.Context, uploadedBy string, files []*multipart.FileHeader) ([]*entity.Image, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	images := make([]*entity.Image, 0, len(files))

	for _, file := range files {
		// 1. Calculate file hash using SHA-256
		f, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s for hashing: %w", file.Filename, err)
		}

		hash := sha256.New()
		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to calculate hash for %s: %w", file.Filename, err)
		}
		f.Close()

		hashString := fmt.Sprintf("%x", hash.Sum(nil))

		// 2. Check if this hash already exists in DB
		existingImg, err := s.repo.FindByHash(ctx, hashString)
		if err == nil && existingImg != nil {
			slog.Info("Deduplication triggered! Reusing existing image upload",
				"filename", file.Filename,
				"hash", hashString,
				"existingID", existingImg.ID,
				"url", existingImg.URL,
			)
			images = append(images, existingImg)
			continue
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("database lookup error for hash: %w", err)
		}

		// 3. Not a duplicate: upload file to Cloudinary
		rs, err := s.cldUtil.UploadFile(file, model.UploadParams{
			Folder: ImageFolder,
			Type:   MediaTypeImage,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", file.Filename, err)
		}

		mimeType := file.Header.Get("Content-Type")

		img := &entity.Image{
			URL:        rs.URL,
			Type:       MediaTypeImage,
			Format:     &rs.Format,
			Provider:   "cloudinary",
			PublicID:   &rs.PublicId,
			MimeType:   &mimeType,
			FileSize:   &file.Size,
			Width:      &rs.Width,
			Height:     &rs.Height,
			UploadedBy: uploadedBy,
			FileHash:   &hashString,
		}

		if err := s.repo.Create(ctx, img); err != nil {
			return nil, fmt.Errorf("failed to save metadata for %s: %w", file.Filename, err)
		}

		images = append(images, img)
	}

	return images, nil
}

func (s *ImageService) GetMetadata(ctx context.Context, id string) (*entity.Image, bool) {
	img, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, false
	}
	return img, true
}

func (s *ImageService) List(ctx context.Context, filter repository.GetImageFilter) ([]entity.Image, error) {
	return s.repo.List(ctx, filter)
}

func (s *ImageService) Resize(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}

	url, err := s.cldUtil.Resize(img.URL, req.Width, req.Height)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *ImageService) Convert(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}

	url, err := s.cldUtil.Convert(img.URL, req.Format)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *ImageService) Filter(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}

	url, err := s.cldUtil.ApplyFilter(img.URL, req.Filter)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *ImageService) Crop(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}
	url, err := s.cldUtil.Crop(img.URL, req.Width, req.Height, req.X, req.Y)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *ImageService) Rotate(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}
	url, err := s.cldUtil.Rotate(img.URL, req.Angle)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *ImageService) Flip(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}
	url, err := s.cldUtil.Flip(img.URL, req.FlipAxis)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *ImageService) Watermark(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}
	url, err := s.cldUtil.Watermark(img.URL, req.Watermark)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *ImageService) Compress(ctx context.Context, req model.TransformRequest) (string, error) {
	img, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return "", fmt.Errorf("image not found")
	}
	url, err := s.cldUtil.Compress(img.URL, req.Quality)
	if err != nil {
		return "", err
	}
	return url, nil
}
