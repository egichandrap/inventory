package repository

import (
	"context"
	"sync"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// MemoryTokenRepository is an in-memory implementation of TokenRepository
type MemoryTokenRepository struct {
	mu         sync.RWMutex
	tokens     map[string]*storedToken
	blacklist  map[string]time.Time
}

type storedToken struct {
	userID    string
	token     string
	tokenType string
	expiresAt time.Time
}

// NewMemoryTokenRepository creates a new in-memory token repository
func NewMemoryTokenRepository() *MemoryTokenRepository {
	return &MemoryTokenRepository{
		tokens:    make(map[string]*storedToken),
		blacklist: make(map[string]time.Time),
	}
}

// Store saves a token with its metadata
func (r *MemoryTokenRepository) Store(ctx context.Context, userID string, token string, tokenType string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.makeKey(userID, tokenType)
	r.tokens[key] = &storedToken{
		userID:    userID,
		token:     token,
		tokenType: tokenType,
		expiresAt: expiresAt,
	}

	return nil
}

// Find retrieves a token by user ID and type
func (r *MemoryTokenRepository) Find(ctx context.Context, userID string, tokenType string) (*model.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.makeKey(userID, tokenType)
	stored, exists := r.tokens[key]
	if !exists {
		return nil, nil
	}

	if time.Now().After(stored.expiresAt) {
		delete(r.tokens, key)
		return nil, nil
	}

	return &model.Token{
		AccessToken: stored.token,
		ExpiresAt:   stored.expiresAt,
	}, nil
}

// Delete removes a token
func (r *MemoryTokenRepository) Delete(ctx context.Context, userID string, tokenType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.makeKey(userID, tokenType)
	delete(r.tokens, key)

	return nil
}

// IsBlacklisted checks if a token is blacklisted
func (r *MemoryTokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expiresAt, exists := r.blacklist[token]
	if !exists {
		return false, nil
	}

	// Remove expired blacklist entries
	if time.Now().After(expiresAt) {
		delete(r.blacklist, token)
		return false, nil
	}

	return true, nil
}

// Blacklist adds a token to the blacklist
func (r *MemoryTokenRepository) Blacklist(ctx context.Context, token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.blacklist[token] = expiresAt

	return nil
}

func (r *MemoryTokenRepository) makeKey(userID, tokenType string) string {
	return userID + ":" + tokenType
}
