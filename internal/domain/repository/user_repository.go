package repository

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByUsername retrieves a user by their username
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *model.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error

	// Delete removes a user
	Delete(ctx context.Context, id string) error
}
