package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// MockJWTProvider is a mock implementation of JWTProvider
type MockJWTProvider struct {
	generateTokenFunc func(claims *model.TokenClaims, expiresAt time.Time) (string, error)
	validateTokenFunc func(token string) (*model.TokenClaims, error)
	getExpirationFunc func(token string) (time.Time, error)
}

func (m *MockJWTProvider) GenerateToken(claims *model.TokenClaims, expiresAt time.Time) (string, error) {
	if m.generateTokenFunc != nil {
		return m.generateTokenFunc(claims, expiresAt)
	}
	return "mock-token", nil
}

func (m *MockJWTProvider) ValidateToken(token string) (*model.TokenClaims, error) {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(token)
	}
	return &model.TokenClaims{
		UserID:   "user-123",
		Username: "testuser",
		Role:     "user",
	}, nil
}

func (m *MockJWTProvider) GetExpiration(token string) (time.Time, error) {
	if m.getExpirationFunc != nil {
		return m.getExpirationFunc(token)
	}
	return time.Now().Add(24 * time.Hour), nil
}

// MockTokenRepository is a mock implementation of TokenRepository
type MockTokenRepository struct {
	storeFunc         func(ctx context.Context, userID string, token string, tokenType string, expiresAt time.Time) error
	findFunc          func(ctx context.Context, userID string, tokenType string) (*model.Token, error)
	deleteFunc        func(ctx context.Context, userID string, tokenType string) error
	isBlacklistedFunc func(ctx context.Context, token string) (bool, error)
	blacklistFunc     func(ctx context.Context, token string, expiresAt time.Time) error
}

func (m *MockTokenRepository) Store(ctx context.Context, userID string, token string, tokenType string, expiresAt time.Time) error {
	if m.storeFunc != nil {
		return m.storeFunc(ctx, userID, token, tokenType, expiresAt)
	}
	return nil
}

func (m *MockTokenRepository) Find(ctx context.Context, userID string, tokenType string) (*model.Token, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, userID, tokenType)
	}
	return nil, nil
}

func (m *MockTokenRepository) Delete(ctx context.Context, userID string, tokenType string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, userID, tokenType)
	}
	return nil
}

func (m *MockTokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	if m.isBlacklistedFunc != nil {
		return m.isBlacklistedFunc(ctx, token)
	}
	return false, nil
}

func (m *MockTokenRepository) Blacklist(ctx context.Context, token string, expiresAt time.Time) error {
	if m.blacklistFunc != nil {
		return m.blacklistFunc(ctx, token, expiresAt)
	}
	return nil
}

func TestTokenService_GenerateTokens(t *testing.T) {
	t.Run("should generate token pair successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		user := model.NewUser("user-123", "testuser", "test@example.com", "user")

		// Act
		tokenPair, err := tokenService.GenerateTokens(context.Background(), user)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, tokenPair)
		assert.NotNil(t, tokenPair.Access)
		assert.NotNil(t, tokenPair.Refresh)
		assert.NotEmpty(t, tokenPair.Access.AccessToken)
		assert.NotEmpty(t, tokenPair.Refresh.AccessToken)
	})

	t.Run("should return error when jwt provider fails to generate access token", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{}
		mockJWT := &MockJWTProvider{
			generateTokenFunc: func(claims *model.TokenClaims, expiresAt time.Time) (string, error) {
				return "", apperrors.ErrTokenGenerationErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)
		user := model.NewUser("user-123", "testuser", "test@example.com", "user")

		// Act
		tokenPair, err := tokenService.GenerateTokens(context.Background(), user)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
		assert.ErrorIs(t, err, apperrors.ErrTokenGenerationErr)
	})

	t.Run("should return error when jwt provider fails to generate refresh token", func(t *testing.T) {
		// Arrange
		callCount := 0
		mockRepo := &MockTokenRepository{}
		mockJWT := &MockJWTProvider{
			generateTokenFunc: func(claims *model.TokenClaims, expiresAt time.Time) (string, error) {
				callCount++
				// First call (access token) succeeds
				if callCount == 1 {
					return "access-token", nil
				}
				// Second call (refresh token) fails
				return "", apperrors.ErrTokenGenerationErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)
		user := model.NewUser("user-123", "testuser", "test@example.com", "user")

		// Act
		tokenPair, err := tokenService.GenerateTokens(context.Background(), user)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
	})

	t.Run("should return error when storing refresh token fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			storeFunc: func(ctx context.Context, userID string, token string, tokenType string, expiresAt time.Time) error {
				return apperrors.ErrTokenStorageErr
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)
		user := model.NewUser("user-123", "testuser", "test@example.com", "user")

		// Act
		tokenPair, err := tokenService.GenerateTokens(context.Background(), user)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
		assert.ErrorIs(t, err, apperrors.ErrTokenStorageErr)
	})
}

func TestTokenService_ValidateToken(t *testing.T) {
	t.Run("should validate token successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		claims, err := tokenService.ValidateToken(context.Background(), "valid-token")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "user-123", claims.UserID)
		assert.Equal(t, "testuser", claims.Username)
		assert.Equal(t, "user", claims.Role)
	})

	t.Run("should return error when token is blacklisted", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return true, nil
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		claims, err := tokenService.ValidateToken(context.Background(), "blacklisted-token")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrRevokedTokenErr, err)
		assert.Nil(t, claims)
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
		}
		mockJWT := &MockJWTProvider{
			validateTokenFunc: func(token string) (*model.TokenClaims, error) {
				return nil, apperrors.ErrInvalidTokenErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		claims, err := tokenService.ValidateToken(context.Background(), "invalid-token")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrInvalidTokenErr, err)
		assert.Nil(t, claims)
	})

	t.Run("should return error when token is expired", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
		}
		mockJWT := &MockJWTProvider{
			validateTokenFunc: func(token string) (*model.TokenClaims, error) {
				return nil, apperrors.ErrExpiredTokenErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		claims, err := tokenService.ValidateToken(context.Background(), "expired-token")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrExpiredTokenErr, err)
		assert.Nil(t, claims)
	})

	t.Run("should return error when checking blacklist fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, apperrors.ErrInternalErr
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		claims, err := tokenService.ValidateToken(context.Background(), "token")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestTokenService_RefreshToken(t *testing.T) {
	t.Run("should refresh token successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
			blacklistFunc: func(ctx context.Context, token string, expiresAt time.Time) error {
				return nil
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		tokenPair, err := tokenService.RefreshToken(context.Background(), "refresh-token")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, tokenPair)
		assert.NotNil(t, tokenPair.Access)
		assert.NotNil(t, tokenPair.Refresh)
	})

	t.Run("should return error when refresh token is invalid", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{}
		mockJWT := &MockJWTProvider{
			validateTokenFunc: func(token string) (*model.TokenClaims, error) {
				return nil, apperrors.ErrInvalidTokenErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		tokenPair, err := tokenService.RefreshToken(context.Background(), "invalid-refresh-token")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrInvalidTokenErr, err)
		assert.Nil(t, tokenPair)
	})

	t.Run("should return error when refresh token is expired", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{}
		mockJWT := &MockJWTProvider{
			validateTokenFunc: func(token string) (*model.TokenClaims, error) {
				return nil, apperrors.ErrExpiredTokenErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		tokenPair, err := tokenService.RefreshToken(context.Background(), "expired-refresh-token")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrExpiredTokenErr, err)
		assert.Nil(t, tokenPair)
	})

	t.Run("should return error when refresh token is blacklisted", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return true, nil
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		tokenPair, err := tokenService.RefreshToken(context.Background(), "blacklisted-refresh-token")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrRevokedTokenErr, err)
		assert.Nil(t, tokenPair)
	})

	t.Run("should return error when blacklisting old token fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			isBlacklistedFunc: func(ctx context.Context, token string) (bool, error) {
				return false, nil
			},
			blacklistFunc: func(ctx context.Context, token string, expiresAt time.Time) error {
				return apperrors.ErrTokenStorageErr
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		tokenPair, err := tokenService.RefreshToken(context.Background(), "refresh-token")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
		assert.ErrorIs(t, err, apperrors.ErrTokenStorageErr)
	})
}

func TestTokenService_RevokeToken(t *testing.T) {
	t.Run("should revoke token successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			blacklistFunc: func(ctx context.Context, token string, expiresAt time.Time) error {
				return nil
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		err := tokenService.RevokeToken(context.Background(), "token-to-revoke")

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error when getting expiration fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{}
		mockJWT := &MockJWTProvider{
			getExpirationFunc: func(token string) (time.Time, error) {
				return time.Time{}, apperrors.ErrInvalidTokenErr
			},
		}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		err := tokenService.RevokeToken(context.Background(), "token")

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrInvalidTokenErr)
	})

	t.Run("should return error when blacklisting token fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			blacklistFunc: func(ctx context.Context, token string, expiresAt time.Time) error {
				return apperrors.ErrTokenStorageErr
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		err := tokenService.RevokeToken(context.Background(), "token")

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrTokenStorageErr)
	})
}

func TestTokenService_RevokeAllUserTokens(t *testing.T) {
	t.Run("should revoke all user tokens successfully", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			deleteFunc: func(ctx context.Context, userID string, tokenType string) error {
				return nil
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		err := tokenService.RevokeAllUserTokens(context.Background(), "user-123")

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		// Arrange
		mockRepo := &MockTokenRepository{
			deleteFunc: func(ctx context.Context, userID string, tokenType string) error {
				return apperrors.ErrTokenStorageErr
			},
		}
		mockJWT := &MockJWTProvider{}

		tokenService := service.NewTokenService(mockRepo, mockJWT, 15*time.Minute, 7*24*time.Hour)

		// Act
		err := tokenService.RevokeAllUserTokens(context.Background(), "user-123")

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrTokenStorageErr)
	})
}
