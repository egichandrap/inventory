package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// PostgresInventoryRepository implements InventoryRepository using PostgreSQL
type PostgresInventoryRepository struct {
	db *sql.DB
}

// NewPostgresInventoryRepository creates a new PostgreSQL inventory repository
func NewPostgresInventoryRepository(db *sql.DB) repository.InventoryRepository {
	return &PostgresInventoryRepository{
		db: db,
	}
}

// Create creates a new inventory item
func (r *PostgresInventoryRepository) Create(ctx context.Context, inventory *model.Inventory) error {
	query := `
		INSERT INTO inventories (id, sku, name, description, quantity, unit, location, min_stock, max_stock, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	inventory.CreatedAt = now
	inventory.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		inventory.ID,
		inventory.SKU,
		inventory.Name,
		inventory.Description,
		inventory.Quantity,
		inventory.Unit,
		inventory.Location,
		inventory.MinStock,
		inventory.MaxStock,
		inventory.Price,
		inventory.CreatedAt,
		inventory.UpdatedAt,
	)

	return err
}

// GetByID retrieves an inventory item by its ID
func (r *PostgresInventoryRepository) GetByID(ctx context.Context, id string) (*model.Inventory, error) {
	query := `
		SELECT id, sku, name, description, quantity, unit, location, min_stock, max_stock, price, created_at, updated_at
		FROM inventories
		WHERE id = $1
	`

	inv := &model.Inventory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID,
		&inv.SKU,
		&inv.Name,
		&inv.Description,
		&inv.Quantity,
		&inv.Unit,
		&inv.Location,
		&inv.MinStock,
		&inv.MaxStock,
		&inv.Price,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return inv, nil
}

// GetBySKU retrieves an inventory item by its SKU
func (r *PostgresInventoryRepository) GetBySKU(ctx context.Context, sku string) (*model.Inventory, error) {
	query := `
		SELECT id, sku, name, description, quantity, unit, location, min_stock, max_stock, price, created_at, updated_at
		FROM inventories
		WHERE sku = $1
	`

	inv := &model.Inventory{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&inv.ID,
		&inv.SKU,
		&inv.Name,
		&inv.Description,
		&inv.Quantity,
		&inv.Unit,
		&inv.Location,
		&inv.MinStock,
		&inv.MaxStock,
		&inv.Price,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return inv, nil
}

// Update updates an existing inventory item
func (r *PostgresInventoryRepository) Update(ctx context.Context, inventory *model.Inventory) error {
	query := `
		UPDATE inventories
		SET sku = $1, name = $2, description = $3, quantity = $4, unit = $5, location = $6, min_stock = $7, max_stock = $8, price = $9, updated_at = $10
		WHERE id = $11
	`

	inventory.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		inventory.SKU,
		inventory.Name,
		inventory.Description,
		inventory.Quantity,
		inventory.Unit,
		inventory.Location,
		inventory.MinStock,
		inventory.MaxStock,
		inventory.Price,
		inventory.UpdatedAt,
		inventory.ID,
	)

	return err
}

// Delete removes an inventory item
func (r *PostgresInventoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM inventories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves a list of inventory items with optional filtering
func (r *PostgresInventoryRepository) List(ctx context.Context, filter *model.InventoryFilter) ([]*model.Inventory, error) {
	query := `
		SELECT id, sku, name, description, quantity, unit, location, min_stock, max_stock, price, created_at, updated_at
		FROM inventories
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.SKU != nil {
			query += " AND sku LIKE $" + string(rune('0'+argCount))
			args = append(args, "%"+*filter.SKU+"%")
			argCount++
		}
		if filter.Name != nil {
			query += " AND name LIKE $" + string(rune('0'+argCount))
			args = append(args, "%"+*filter.Name+"%")
			argCount++
		}
		if filter.Location != nil {
			query += " AND location = $" + string(rune('0'+argCount))
			args = append(args, *filter.Location)
			argCount++
		}
		if filter.MinQty != nil {
			query += " AND quantity >= $" + string(rune('0'+argCount))
			args = append(args, *filter.MinQty)
			argCount++
		}
		if filter.MaxQty != nil {
			query += " AND quantity <= $" + string(rune('0'+argCount))
			args = append(args, *filter.MaxQty)
			argCount++
		}
	}

	query += " ORDER BY created_at DESC"

	if filter != nil && filter.Limit > 0 {
		query += " LIMIT $" + string(rune('0'+argCount))
		args = append(args, filter.Limit)
		argCount++

		if filter.Offset > 0 {
			query += " OFFSET $" + string(rune('0'+argCount))
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []*model.Inventory
	for rows.Next() {
		inv := &model.Inventory{}
		err := rows.Scan(
			&inv.ID,
			&inv.SKU,
			&inv.Name,
			&inv.Description,
			&inv.Quantity,
			&inv.Unit,
			&inv.Location,
			&inv.MinStock,
			&inv.MaxStock,
			&inv.Price,
			&inv.CreatedAt,
			&inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}

	return inventories, rows.Err()
}

// Count returns the total count of inventory items
func (r *PostgresInventoryRepository) Count(ctx context.Context, filter *model.InventoryFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM inventories WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.SKU != nil {
			query += " AND sku LIKE $" + string(rune('0'+argCount))
			args = append(args, "%"+*filter.SKU+"%")
			argCount++
		}
		if filter.Name != nil {
			query += " AND name LIKE $" + string(rune('0'+argCount))
			args = append(args, "%"+*filter.Name+"%")
			argCount++
		}
		if filter.Location != nil {
			query += " AND location = $" + string(rune('0'+argCount))
			args = append(args, *filter.Location)
			argCount++
		}
		if filter.MinQty != nil {
			query += " AND quantity >= $" + string(rune('0'+argCount))
			args = append(args, *filter.MinQty)
			argCount++
		}
		if filter.MaxQty != nil {
			query += " AND quantity <= $" + string(rune('0'+argCount))
			args = append(args, *filter.MaxQty)
			argCount++
		}
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// UpdateQuantity updates the quantity of an inventory item
func (r *PostgresInventoryRepository) UpdateQuantity(ctx context.Context, id string, quantity int) error {
	query := `UPDATE inventories SET quantity = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, quantity, time.Now(), id)
	return err
}

// ExistsBySKU checks if an inventory item with the given SKU exists
func (r *PostgresInventoryRepository) ExistsBySKU(ctx context.Context, sku string, excludeID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM inventories WHERE sku = $1`
	args := []interface{}{sku}

	if excludeID != "" {
		query += " AND id != $2"
		args = append(args, excludeID)
	}
	query += ")"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	return exists, err
}

// Helper function to build parameterized query placeholders for PostgreSQL
func buildPlaceholders(start int, count int) []string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "$" + string(rune('0'+start+i))
	}
	return placeholders
}

// MemoryInventoryRepository is an in-memory implementation for testing
type MemoryInventoryRepository struct {
	items map[string]*model.Inventory
}

// NewMemoryInventoryRepository creates a new in-memory inventory repository
func NewMemoryInventoryRepository() repository.InventoryRepository {
	return &MemoryInventoryRepository{
		items: make(map[string]*model.Inventory),
	}
}

// Create creates a new inventory item
func (r *MemoryInventoryRepository) Create(ctx context.Context, inventory *model.Inventory) error {
	now := time.Now()
	inventory.CreatedAt = now
	inventory.UpdatedAt = now
	r.items[inventory.ID] = inventory
	return nil
}

// GetByID retrieves an inventory item by its ID
func (r *MemoryInventoryRepository) GetByID(ctx context.Context, id string) (*model.Inventory, error) {
	if inv, ok := r.items[id]; ok {
		return inv, nil
	}
	return nil, nil
}

// GetBySKU retrieves an inventory item by its SKU
func (r *MemoryInventoryRepository) GetBySKU(ctx context.Context, sku string) (*model.Inventory, error) {
	for _, inv := range r.items {
		if inv.SKU == sku {
			return inv, nil
		}
	}
	return nil, nil
}

// Update updates an existing inventory item
func (r *MemoryInventoryRepository) Update(ctx context.Context, inventory *model.Inventory) error {
	inventory.UpdatedAt = time.Now()
	r.items[inventory.ID] = inventory
	return nil
}

// Delete removes an inventory item
func (r *MemoryInventoryRepository) Delete(ctx context.Context, id string) error {
	delete(r.items, id)
	return nil
}

// List retrieves a list of inventory items with optional filtering
func (r *MemoryInventoryRepository) List(ctx context.Context, filter *model.InventoryFilter) ([]*model.Inventory, error) {
	var result []*model.Inventory

	for _, inv := range r.items {
		if filter != nil {
			if filter.SKU != nil && !strings.Contains(inv.SKU, *filter.SKU) {
				continue
			}
			if filter.Name != nil && !strings.Contains(inv.Name, *filter.Name) {
				continue
			}
			if filter.Location != nil && inv.Location != *filter.Location {
				continue
			}
			if filter.MinQty != nil && inv.Quantity < *filter.MinQty {
				continue
			}
			if filter.MaxQty != nil && inv.Quantity > *filter.MaxQty {
				continue
			}
		}
		result = append(result, inv)
	}

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		start := filter.Offset
		end := start + filter.Limit

		if start >= len(result) {
			return []*model.Inventory{}, nil
		}
		if end > len(result) {
			end = len(result)
		}

		result = result[start:end]
	}

	return result, nil
}

// Count returns the total count of inventory items
func (r *MemoryInventoryRepository) Count(ctx context.Context, filter *model.InventoryFilter) (int64, error) {
	list, err := r.List(ctx, filter)
	return int64(len(list)), err
}

// UpdateQuantity updates the quantity of an inventory item
func (r *MemoryInventoryRepository) UpdateQuantity(ctx context.Context, id string, quantity int) error {
	if inv, ok := r.items[id]; ok {
		inv.Quantity = quantity
		inv.UpdatedAt = time.Now()
	}
	return nil
}

// ExistsBySKU checks if an inventory item with the given SKU exists
func (r *MemoryInventoryRepository) ExistsBySKU(ctx context.Context, sku string, excludeID string) (bool, error) {
	for id, inv := range r.items {
		if inv.SKU == sku && id != excludeID {
			return true, nil
		}
	}
	return false, nil
}
