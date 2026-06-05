package service

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJWKService interface {
	GenSignedToken(authTokens model.TokenClaims) (*model.AuthTokenSignedResponse, error)
	VerifyToken(tokenStr string) (*model.JWTClaims, error)
}

type JWKService struct {
	Config *app.Config
}

func NewJWKService(config *app.Config) *JWKService {
	return &JWKService{Config: config}
}

func (s *JWKService) GenSignedToken(authTokens model.TokenClaims) (*model.AuthTokenSignedResponse, error) {
	secretStr := s.Config.JWTConfig.Secret
	if secretStr == "" {
		return nil, errors.New("jwt secret is empty")
	}
	secret, decErr := base64.RawURLEncoding.DecodeString(secretStr)
	if decErr != nil {
		if b, err := base64.StdEncoding.DecodeString(secretStr); err == nil {
			secret = b
		} else {
			secret = []byte(secretStr)
		}
	}

	alg := s.Config.JWTConfig.Alg
	method := jwt.GetSigningMethod(alg)
	if method == nil {
		return nil, errors.New("unsupported signing algorithm")
	}

	kid := s.Config.JWTConfig.Kid
	now := time.Now()
	accessExp := now.Add(time.Duration(s.Config.JWTConfig.DurationInMinutes) * time.Minute)
	refreshExp := now.Add(time.Duration(s.Config.JWTConfig.DurationInMinutesForRefreshToken) * time.Minute)

	accessClaims := model.JWTClaims{
		ID:        authTokens.AccessToken.ID,
		Email:     authTokens.AccessToken.Email,
		Status:    authTokens.AccessToken.Status,
		Roles:     authTokens.AccessToken.Roles,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	access := jwt.NewWithClaims(method, accessClaims)
	if kid != "" {
		access.Header["kid"] = kid
	}
	accessStr, err := access.SignedString(secret)
	if err != nil {
		return nil, errors.New("sign access token error")
	}

	// refresh token claims
	refreshClaims := model.JWTClaims{
		ID:        authTokens.AccessToken.ID,
		Email:     authTokens.AccessToken.Email,
		Status:    authTokens.AccessToken.Status,
		Roles:     authTokens.AccessToken.Roles,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refresh := jwt.NewWithClaims(method, refreshClaims)
	if kid != "" {
		refresh.Header["kid"] = kid
	}
	refreshStr, err := refresh.SignedString(secret)
	if err != nil {
		return nil, errors.New("sign refresh token error")
	}

	return &model.AuthTokenSignedResponse{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		Expiration:   refreshExp.Unix(),
	}, nil
}

func (s *JWKService) VerifyToken(tokenStr string) (*model.JWTClaims, error) {
	secretStr := s.Config.JWTConfig.Secret
	secret, _ := base64.RawURLEncoding.DecodeString(secretStr)
	if len(secret) == 0 {
		if b, err := base64.StdEncoding.DecodeString(secretStr); err == nil {
			secret = b
		} else {
			secret = []byte(secretStr)
		}
	}
	claims := &model.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != s.Config.JWTConfig.Alg {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
