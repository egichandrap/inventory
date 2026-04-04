package repository

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// TableRepository defines the interface for table data operations
type TableRepository interface {
	// Create creates a new table
	Create(ctx context.Context, table *model.Table) error

	// GetByID retrieves a table by ID
	GetByID(ctx context.Context, id string) (*model.Table, error)

	// GetByNumber retrieves a table by number
	GetByNumber(ctx context.Context, number int) (*model.Table, error)

	// Update updates an existing table
	Update(ctx context.Context, table *model.Table) error

	// Delete removes a table
	Delete(ctx context.Context, id string) error

	// List retrieves all tables with optional filtering
	List(ctx context.Context, filter *TableFilter) ([]*model.Table, error)

	// Count returns total number of tables
	Count(ctx context.Context, filter *TableFilter) (int64, error)

	// GetAvailableTables retrieves all available tables
	GetAvailableTables(ctx context.Context, location *model.TableLocation) ([]*model.Table, error)

	// ExistsByNumber checks if a table with the given number exists
	ExistsByNumber(ctx context.Context, number int, excludeID string) (bool, error)
}

// TableFilter defines filter options for listing tables
type TableFilter struct {
	Location *model.TableLocation
	Status   *model.TableStatus
	MinCapacity int
	MaxCapacity int
	Search   string // search by description
	Limit    int
	Offset   int
}
