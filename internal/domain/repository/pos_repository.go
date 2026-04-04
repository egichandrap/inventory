package repository

import (
	"context"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// CartRepository defines the interface for cart data operations
type CartRepository interface {
	// Create creates a new cart
	Create(ctx context.Context, cart *model.Cart) error

	// GetByID retrieves a cart by ID
	GetByID(ctx context.Context, id string) (*model.Cart, error)

	// GetByUserID retrieves active cart by user ID
	GetByUserID(ctx context.Context, userID string) (*model.Cart, error)

	// Update updates an existing cart
	Update(ctx context.Context, cart *model.Cart) error

	// Delete removes a cart
	Delete(ctx context.Context, id string) error

	// ClearItems removes all items from a cart
	ClearItems(ctx context.Context, cartID string) error

	// ListByStatus retrieves carts by status
	ListByStatus(ctx context.Context, status model.CartStatus, limit, offset int) ([]*model.Cart, error)
}

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, transaction *model.Transaction) error

	// GetByID retrieves a transaction by ID
	GetByID(ctx context.Context, id string) (*model.Transaction, error)

	// GetByTransactionNo retrieves a transaction by transaction number
	GetByTransactionNo(ctx context.Context, transactionNo string) (*model.Transaction, error)

	// Update updates an existing transaction
	Update(ctx context.Context, transaction *model.Transaction) error

	// Delete removes a transaction
	Delete(ctx context.Context, id string) error

	// List retrieves transactions with filtering and pagination
	List(ctx context.Context, filter TransactionFilter) ([]*model.Transaction, error)

	// Count returns total number of transactions
	Count(ctx context.Context, filter TransactionFilter) (int64, error)

	// ListWithPagination retrieves transactions with pagination
	ListWithPagination(ctx context.Context, filter TransactionFilter) (*PaginatedTransactions, error)

	// GetByDateRange retrieves transactions within a date range
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.Transaction, error)

	// GetByCashierID retrieves transactions by cashier
	GetByCashierID(ctx context.Context, cashierID string, limit int) ([]*model.Transaction, error)

	// GenerateTransactionNo generates a unique transaction number
	GenerateTransactionNo(ctx context.Context) (string, error)
}

// TransactionFilter defines filter options for listing transactions
type TransactionFilter struct {
	Status      model.TransactionStatus
	PaymentMethod model.PaymentMethod
	CashierID   string
	StartDate   time.Time
	EndDate     time.Time
	Search      string // search by transaction_no or customer_name
	Limit       int
	Offset      int
}

// PaginatedTransactions represents paginated transaction list
type PaginatedTransactions struct {
	Transactions []*model.Transaction `json:"transactions"`
	Total        int64                `json:"total"`
	Limit        int                  `json:"limit"`
	Offset       int                  `json:"offset"`
	TotalPages   int                  `json:"total_pages"`
}
