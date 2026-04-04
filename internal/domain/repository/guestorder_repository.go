package repository

import (
	"context"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// GuestOrderRepository defines the interface for guest order data operations
type GuestOrderRepository interface {
	// Create creates a new guest order
	Create(ctx context.Context, order *model.GuestOrder) error

	// GetByID retrieves a guest order by ID
	GetByID(ctx context.Context, id string) (*model.GuestOrder, error)

	// GetByOrderNumber retrieves a guest order by order number
	GetByOrderNumber(ctx context.Context, orderNumber string) (*model.GuestOrder, error)

	// Update updates an existing guest order
	Update(ctx context.Context, order *model.GuestOrder) error

	// Delete removes a guest order
	Delete(ctx context.Context, id string) error

	// List retrieves guest orders with filtering and pagination
	List(ctx context.Context, filter GuestOrderFilter) ([]*model.GuestOrder, error)

	// Count returns total number of guest orders
	Count(ctx context.Context, filter GuestOrderFilter) (int64, error)

	// ListWithPagination retrieves guest orders with pagination
	ListWithPagination(ctx context.Context, filter GuestOrderFilter) (*PaginatedGuestOrders, error)

	// GetByTableID retrieves orders for a specific table
	GetByTableID(ctx context.Context, tableID string, limit int) ([]*model.GuestOrder, error)

	// GetByDateRange retrieves orders within a date range
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.GuestOrder, error)

	// GetByStatus retrieves orders by status
	GetByStatus(ctx context.Context, status model.GuestOrderStatus, limit int) ([]*model.GuestOrder, error)

	// GenerateOrderNumber generates a unique order number
	GenerateOrderNumber(ctx context.Context) (string, error)

	// GetPendingOrders retrieves all pending orders
	GetPendingOrders(ctx context.Context) ([]*model.GuestOrder, error)

	// GetActiveOrders retrieves all active orders (not completed)
	GetActiveOrders(ctx context.Context) ([]*model.GuestOrder, error)
}

// GuestOrderFilter defines filter options for listing guest orders
type GuestOrderFilter struct {
	Status        model.GuestOrderStatus
	PaymentStatus model.GuestOrderPaymentStatus
	TableID       string
	StartDate     time.Time
	EndDate       time.Time
	PaymentMethod model.PaymentMethod
	Search        string // search by order_number or customer_name
	Limit         int
	Offset        int
}

// PaginatedGuestOrders represents paginated guest order list
type PaginatedGuestOrders struct {
	Orders       []*model.GuestOrder `json:"orders"`
	Total        int64               `json:"total"`
	Limit        int                 `json:"limit"`
	Offset       int                 `json:"offset"`
	TotalPages   int                 `json:"total_pages"`
}
