package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// PostgresCartRepository implements repository.CartRepository with PostgreSQL
type PostgresCartRepository struct {
	db *sql.DB
}

// NewPostgresCartRepository creates a new PostgresCartRepository
func NewPostgresCartRepository(db *sql.DB) *PostgresCartRepository {
	return &PostgresCartRepository{db: db}
}

// Create creates a new cart in PostgreSQL
func (r *PostgresCartRepository) Create(ctx context.Context, cart *model.Cart) error {
	itemsJSON, err := json.Marshal(cart.Items())
	if err != nil {
		return fmt.Errorf("failed to marshal cart items: %w", err)
	}

	query := `
		INSERT INTO carts (id, user_id, customer_name, items, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.ExecContext(ctx, query,
		cart.ID(),
		cart.UserID(),
		cart.CustomerName(),
		itemsJSON,
		cart.Total(),
		cart.CreatedAt(),
		cart.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to create cart: %w", err)
	}

	return nil
}

// GetByID retrieves a cart by ID from PostgreSQL
func (r *PostgresCartRepository) GetByID(ctx context.Context, id string) (*model.Cart, error) {
	query := `
		SELECT id, user_id, customer_name, items, total_amount, created_at, updated_at
		FROM carts
		WHERE id = $1
	`

	var cart model.Cart
	var itemsJSON []byte
	var userID, customerName sql.NullString
	var totalAmount float64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cart,
		&userID,
		&customerName,
		&itemsJSON,
		&totalAmount,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Parse items
	var items []model.CartItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to parse cart items: %w", err)
	}

	// Reconstruct cart
	custName := ""
	if customerName.Valid {
		custName = customerName.String
	}
	
	uid := ""
	if userID.Valid {
		uid = userID.String
	}

	cart = *model.ReconstructCart(
		id,
		uid,
		custName,
		items,
		totalAmount,
		createdAt,
		updatedAt,
	)

	return &cart, nil
}

// GetByUserID retrieves active cart by user ID
func (r *PostgresCartRepository) GetByUserID(ctx context.Context, userID string) (*model.Cart, error) {
	query := `
		SELECT id, user_id, customer_name, items, total_amount, created_at, updated_at
		FROM carts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var cart model.Cart
	var itemsJSON []byte
	var custName, uid sql.NullString
	var totalAmount float64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&cart,
		&uid,
		&custName,
		&itemsJSON,
		&totalAmount,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart by user: %w", err)
	}

	// Parse items
	var items []model.CartItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to parse cart items: %w", err)
	}

	customerName := ""
	if custName.Valid {
		customerName = custName.String
	}

	cartID := ""
	if uid.Valid {
		cartID = uid.String
	}

	cart = *model.ReconstructCart(
		cartID,
		userID,
		customerName,
		items,
		totalAmount,
		createdAt,
		updatedAt,
	)

	return &cart, nil
}

// Update updates an existing cart
func (r *PostgresCartRepository) Update(ctx context.Context, cart *model.Cart) error {
	itemsJSON, err := json.Marshal(cart.Items())
	if err != nil {
		return fmt.Errorf("failed to marshal cart items: %w", err)
	}

	query := `
		UPDATE carts
		SET customer_name = $2, items = $3, total_amount = $4, updated_at = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		cart.ID(),
		cart.CustomerName(),
		itemsJSON,
		cart.Total(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("cart not found")
	}

	return nil
}

// Delete removes a cart
func (r *PostgresCartRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM carts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("cart not found")
	}

	return nil
}

// ClearItems removes all items from a cart
func (r *PostgresCartRepository) ClearItems(ctx context.Context, cartID string) error {
	query := `
		UPDATE carts
		SET items = '[]', total_amount = 0, updated_at = $2
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, cartID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to clear cart items: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check clear result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("cart not found")
	}

	return nil
}

// ListByStatus retrieves carts by status
func (r *PostgresCartRepository) ListByStatus(ctx context.Context, status model.CartStatus, limit, offset int) ([]*model.Cart, error) {
	// TODO: Implement PostgreSQL query
	return nil, fmt.Errorf("ListByStatus not implemented yet")
}

// Ensure interface implementation
var _ repository.CartRepository = (*PostgresCartRepository)(nil)
