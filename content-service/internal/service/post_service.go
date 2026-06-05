package service

import (
	"content-service/internal/entity"
	"content-service/internal/gateway/auth"
	"content-service/internal/model"
	"content-service/internal/repository"
	"content-service/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*entity.Post, error, int)
	UpdatePost(ctx context.Context, userID uuid.UUID, postID string, req model.UpdatePostRequest) error
	GetOne(ctx context.Context, postID string) (*model.PostResponse, error, int)
	GetPostBySlug(ctx context.Context, slug string) (*model.PostResponse, error)
	GetList(ctx context.Context, req repository.GetPostFilter) ([]model.PostResponse, int64, error)
	Publish(ctx context.Context, postID, userID string) (error, int)
	Unpublish(ctx context.Context, postID, userID string) (error, int)
	DeletePost(ctx context.Context, postID string) error
}

type PostService struct {
	baseRepo    repository.IBaseRepository
	postRepo    repository.IPostRepository
	versionRepo repository.IPostVersionRepository
	mediaRepo   repository.IMediaRepository
	tagRepo     repository.ITagRepository
	authClient  auth.Client
}

func NewPostService(
	baseRepo repository.IBaseRepository,
	postRepo repository.IPostRepository,
	versionRepo repository.IPostVersionRepository,
	mediaRepo repository.IMediaRepository,
	tagRepo repository.ITagRepository,
	authClient auth.Client,
) *PostService {
	return &PostService{baseRepo, postRepo, versionRepo, mediaRepo, tagRepo, authClient}
}

func (s *PostService) CreatePost(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*entity.Post, error, int) {
	uniqueSlug, err := s.generateUniqueSlug(ctx, req.Title)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	var post *entity.Post
	err = s.postRepo.Transaction(ctx, func(tx *gorm.DB) error {
		post = &entity.Post{
			UserId:        userID,
			Title:         req.Title,
			Slug:          uniqueSlug,
			Visibility:    req.Visibility,
			Excerpt:       req.Excerpt,
			Status:        "draft",
			CoverImageURL: req.CoverImageURL,
		}
		if err := s.baseRepo.CreateTx(tx, post); err != nil {
			return err
		}

		version := &entity.PostVersion{
			PostID:    post.ID,
			Version:   1,
			Content:   req.Content,
			CreatedBy: userID,
		}
		if err := s.baseRepo.CreateTx(tx, version); err != nil {
			return err
		}

		if len(req.TagIDs) > 0 {
			for _, tagID := range req.TagIDs {
				postTag := &entity.PostTag{
					PostID: post.ID,
					TagID:  tagID,
				}
				if err := s.baseRepo.CreateTx(tx, postTag); err != nil {
					return err
				}
			}
		}

		if len(req.MediaIDs) > 0 {
			for i, mediaID := range req.MediaIDs {
				postMedia := &entity.PostMedia{
					PostID:       post.ID,
					MediaID:      mediaID,
					DisplayOrder: i,
				}
				if err := s.baseRepo.CreateTx(tx, postMedia); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return post, nil, http.StatusCreated
}

func (s *PostService) UpdatePost(ctx context.Context, userID uuid.UUID, postID string, req model.UpdatePostRequest) error {
	post, err := s.postRepo.First(ctx, repository.GetPostFilter{ID: postID})
	if err != nil {
		return err
	}

	if post.UserId != userID {
		return fmt.Errorf("unauthorized")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	post.Subtitle = req.Subtitle
	post.CoverImageURL = req.CoverImageURL

	if err := s.postRepo.Update(ctx, post); err != nil {
		return err
	}

	if req.Content != "" {
		latestVersion, _ := s.versionRepo.First(ctx, repository.GetPostVersionFilter{
			PostID: postID,
			Sort:   "version DESC",
		})

		newVersion := &entity.PostVersion{
			PostID:    post.ID,
			Version:   latestVersion.Version + 1,
			Content:   req.Content,
			CreatedBy: userID,
		}

		if err := s.versionRepo.Create(ctx, newVersion); err != nil {
			return err
		}
	}

	if len(req.TagIDs) > 0 {
		if err := s.postRepo.RemoveTags(ctx, post.ID); err != nil {
			return err
		}
		if err := s.postRepo.AddTags(ctx, post.ID, req.TagIDs); err != nil {
			return err
		}
	}

	if len(req.MediaIDs) > 0 {
		if err := s.postRepo.RemoveMedias(ctx, post.ID); err != nil {
			return err
		}
		if err := s.postRepo.AddMedias(ctx, post.ID, req.MediaIDs); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostService) GetOne(ctx context.Context, postID string) (*model.PostResponse, error, int) {
	post, err := s.postRepo.First(ctx, repository.GetPostFilter{ID: postID})
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err, http.StatusNotFound
	}

	userProfile, err, status := s.authClient.GetUserProfile(model.GetUserRequest{
		ID: post.UserId.String(),
	})
	if err != nil {
		return nil, err, status
	}

	version, err := s.versionRepo.First(ctx, repository.GetPostVersionFilter{
		PostID: post.ID.String(),
		Sort:   "version DESC",
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var media []model.Media
	for _, m := range post.PostMedia {
		media = append(media, model.Media{
			PostID:       m.PostID.String(),
			MediaID:      m.MediaID.String(),
			DisplayOrder: m.DisplayOrder,
			URL:          m.Media.URL,
		})
	}

	resp := &model.PostResponse{
		ID:            post.ID,
		UserID:        post.UserId,
		Title:         post.Title,
		Subtitle:      post.Subtitle,
		Slug:          post.Slug,
		Excerpt:       post.Excerpt,
		Content:       version.Content,
		Status:        post.Status,
		PublishedAt:   post.PublishedAt,
		CoverImageURL: post.CoverImageURL,
		User:          userProfile,
		Version:       version.Version,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		Media:         media,
	}

	return resp, nil, http.StatusOK
}

func (s *PostService) GetPostBySlug(ctx context.Context, slug string) (*model.PostResponse, error) {
	post, err := s.postRepo.First(ctx, repository.GetPostFilter{Slug: slug, Status: "published"})
	if err != nil {
		return nil, err
	}

	return s.buildPostResponse(ctx, post)
}

func (s *PostService) GetList(ctx context.Context, req repository.GetPostFilter) ([]model.PostResponse, int64, error) {
	if req.Status == "" {
		req.Status = "published"
	}

	posts, total, err := s.postRepo.Find(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println("posts", posts)
	if len(posts) == 0 {
		return []model.PostResponse{}, total, nil
	}

	postMap := make(map[string]*entity.Post, len(posts))
	userSet := make(map[string]struct{})
	userIds := make([]string, 0, len(posts))
	for _, post := range posts {
		if _, ok := userSet[post.UserId.String()]; !ok {
			userIds = append(userIds, post.UserId.String())
			userSet[post.UserId.String()] = struct{}{}
		}
		postMap[post.ID.String()] = &post

	}

	users, err, _ := s.authClient.GetListUser(model.GetUserRequest{
		IDs: userIds,
	})
	if err != nil {
		logger.Error(ctx, "[GetList] get list user failed", zap.Error(err))
		return nil, 0, err
	}
	if len(users) == 0 {
		logger.Error(ctx, "[GetList] get list user empty")
		return nil, 0, errors.New("get list user empty")
	}
	userMap := make(map[string]*model.User, len(users))
	for _, user := range users {
		userMap[user.ID] = &user
	}

	responses := make([]model.PostResponse, 0, len(posts))
	for _, post := range posts {

		version, err := s.versionRepo.First(ctx, repository.GetPostVersionFilter{
			PostID: post.ID.String(),

			Sort: "version DESC",
		})
		if err != nil {
			return nil, 0, err
		}

		user, ok := userMap[post.UserId.String()]
		if !ok {
			logger.Error(ctx, "[GetList] user not found", zap.String("user_id", post.UserId.String()))
			return nil, 0, err
		}

		var media []model.Media
		for _, m := range post.PostMedia {
			media = append(media, model.Media{
				PostID:       m.PostID.String(),
				MediaID:      m.MediaID.String(),
				DisplayOrder: m.DisplayOrder,
				URL:          m.Media.URL,
			})
		}

		resp := &model.PostResponse{
			ID:            post.ID,
			UserID:        post.UserId,
			Title:         post.Title,
			Subtitle:      post.Subtitle,
			Slug:          post.Slug,
			Excerpt:       post.Excerpt,
			Content:       version.Content,
			Status:        post.Status,
			PublishedAt:   post.PublishedAt,
			CoverImageURL: post.CoverImageURL,
			Version:       version.Version,
			CreatedAt:     post.CreatedAt,
			UpdatedAt:     post.UpdatedAt,
			User:          user,
			Media:         media,
		}
		responses = append(responses, *resp)
	}

	return responses, total, nil
}

func (s *PostService) Publish(ctx context.Context, postID, userID string) (error, int) {
	now := time.Now()
	if err := s.postRepo.UpdateStatus(ctx, postID, userID, "published", &now); err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (s *PostService) Unpublish(ctx context.Context, postID, userID string) (error, int) {
	if err := s.postRepo.UpdateStatus(ctx, postID, userID, "draft", nil); err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (s *PostService) DeletePost(ctx context.Context, postID string) error {
	return s.postRepo.Delete(ctx, postID)
}

func (p *PostService) generateUniqueSlug(ctx context.Context, title string) (string, error) {
	baseSlug := slug.Make(title)
	if baseSlug == "" {
		baseSlug = "untitled"
	}
	finalSlug := baseSlug

	exists := true
	suffix := 1
	for exists {
		count, err := p.postRepo.Count(ctx, repository.GetPostFilter{
			ExactSlug: finalSlug,
		})
		if err != nil {
			return "", err
		}
		if count == 0 {
			exists = false
		} else {
			suffix++
			finalSlug = fmt.Sprintf("%s-%d", baseSlug, suffix)
		}
	}

	return finalSlug, nil
}

func (s *PostService) buildPostResponse(ctx context.Context, post *entity.Post) (*model.PostResponse, error) {
	version, err := s.versionRepo.First(ctx, repository.GetPostVersionFilter{
		PostID: post.ID.String(),
		Sort:   "version DESC",
	})
	if err != nil {
		return nil, err
	}

	resp := &model.PostResponse{
		ID:            post.ID,
		UserID:        post.UserId,
		Title:         post.Title,
		Slug:          post.Slug,
		Excerpt:       post.Excerpt,
		Content:       version.Content,
		Status:        post.Status,
		PublishedAt:   post.PublishedAt,
		CoverImageURL: post.CoverImageURL,
		Version:       version.Version,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}

	tags, _ := s.tagRepo.FindByPostID(ctx, post.ID)
	for _, tag := range tags {
		resp.Tags = append(resp.Tags, model.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return resp, nil
}
