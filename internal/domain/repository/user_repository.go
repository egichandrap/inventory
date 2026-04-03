package repository

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// UserFilter defines filter options for listing users
type UserFilter struct {
	Role   model.UserRole
	Status model.UserStatus
	Search string // search by username, email, or full_name
	Limit  int
	Offset int
}

// PaginatedUsers represents paginated user list
type PaginatedUsers struct {
	Users      []*model.User `json:"users"`
	Total      int64         `json:"total"`
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	TotalPages int           `json:"total_pages"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByUsername retrieves a user by their username
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// FindByEmail retrieves a user by their email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *model.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error

	// Delete removes a user
	Delete(ctx context.Context, id string) error

	// List retrieves users with filtering and pagination
	List(ctx context.Context, filter UserFilter) ([]*model.User, error)

	// Count returns total number of users with optional filter
	Count(ctx context.Context, filter UserFilter) (int64, error)

	// ListWithPagination retrieves users with pagination
	ListWithPagination(ctx context.Context, filter UserFilter) (*PaginatedUsers, error)

	// ExistsByUsername checks if a username already exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail checks if an email already exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, id, hashedPassword string) error

	// UpdateLastLogin updates the last login timestamp
	UpdateLastLogin(ctx context.Context, id string) error
}
