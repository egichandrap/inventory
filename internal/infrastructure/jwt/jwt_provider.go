package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// Config holds JWT configuration
type Config struct {
	SecretKey string
	Issuer    string
	Audience  string
	Algorithm string
}

// Provider implements JWT token generation and validation
type Provider struct {
	config Config
}

// CustomClaims extends JWT claims with custom fields
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// NewProvider creates a new JWT provider
func NewProvider(config Config) *Provider {
	return &Provider{
		config: config,
	}
}

// GenerateToken generates a new JWT token with the given claims
func (p *Provider) GenerateToken(claims *model.TokenClaims, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    p.config.Issuer,
			Subject:   claims.UserID,
		},
	})

	tokenString, err := token.SignedString([]byte(p.config.SecretKey))
	if err != nil {
		return "", apperrors.Wrap(err, apperrors.ErrTokenGeneration, "Failed to sign JWT token", apperrors.ErrTokenGenerationErr.GetHTTPStatus())
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns its claims
func (p *Provider) ValidateToken(tokenString string) (*model.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrInvalidTokenErr
		}
		return []byte(p.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperrors.ErrExpiredTokenErr
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, apperrors.ErrInvalidTokenErr
		}
		return nil, apperrors.Wrap(err, apperrors.ErrInvalidToken, "Failed to parse JWT token", apperrors.ErrInvalidTokenErr.GetHTTPStatus())
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidTokenErr
	}

	return &model.TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

// GetExpiration extracts the expiration time from a token
func (p *Provider) GetExpiration(tokenString string) (time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.config.SecretKey), nil
	})

	if err != nil {
		return time.Time{}, apperrors.Wrap(err, apperrors.ErrInvalidToken, "Failed to parse JWT token", apperrors.ErrInvalidTokenErr.GetHTTPStatus())
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return time.Time{}, apperrors.ErrInvalidTokenErr
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, apperrors.New(apperrors.ErrInvalidToken, "Token has no expiration", apperrors.ErrInvalidTokenErr.GetHTTPStatus())
	}

	return claims.ExpiresAt.Time, nil
}
