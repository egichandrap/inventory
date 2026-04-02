package handler

import (
	"errors"
	"time"

	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	"github.com/example/jwt-ddd-clean/internal/dto"
)

// TokenHandler handles HTTP requests for token operations
type TokenHandler struct {
	tokenService *service.TokenService
	userService  *UserService
}

// UserService is a placeholder for user authentication service
type UserService struct {
	// In a real application, this would interact with a user repository
}

// NewTokenHandler creates a new TokenHandler
func NewTokenHandler(tokenService *service.TokenService, userService *UserService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
		userService:  userService,
	}
}

// GenerateToken handles token generation requests
// This is a simplified version - in production, you'd validate credentials
func (h *TokenHandler) GenerateToken(username, password string) (*dto.TokenResponse, error) {
	// In a real application, authenticate the user first
	user := &model.User{
		ID:       "user-123",
		Username: username,
		Email:    username + "@example.com",
		Role:     "user",
	}

	tokenPair, err := h.tokenService.GenerateTokens(nil, user)
	if err != nil {
		return nil, err
	}

	expiresIn := int64(time.Until(tokenPair.Access.ExpiresAt).Seconds())

	return &dto.TokenResponse{
		AccessToken:  tokenPair.Access.AccessToken,
		RefreshToken: tokenPair.Refresh.AccessToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken handles token refresh requests
func (h *TokenHandler) RefreshToken(refreshToken string) (*dto.TokenResponse, error) {
	tokenPair, err := h.tokenService.RefreshToken(nil, refreshToken)
	if err != nil {
		return nil, err
	}

	expiresIn := int64(time.Until(tokenPair.Access.ExpiresAt).Seconds())

	return &dto.TokenResponse{
		AccessToken:  tokenPair.Access.AccessToken,
		RefreshToken: tokenPair.Refresh.AccessToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken handles token validation requests
func (h *TokenHandler) ValidateToken(token string) (*dto.ValidateTokenResponse, error) {
	claims, err := h.tokenService.ValidateToken(nil, token)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == apperrors.ErrInvalidToken || appErr.Code == apperrors.ErrExpiredToken {
				return &dto.ValidateTokenResponse{
					Valid: false,
				}, nil
			}
		}
		return nil, err
	}

	return &dto.ValidateTokenResponse{
		Valid:    true,
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

// RevokeToken handles token revocation requests
func (h *TokenHandler) RevokeToken(token string) error {
	return h.tokenService.RevokeToken(nil, token)
}
