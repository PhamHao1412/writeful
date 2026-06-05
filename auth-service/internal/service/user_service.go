package service

import (
	"auth-service/internal/app"
	"auth-service/internal/entity"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

const (
	defaultUserAvt = "http://res.cloudinary.com/dlqwa0yhj/image/upload/v1768639967/image-service/originals/images/am3sxdok4wrbd16l5mrt.jpg"
)

type IUserService interface {
	GetOne(ctx context.Context, filter repository.GetUsersFilter) (*entity.User, error, int)
	GetList(ctx context.Context, req model.GetUsersRequest) ([]model.UserResponse, int64, error)
	SignUp(ctx context.Context, req model.SignUpRequest) (*model.UserInfo, error, int)
	Login(ctx context.Context, req model.LoginRequest) (*model.AuthTokenSignedResponse, error, int)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokenSignedResponse, error, int)
	Logout(ctx context.Context, userID string) (error, int)
	UpdateInfo(ctx context.Context, userID string, req model.UpdateUserRequest) (*model.UserInfo, error, int)
	UpdateActiveStatus(ctx context.Context, userID string, lastActiveAt time.Time) (error, int)
}

type UserService struct {
	config      *app.Config
	userRepo    repository.IUserRepository
	roleRepo    repository.IRoleRepository
	refreshRepo repository.IRefreshTokenRepository
	jwkService  IJWKService
}

func NewUserService(config *app.Config, userRepo repository.IUserRepository, roleRepo repository.IRoleRepository, refreshTokenRepo repository.IRefreshTokenRepository, jwkService IJWKService) *UserService {
	return &UserService{
		config:      config,
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		refreshRepo: refreshTokenRepo,
		jwkService:  jwkService,
	}
}

func (s *UserService) GetOne(ctx context.Context, filter repository.GetUsersFilter) (*entity.User, error, int) {
	user, err := s.userRepo.First(ctx, repository.GetUsersFilter{
		ID:       filter.ID,
		Username: filter.Username,
		Email:    filter.Email,
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return user, nil, http.StatusOK
}

func (s *UserService) GetList(ctx context.Context, req model.GetUsersRequest) ([]model.UserResponse, int64, error) {

	// Build filter
	filter := repository.GetUsersFilter{
		IDs:    strings.Split(req.IDs, ","),
		Limit:  req.PageSize,
		Offset: req.Page,
	}

	users, total, err := s.userRepo.Find(ctx, filter)
	if err != nil {
		return nil, 0, errors.New("failed to get users")
	}

	// Convert to response
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
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Status:       user.Status,
			AvatarURL:    user.AvatarURL,
			DisplayName:  user.DisplayName,
			LastActiveAt: user.LastActiveAt,
			Roles:        roles,
		})
	}

	return rs, total, nil

}

func (s *UserService) SignUp(ctx context.Context, req model.SignUpRequest) (*model.UserInfo, error, int) {
	// Check if email already exists
	existingUserByEmail, err := s.userRepo.First(ctx, repository.GetUsersFilter{
		Email: req.Email,
	})
	if err == nil && existingUserByEmail != nil {
		return nil, errors.New("user with this email already exists"), http.StatusBadRequest
	}

	// Check if username already exists
	existingUserByUsername, err := s.userRepo.First(ctx, repository.GetUsersFilter{
		Username: req.Username,
	})
	if err == nil && existingUserByUsername != nil {
		return nil, errors.New("user with this username already exists"), http.StatusBadRequest
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password"), http.StatusInternalServerError
	}

	if req.AvatarURL == "" {
		req.AvatarURL = defaultUserAvt
	}
	// Create user
	user := &entity.User{
		Email:        req.Email,
		Username:     "@" + req.Username,
		AvatarURL:    req.AvatarURL,
		DisplayName:  req.DisplayName,
		Bio:          req.Bio,
		PasswordHash: string(hashedPassword),
		Status:       "active",
	}

	roleCode := "writer"
	if req.Key == s.config.JWTConfig.AdminKey {
		roleCode = "super_admin"
	}

	// Get default user role
	defaultRole, err := s.roleRepo.First(ctx, roleCode)
	if err != nil {

		return nil, errors.New("default role not found"), http.StatusInternalServerError
	}

	// Create user and assign role in transaction
	if err = s.userRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err = tx.Create(user).Error; err != nil {
			return err
		}

		userRole := &entity.UserRole{
			UserID: user.ID,
			RoleID: defaultRole.ID,
		}

		if err := tx.Create(userRole).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.New("failed to create user"), http.StatusInternalServerError
	}

	// Get user with roles
	user, err = s.userRepo.First(ctx, repository.GetUsersFilter{
		ID: user.ID.String(),
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &model.UserInfo{
		ID:     user.ID,
		Email:  user.Email,
		Status: user.Status,
	}, nil, http.StatusOK
}

func (s *UserService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthTokenSignedResponse, error, int) {
	// Find user by email
	user, err := s.userRepo.First(ctx, repository.GetUsersFilter{
		Email: req.Email,
	})
	if err != nil {
		return nil, errors.New("invalid email or password"), http.StatusUnauthorized
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user account is not active"), http.StatusUnauthorized
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password"), http.StatusUnauthorized
	}

	// Generate tokens

	roles := s.extractRoles(user)
	response, err := s.jwkService.GenSignedToken(
		model.TokenClaims{
			AccessToken: model.Token{
				ID:     user.ID.String(),
				Email:  user.Email,
				Status: user.Status,
				Roles:  roles,
			},
			RefreshToken: model.Token{
				ID:     user.ID.String(),
				Email:  user.Email,
				Status: user.Status,
				Roles:  roles,
			},
		})
	if err != nil {
		return nil, errors.New("failed to generate access token"), http.StatusInternalServerError
	}

	refreshToken := entity.RefreshToken{
		UserID:    user.ID,
		Token:     response.RefreshToken,
		ExpiredAt: response.Expiration,
		IsRevoked: false,
	}

	if err = s.refreshRepo.Create(ctx, &refreshToken); err != nil {
		return nil, errors.New("failed to create refresh token"), http.StatusInternalServerError
	}

	return &model.AuthTokenSignedResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		Expiration:   response.Expiration,
	}, nil, http.StatusOK
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokenSignedResponse, error, int) {
	claims, err := s.jwkService.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid token"), http.StatusUnauthorized
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token"), http.StatusUnauthorized
	}

	oldToken, err := s.refreshRepo.First(ctx, repository.GetTokensFilter{
		Token: refreshToken,
	})
	if err != nil {
		_ = s.refreshRepo.RevokeAll(ctx, claims.ID)
		return nil, errors.New("invalid token"), http.StatusUnauthorized
	}

	now := time.Now().Unix()
	if oldToken.IsRevoked || oldToken.ExpiredAt < now {
		_ = s.refreshRepo.RevokeAll(ctx, claims.ID)
		return nil, errors.New("token expired or revoked"), http.StatusUnauthorized
	}

	claimBytes, err := json.Marshal(claims)
	if err != nil {
		return nil, errors.New("failed to parse token claims"), http.StatusInternalServerError
	}

	// revoke old token
	if err := s.refreshRepo.Revoke(ctx, oldToken.ID.String()); err != nil {
		return nil, errors.New("failed to revoke refresh token"), http.StatusInternalServerError
	}

	var user model.UserInfo
	if err = json.Unmarshal(claimBytes, &user); err != nil {
		return nil, errors.New("failed to parse token claims"), http.StatusInternalServerError
	}

	response, err := s.jwkService.GenSignedToken(
		model.TokenClaims{
			AccessToken: model.Token{
				ID:     user.ID.String(),
				Email:  user.Email,
				Status: user.Status,
				Roles:  user.Roles,
			},
			RefreshToken: model.Token{
				ID:     user.ID.String(),
				Email:  user.Email,
				Status: user.Status,
				Roles:  user.Roles,
			},
		})
	if err != nil {
		return nil, errors.New("failed to generate access token"), http.StatusInternalServerError
	}

	newToken := entity.RefreshToken{
		UserID:    user.ID,
		Token:     response.RefreshToken,
		ExpiredAt: response.Expiration,
		IsRevoked: false,
	}

	if err = s.refreshRepo.Create(ctx, &newToken); err != nil {
		return nil, errors.New("failed to create refresh token"), http.StatusInternalServerError
	}

	return &model.AuthTokenSignedResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		Expiration:   response.Expiration,
	}, nil, http.StatusOK
}

func (s *UserService) extractRoles(user *entity.User) []string {
	roles := make([]string, 0)
	for _, userRole := range user.Roles {
		roles = append(roles, userRole.Code)
	}
	return roles
}

func (s *UserService) Logout(ctx context.Context, userID string) (error, int) {
	// Revoke all refresh tokens for this user
	if err := s.refreshRepo.RevokeAll(ctx, userID); err != nil {
		return errors.New("failed to logout"), http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (s *UserService) UpdateInfo(ctx context.Context, userID string, req model.UpdateUserRequest) (*model.UserInfo, error, int) {
	user, err := s.userRepo.First(ctx, repository.GetUsersFilter{
		ID: userID,
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	user.DisplayName = req.DisplayName
	user.Bio = req.Bio
	user.AvatarURL = req.AvatarURL
	if err = s.userRepo.Update(ctx, user); err != nil {
	}
	return nil, nil, http.StatusOK
}

func (s *UserService) UpdateActiveStatus(ctx context.Context, userID string, lastActiveAt time.Time) (error, int) {
	err := s.userRepo.GetDB().WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("last_active_at", lastActiveAt).Error
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
