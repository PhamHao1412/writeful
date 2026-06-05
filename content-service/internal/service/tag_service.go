package service

import (
	"content-service/internal/entity"
	"content-service/internal/model"
	"content-service/internal/repository"
	"context"
)

type ITagService interface {
	CreateTag(ctx context.Context, req model.CreateTagRequest) (*model.TagResponse, error)
	GetTag(ctx context.Context, id string) (*model.TagResponse, error)
	ListTags(ctx context.Context) ([]model.TagResponse, error)
	DeleteTag(ctx context.Context, id string) error
}

type TagService struct {
	tagRepo repository.ITagRepository
}

func NewTagService(tagRepo repository.ITagRepository) *TagService {
	return &TagService{tagRepo}
}

func (s *TagService) CreateTag(ctx context.Context, req model.CreateTagRequest) (*model.TagResponse, error) {
	tag := &entity.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}

	return &model.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

func (s *TagService) GetTag(ctx context.Context, id string) (*model.TagResponse, error) {
	tag, err := s.tagRepo.First(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

func (s *TagService) ListTags(ctx context.Context) ([]model.TagResponse, error) {
	tags, err := s.tagRepo.Find(ctx)
	if err != nil {
		return nil, err
	}

	var responses []model.TagResponse
	for _, tag := range tags {
		responses = append(responses, model.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return responses, nil
}

func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	return s.tagRepo.Delete(ctx, id)
}
