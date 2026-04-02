// Package jwt provides JWT token utilities for external use
package jwt

import (
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	"github.com/example/jwt-ddd-clean/internal/handler"
	infrastructurejwt "github.com/example/jwt-ddd-clean/internal/infrastructure/jwt"
	repo "github.com/example/jwt-ddd-clean/internal/infrastructure/repository"
)

// Config holds the JWT configuration
type Config struct {
	SecretKey       string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// JWT provides JWT token operations
type JWT struct {
	handler *handler.TokenHandler
	config  Config
}

// New creates a new JWT instance
func New(config Config) *JWT {
	// Infrastructure layer
	jwtProvider := infrastructurejwt.NewProvider(infrastructurejwt.Config{
		SecretKey: config.SecretKey,
		Issuer:    config.Issuer,
		Algorithm: "HS256",
	})

	tokenRepository := repo.NewMemoryTokenRepository()

	// Domain layer
	tokenService := service.NewTokenService(
		tokenRepository,
		jwtProvider,
		config.AccessTokenTTL,
		config.RefreshTokenTTL,
	)

	// Handler layer
	tokenHandler := handler.NewTokenHandler(tokenService, &handler.UserService{})

	return &JWT{
		handler: tokenHandler,
		config:  config,
	}
}

// GenerateToken generates a new JWT token pair
func (j *JWT) GenerateToken(username, password string) (*TokenPair, error) {
	response, err := j.handler.GenerateToken(username, password)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
	}, nil
}

// ValidateToken validates a JWT token
func (j *JWT) ValidateToken(token string) (*TokenClaims, error) {
	response, err := j.handler.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if !response.Valid {
		return nil, ErrInvalidToken
	}

	return &TokenClaims{
		UserID:   response.UserID,
		Username: response.Username,
		Role:     response.Role,
	}, nil
}

// RefreshToken refreshes an expired access token
func (j *JWT) RefreshToken(refreshToken string) (*TokenPair, error) {
	response, err := j.handler.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
	}, nil
}

// RevokeToken revokes a token
func (j *JWT) RevokeToken(token string) error {
	return j.handler.RevokeToken(token)
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID   string
	Username string
	Role     string
}

// ErrInvalidToken is returned when a token is invalid
var ErrInvalidToken = model.ErrInvalidToken
