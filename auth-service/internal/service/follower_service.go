package service

import (
	"auth-service/internal/entity"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"context"
	"errors"
	"net/http"
)

type IFollowerService interface {
	Follow(ctx context.Context, followerID, followingID string) (error, int)
	Unfollow(ctx context.Context, followerID, followingID string) (error, int)
	GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]model.UserResponse, int64, error, int)
	GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]model.UserResponse, int64, error, int)
}

type FollowerService struct {
	followerRepo repository.IFollowerRepository
	userRepo     repository.IUserRepository
}

func NewFollowerService(followerRepo repository.IFollowerRepository, userRepo repository.IUserRepository) *FollowerService {
	return &FollowerService{followerRepo: followerRepo, userRepo: userRepo}
}

func (s *FollowerService) Follow(ctx context.Context, followerID, followingID string) (error, int) {
	if followerID == followingID {
		return errors.New("cannot follow yourself"), http.StatusBadRequest
	}

	// Check if user to follow exists
	_, err := s.userRepo.First(ctx, repository.GetUsersFilter{ID: followingID})
	if err != nil {
		return errors.New("user not found"), http.StatusNotFound
	}

	// Check if already following
	isFollowing, err := s.followerRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if isFollowing {
		return errors.New("already following"), http.StatusBadRequest
	}

	err = s.followerRepo.Create(ctx, &entity.Follower{
		FollowerID:  followerID,
		FollowingID: followingID,
	})
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (s *FollowerService) Unfollow(ctx context.Context, followerID, followingID string) (error, int) {
	if followerID == followingID {
		return errors.New("cannot unfollow yourself"), http.StatusBadRequest
	}

	// Check if already following
	isFollowing, err := s.followerRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if !isFollowing {
		return errors.New("not following"), http.StatusBadRequest
	}

	err = s.followerRepo.Delete(ctx, followerID, followingID)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (s *FollowerService) GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]model.UserResponse, int64, error, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, total, err := s.followerRepo.GetFollowers(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err, http.StatusInternalServerError
	}

	return s.toUserResponses(users), total, nil, http.StatusOK
}

func (s *FollowerService) GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]model.UserResponse, int64, error, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, total, err := s.followerRepo.GetFollowing(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err, http.StatusInternalServerError
	}

	return s.toUserResponses(users), total, nil, http.StatusOK
}

func (s *FollowerService) toUserResponses(users []entity.User) []model.UserResponse {
	rs := make([]model.UserResponse, 0, len(users))
	for _, user := range users {
		roles := make([]model.Role, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, model.Role{
				Code: role.Code,
				Name: role.Name,
			})
		}

		rs = append(rs, model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			AvatarURL: user.AvatarURL,
			Roles:     roles,
		})
	}
	return rs
}
