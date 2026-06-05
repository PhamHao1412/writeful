package auth

import (
	"content-service/internal/app"
	"content-service/internal/model"
	"content-service/pkg/util"
	"strings"

	"github.com/sendgrid/rest"
)

const (
	GetUserProfileEndPoint = "/api/v1/user/profile"
	GetListUserEndPoint    = "/api/v1/user/list"
)

type Client interface {
	GetUserProfile(req model.GetUserRequest) (*model.User, error, int)
	GetListUser(req model.GetUserRequest) ([]model.User, error, int)
}

type client struct {
	config *app.Config
}

func NewClient(config *app.Config) Client {
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
