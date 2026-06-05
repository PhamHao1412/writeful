package v1

import (
	"auth-service/internal/handler/rest"
	"auth-service/internal/handler/rest/v1/dto"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userSvc service.IUserService
}

func NewAuthHandler(userSvc service.IUserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// Profile godoc
// @Summary Get current user info
// @Description Get authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /auth/me [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userIdStr := c.GetHeader(rest.HeaderUserID)
	if userIdStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid query parameters")))
		return
	}

	user, err, status := h.userSvc.GetOne(c.Request.Context(), repository.GetUsersFilter{
		ID:       userIdStr,
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(user).WithMessage("user info retrieved successfully"))
}

func (h *AuthHandler) Detail(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid username")))
		return
	}

	user, err, status := h.userSvc.GetOne(c.Request.Context(), repository.GetUsersFilter{
		Username: username,
	})
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(user).WithMessage("user info retrieved successfully"))
}

// List godoc
// @Summary Get list of users
// @Description Get paginated list of users with optional filters
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username query string false "Filter by username"
// @Param email query string false "Filter by email"
// @Param status query string false "Filter by status"
// @Param name query string false "Filter by display name"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (h *AuthHandler) List(c *gin.Context) {
	var req model.GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid query parameters")))
		return
	}

	ctx := c.Request.Context()
	response, total, err := h.userSvc.GetList(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(response).TotalItem(total))
}

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user account with email, username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SignUpRequest true "Sign up request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	keyId := c.GetHeader(rest.HeaderXKeyID)
	var req model.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}
	req.Key = keyId
	response, err, status := h.userSvc.SignUp(c.Request.Context(), req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseOK(response).WithMessage("user registered successfully"))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token and refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	response, err, status := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(response).WithMessage("login successful"))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token and refresh token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}
	ctx := c.Request.Context()
	response, err, status := h.userSvc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(response).WithMessage("token refreshed successfully"))
}

// Logout godoc
// @Summary Logout user
// @Description Revoke all refresh tokens for the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userId := c.GetHeader(rest.HeaderUserID)
	if userId == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	ctx := c.Request.Context()
	err, status := h.userSvc.Logout(ctx, userId)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("logout successful"))
}

func (h *AuthHandler) UpdateInfo(c *gin.Context) {
	userId := c.GetHeader(rest.HeaderUserID)
	if userId == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid request body")))
		return
	}

	ctx := c.Request.Context()
	_, err, status := h.userSvc.UpdateInfo(ctx, userId, req)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("user info updated successfully"))
}

func (h *AuthHandler) UpdateActiveStatus(c *gin.Context) {
	userId := c.GetHeader(rest.HeaderUserID)
	if userId == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}

	var req struct {
		LastActiveAt time.Time `json:"last_active_at" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(err))
		return
	}

	ctx := c.Request.Context()
	err, status := h.userSvc.UpdateActiveStatus(ctx, userId, req.LastActiveAt)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("user active status updated successfully"))
}
