package service

import (
	"content-service/internal/entity"
	"content-service/internal/gateway/auth"
	"content-service/internal/model"
	"content-service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type IStoryService interface {
	CreateStory(ctx context.Context, userID uuid.UUID, req model.CreateStoryRequest) (*entity.Story, error)
	GetFeedStories(ctx context.Context, viewerID uuid.UUID) ([]model.UserStoriesGroup, error)
	MarkStoryAsSeen(ctx context.Context, storyID, viewerID uuid.UUID) error
	DeleteStory(ctx context.Context, storyID, userID uuid.UUID) error
}

type StoryService struct {
	storyRepo  repository.IStoryRepository
	authClient auth.Client
}

func NewStoryService(storyRepo repository.IStoryRepository, authClient auth.Client) *StoryService {
	return &StoryService{
		storyRepo:  storyRepo,
		authClient: authClient,
	}
}

func (s *StoryService) CreateStory(ctx context.Context, userID uuid.UUID, req model.CreateStoryRequest) (*entity.Story, error) {
	story := &entity.Story{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        "image",
		MediaURL:    req.MediaURL,
		Caption:     req.Caption,
		AudioURL:    req.AudioURL,
		AudioTitle:  req.AudioTitle,
		AudioArtist: req.AudioArtist,
		AudioOffset: req.AudioOffset,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Status:      "active",
	}

	if err := s.storyRepo.Create(ctx, story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *StoryService) GetFeedStories(ctx context.Context, viewerID uuid.UUID) ([]model.UserStoriesGroup, error) {
	stories, err := s.storyRepo.GetActiveStories(ctx)
	if err != nil {
		return nil, err
	}

	seenMap, err := s.storyRepo.GetSeenStoryIDs(ctx, viewerID)
	if err != nil {
		seenMap = make(map[uuid.UUID]bool)
	}

	// Group stories by UserID
	userGroups := make(map[uuid.UUID][]entity.Story)
	var userIDs []uuid.UUID
	seenUserIDs := make(map[uuid.UUID]bool)

	for _, st := range stories {
		if _, ok := userGroups[st.UserID]; !ok {
			userGroups[st.UserID] = []entity.Story{}
			if !seenUserIDs[st.UserID] {
				userIDs = append(userIDs, st.UserID)
				seenUserIDs[st.UserID] = true
			}
		}
		userGroups[st.UserID] = append(userGroups[st.UserID], st)
	}

	// Fetch User metadata from Auth Service
	usersMap := make(map[string]model.User)
	if len(userIDs) > 0 {
		var idsStr []string
		for _, uid := range userIDs {
			idsStr = append(idsStr, uid.String())
		}
		// Get users list
		users, err, _ := s.authClient.GetListUser(model.GetUserRequest{IDs: idsStr})
		if err == nil {
			for _, u := range users {
				usersMap[u.ID] = u
			}
		}
	}

	var result []model.UserStoriesGroup
	for _, uid := range userIDs {
		gStories := userGroups[uid]
		var dStories []model.StoryDisplayDTO
		hasUnread := false

		for _, st := range gStories {
			seen := seenMap[st.ID]
			if !seen && uid != viewerID {
				hasUnread = true
			}
			dStories = append(dStories, model.StoryDisplayDTO{
				ID:          st.ID.String(),
				UserID:      st.UserID.String(),
				Type:        st.Type,
				MediaURL:    st.MediaURL,
				Caption:     st.Caption,
				AudioURL:    st.AudioURL,
				AudioTitle:  st.AudioTitle,
				AudioArtist: st.AudioArtist,
				AudioOffset: st.AudioOffset,
				CreatedAt:   st.CreatedAt,
				ExpiresAt:   st.ExpiresAt,
				Seen:        seen,
			})
		}

		uMeta, ok := usersMap[uid.String()]
		username := "User"
		avatar := "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=100" // default premium placeholder
		if ok {
			if uMeta.Username != "" {
				username = uMeta.Username
			}
			if uMeta.AvtarURL != "" {
				avatar = uMeta.AvtarURL
			}
		}

		result = append(result, model.UserStoriesGroup{
			UserID:    uid.String(),
			Username:  username,
			AvatarURL: avatar,
			HasUnread: hasUnread,
			Stories:   dStories,
		})
	}

	return result, nil
}

func (s *StoryService) MarkStoryAsSeen(ctx context.Context, storyID, viewerID uuid.UUID) error {
	return s.storyRepo.MarkAsSeen(ctx, storyID, viewerID)
}

func (s *StoryService) DeleteStory(ctx context.Context, storyID, userID uuid.UUID) error {
	return s.storyRepo.Delete(ctx, storyID, userID)
}
