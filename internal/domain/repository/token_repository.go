package repository

import (
	"context"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// TokenRepository defines the interface for token persistence operations
type TokenRepository interface {
	// Store saves a token with its metadata
	Store(ctx context.Context, userID string, token string, tokenType string, expiresAt time.Time) error

	// Find retrieves a token by user ID and type
	Find(ctx context.Context, userID string, tokenType string) (*model.Token, error)

	// Delete removes a token
	Delete(ctx context.Context, userID string, tokenType string) error

	// IsBlacklisted checks if a token is blacklisted
	IsBlacklisted(ctx context.Context, token string) (bool, error)

	// Blacklist adds a token to the blacklist
	Blacklist(ctx context.Context, token string, expiresAt time.Time) error
}
