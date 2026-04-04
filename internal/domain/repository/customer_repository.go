package repository

import (
	"context"
	
	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	// Create creates a new customer
	Create(ctx context.Context, customer *model.Customer) error

	// GetByID retrieves a customer by ID
	GetByID(ctx context.Context, id string) (*model.Customer, error)

	// GetByEmail retrieves a customer by email
	GetByEmail(ctx context.Context, email string) (*model.Customer, error)

	// Update updates an existing customer
	Update(ctx context.Context, customer *model.Customer) error

	// Delete removes a customer
	Delete(ctx context.Context, id string) error

	// List retrieves customers with pagination
	List(ctx context.Context, filter CustomerFilter) (*PaginatedCustomers, error)

	// Count returns total number of customers
	Count(ctx context.Context, filter CustomerFilter) (int64, error)
}

// CustomerFilter defines filter options for listing customers
type CustomerFilter struct {
	Search string
	Limit  int
	Offset int
}

// PaginatedCustomers represents paginated customer list
type PaginatedCustomers struct {
	Customers []*model.Customer `json:"customers"`
	Total     int64             `json:"total"`
	Limit     int               `json:"limit"`
	Offset    int               `json:"offset"`
	TotalPages int              `json:"total_pages"`
}
