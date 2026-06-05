package v1

import (
	"auth-service/internal/handler/rest"
	"auth-service/internal/handler/rest/v1/dto"
	"auth-service/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowerHandler struct {
	followerSvc service.IFollowerService
}

func NewFollowerHandler(followerSvc service.IFollowerService) *FollowerHandler {
	return &FollowerHandler{followerSvc: followerSvc}
}

// Follow godoc
// @Summary Follow a user
// @Description Follow a user by ID
// @Tags followers
// @Security BearerAuth
// @Param id path string true "User ID to follow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/follow [post]
func (h *FollowerHandler) Follow(c *gin.Context) {
	followerID := c.GetHeader(rest.HeaderUserID)
	if followerID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	followingID := c.Param("id")
	if followingID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}

	err, status := h.followerSvc.Follow(c.Request.Context(), followerID, followingID)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("followed successfully"))
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollow a user by ID
// @Tags followers
// @Security BearerAuth
// @Param id path string true "User ID to unfollow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/{id}/unfollow [delete]
func (h *FollowerHandler) Unfollow(c *gin.Context) {
	followerID := c.GetHeader(rest.HeaderUserID)
	if followerID == "" {
		c.JSON(http.StatusUnauthorized, dto.ResponseUnauthorized(errors.New("unauthorized")))
		return
	}
	followingID := c.Param("id")
	if followingID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}

	err, status := h.followerSvc.Unfollow(c.Request.Context(), followerID, followingID)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(gin.H{}).WithMessage("unfollowed successfully"))
}

// GetFollowers godoc
// @Summary Get followers of a user
// @Description Get list of users following the specified user
// @Tags followers
// @Param id path string true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id}/followers [get]
func (h *FollowerHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err, status := h.followerSvc.GetFollowers(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(users).TotalItem(total))
}

// GetFollowing godoc
// @Summary Get following of a user
// @Description Get list of users followed by the specified user
// @Tags followers
// @Param id path string true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id}/following [get]
func (h *FollowerHandler) GetFollowing(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ResponseBadRequest(errors.New("invalid user id")))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err, status := h.followerSvc.GetFollowing(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(status, dto.ResponseError(err, status))
		return
	}

	c.JSON(http.StatusOK, dto.ResponseOK(users).TotalItem(total))
}
