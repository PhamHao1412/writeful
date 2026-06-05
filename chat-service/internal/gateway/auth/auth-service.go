package auth

import (
	"chat-service/internal/config"
	"chat-service/internal/model"
	"chat-service/pkg/util"
	"strings"
	"time"

	"github.com/sendgrid/rest"
)

const (
	GetUserProfileEndPoint = "/api/v1/user/profile"
	GetListUserEndPoint    = "/api/v1/user/list"
)

type Client interface {
	GetUserProfile(req model.GetUserRequest) (*model.User, error, int)
	GetListUser(req model.GetUserRequest) ([]model.User, error, int)
	UpdateActiveStatus(userID string, lastActiveAt time.Time) (error, int)
}

type client struct {
	config *config.Config
}

func NewClient(config *config.Config) Client {
	return &client{config: config}
}

func (c *client) GetUserProfile(req model.GetUserRequest) (*model.User, error, int) {
	header := map[string]string{
		"x-user-id":    req.ID,
		"Content-Type": "application/json",
	}

	queryParams := map[string]string{
		"id": req.ID,
	}

	baseURL := c.config.AuthServiceURL + GetUserProfileEndPoint
	data, err, statusCode := util.SendRequest(header, nil, queryParams, baseURL, rest.Get)
	if err != nil {
		return nil, err, statusCode
	}

	res, err := util.ParseResponse[*model.User](data)
	if err != nil {
		return nil, err, statusCode
	}

	return res, nil, statusCode
}

func (c *client) GetListUser(req model.GetUserRequest) ([]model.User, error, int) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	queryParams := map[string]string{
		"ids": strings.Join(req.IDs, ","),
	}

	baseURL := c.config.AuthServiceURL + GetListUserEndPoint
	data, err, statusCode := util.SendRequest(header, nil, queryParams, baseURL, rest.Get)
	if err != nil {
		return nil, err, statusCode
	}

	res, err := util.ParseResponse[[]model.User](data)
	if err != nil {
		return nil, err, statusCode
	}
	return res, nil, statusCode
}

func (c *client) UpdateActiveStatus(userID string, lastActiveAt time.Time) (error, int) {
	header := map[string]string{
		"x-user-id":    userID,
		"Content-Type": "application/json",
	}

	body := map[string]interface{}{
		"last_active_at": lastActiveAt,
	}

	baseURL := c.config.AuthServiceURL + "/api/v1/user/active"
	_, err, statusCode := util.SendRequest(header, body, nil, baseURL, rest.Post)
	return err, statusCode
}
