package repository

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	// Create creates a new category
	Create(ctx context.Context, category *model.Category) error

	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id string) (*model.Category, error)

	// GetBySlug retrieves a category by slug
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)

	// Update updates an existing category
	Update(ctx context.Context, category *model.Category) error

	// Delete removes a category
	Delete(ctx context.Context, id string) error

	// List retrieves categories with filtering
	List(ctx context.Context, filter CategoryFilter) ([]*model.Category, error)

	// ListByParentID retrieves categories by parent ID
	ListByParentID(ctx context.Context, parentID string) ([]*model.Category, error)

	// Count returns total number of categories
	Count(ctx context.Context, filter CategoryFilter) (int64, error)
}

// CategoryFilter defines filter options for listing categories
type CategoryFilter struct {
	IsActive *bool
	Level    *int
	Search   string
	Limit    int
	Offset   int
}
