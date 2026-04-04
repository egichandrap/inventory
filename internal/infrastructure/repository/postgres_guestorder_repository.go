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

// PostgresGuestOrderRepository implements repository.GuestOrderRepository
type PostgresGuestOrderRepository struct {
	db *sql.DB
}

// NewPostgresGuestOrderRepository creates the repository
func NewPostgresGuestOrderRepository(db *sql.DB) repository.GuestOrderRepository {
	return &PostgresGuestOrderRepository{db: db}
}

func (r *PostgresGuestOrderRepository) Create(ctx context.Context, order *model.GuestOrder) error {
	itemsJSON, err := json.Marshal(order.Items())
	if err != nil {
		return fmt.Errorf("failed to marshal order items: %w", err)
	}

	query := `
		INSERT INTO guest_orders (
			id, order_number, table_id, table_number, customer_name, customer_phone,
			items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			total_amount, payment_method, payment_status, payment_amount, change_amount,
			status, notes, session_id, created_at, updated_at, completed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)
	`

	var completedAt interface{} = nil
	if order.CompletedAt() != nil {
		completedAt = *order.CompletedAt()
	}

	_, err = r.db.ExecContext(ctx, query,
		order.ID(),
		order.OrderNumber(),
		order.TableID(),
		order.TableNumber(),
		order.CustomerName(),
		order.CustomerPhone(),
		itemsJSON,
		order.Subtotal(),
		order.TaxAmount(),
		order.TaxPercent(),
		order.DiscountAmount(),
		order.DiscountPercent(),
		order.TotalAmount(),
		order.PaymentMethod(),
		order.PaymentStatus(),
		order.PaymentAmount(),
		order.ChangeAmount(),
		order.Status(),
		order.Notes(),
		order.SessionID(),
		order.CreatedAt(),
		order.UpdatedAt(),
		completedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create guest order: %w", err)
	}
	return nil
}

func (r *PostgresGuestOrderRepository) GetByID(ctx context.Context, id string) (*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanOrder(row)
}

func (r *PostgresGuestOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE order_number = $1
	`

	row := r.db.QueryRowContext(ctx, query, orderNumber)
	return r.scanOrder(row)
}

func (r *PostgresGuestOrderRepository) Update(ctx context.Context, order *model.GuestOrder) error {
	itemsJSON, err := json.Marshal(order.Items())
	if err != nil {
		return fmt.Errorf("failed to marshal order items: %w", err)
	}

	query := `
		UPDATE guest_orders
		SET table_id = $1, table_number = $2, customer_name = $3, customer_phone = $4,
			items = $5, subtotal = $6, tax_amount = $7, tax_percent = $8,
			discount_amount = $9, discount_percent = $10, total_amount = $11,
			payment_method = $12, payment_status = $13, payment_amount = $14, change_amount = $15,
			status = $16, notes = $17, session_id = $18, updated_at = $19, completed_at = $20
		WHERE id = $21
	`

	var completedAt interface{} = nil
	if order.CompletedAt() != nil {
		completedAt = *order.CompletedAt()
	}

	_, err = r.db.ExecContext(ctx, query,
		order.TableID(),
		order.TableNumber(),
		order.CustomerName(),
		order.CustomerPhone(),
		itemsJSON,
		order.Subtotal(),
		order.TaxAmount(),
		order.TaxPercent(),
		order.DiscountAmount(),
		order.DiscountPercent(),
		order.TotalAmount(),
		order.PaymentMethod(),
		order.PaymentStatus(),
		order.PaymentAmount(),
		order.ChangeAmount(),
		order.Status(),
		order.Notes(),
		order.SessionID(),
		order.UpdatedAt(),
		completedAt,
		order.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update guest order: %w", err)
	}
	return nil
}

func (r *PostgresGuestOrderRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM guest_orders WHERE id = $1", id)
	return err
}

func (r *PostgresGuestOrderRepository) List(ctx context.Context, filter repository.GuestOrderFilter) ([]*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, string(filter.Status))
		argCount++
	}
	if filter.PaymentStatus != "" {
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, string(filter.PaymentStatus))
		argCount++
	}
	if filter.TableID != "" {
		query += fmt.Sprintf(" AND table_id = $%d", argCount)
		args = append(args, filter.TableID)
		argCount++
	}
	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	}
	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (order_number ILIKE $%d OR customer_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filter.Search+"%")
		argCount++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, filter.Offset)
		}
	}

	return r.scanOrders(r.db.QueryContext(ctx, query, args...))
}

func (r *PostgresGuestOrderRepository) Count(ctx context.Context, filter repository.GuestOrderFilter) (int64, error) {
	query := "SELECT COUNT(*) FROM guest_orders WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, string(filter.Status))
		argCount++
	}
	if filter.TableID != "" {
		query += fmt.Sprintf(" AND table_id = $%d", argCount)
		args = append(args, filter.TableID)
		argCount++
	}
	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	}
	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *PostgresGuestOrderRepository) ListWithPagination(ctx context.Context, filter repository.GuestOrderFilter) (*repository.PaginatedGuestOrders, error) {
	orders, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.Limit
	if filter.Limit > 0 && int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &repository.PaginatedGuestOrders{
		Orders:       orders,
		Total:        total,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
		TotalPages:   totalPages,
	}, nil
}

func (r *PostgresGuestOrderRepository) GetByTableID(ctx context.Context, tableID string, limit int) ([]*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE table_id = $1
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	return r.scanOrders(r.db.QueryContext(ctx, query, tableID))
}

func (r *PostgresGuestOrderRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
	`

	return r.scanOrders(r.db.QueryContext(ctx, query, startDate, endDate))
}

func (r *PostgresGuestOrderRepository) GetByStatus(ctx context.Context, status model.GuestOrderStatus, limit int) ([]*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE status = $1
		ORDER BY created_at ASC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	return r.scanOrders(r.db.QueryContext(ctx, query, string(status)))
}

func (r *PostgresGuestOrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	now := time.Now()
	base := fmt.Sprintf("ORD-%s", now.Format("20060102"))

	var maxNum int
	err := r.db.QueryRowContext(ctx,
		"SELECT MAX(CAST(SUBSTRING(order_number FROM LENGTH($1) + 2) AS INTEGER)) FROM guest_orders WHERE order_number LIKE $1 || '%'",
		base,
	).Scan(&maxNum)

	if err != nil || maxNum == 0 {
		maxNum = 0
	}

	return fmt.Sprintf("%s-%04d", base, maxNum+1), nil
}

func (r *PostgresGuestOrderRepository) GetPendingOrders(ctx context.Context) ([]*model.GuestOrder, error) {
	return r.GetByStatus(ctx, model.OrderPending, 0)
}

func (r *PostgresGuestOrderRepository) GetActiveOrders(ctx context.Context) ([]*model.GuestOrder, error) {
	query := `
		SELECT id, order_number, table_id, table_number, customer_name, customer_phone,
			   items, subtotal, tax_amount, tax_percent, discount_amount, discount_percent,
			   total_amount, payment_method, payment_status, payment_amount, change_amount,
			   status, notes, session_id, created_at, updated_at, completed_at
		FROM guest_orders
		WHERE status NOT IN ('SERVED', 'CANCELLED')
		ORDER BY created_at ASC
	`

	return r.scanOrders(r.db.QueryContext(ctx, query))
}

// Helper methods
func (r *PostgresGuestOrderRepository) scanOrder(row *sql.Row) (*model.GuestOrder, error) {
	var id, orderNumber, tableID, customerName, customerPhone, sessionID string
	var tableNumber int
	var itemsJSON []byte
	var subtotal, taxAmount, taxPercent, discountAmount, discountPercent, totalAmount float64
	var paymentMethod, paymentStatus, status string
	var paymentAmount, changeAmount float64
	var notes string
	var createdAt, updatedAt time.Time
	var completedAt sql.NullTime

	err := row.Scan(
		&id, &orderNumber, &tableID, &tableNumber, &customerName, &customerPhone,
		&itemsJSON, &subtotal, &taxAmount, &taxPercent, &discountAmount, &discountPercent,
		&totalAmount, &paymentMethod, &paymentStatus, &paymentAmount, &changeAmount,
		&status, &notes, &sessionID, &createdAt, &updatedAt, &completedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan guest order: %w", err)
	}

	var items []model.GuestOrderItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to parse order items: %w", err)
	}

	reconstructedItems := make([]model.GuestOrderItem, len(items))
	for i, item := range items {
		reconstructedItems[i] = *model.ReconstructGuestOrderItem(
			item.ProductID(), item.ProductName(), item.Quantity(),
			item.UnitPrice(), item.Subtotal(), item.Notes(),
		)
	}

	var completedAtPtr *time.Time
	if completedAt.Valid {
		completedAtPtr = &completedAt.Time
	}

	return model.ReconstructGuestOrder(
		id, orderNumber, tableID, customerName, customerPhone, sessionID,
		tableNumber, reconstructedItems,
		subtotal, taxAmount, taxPercent, discountAmount, discountPercent, totalAmount,
		model.PaymentMethod(paymentMethod),
		model.GuestOrderPaymentStatus(paymentStatus),
		paymentAmount, changeAmount,
		model.GuestOrderStatus(status),
		notes,
		createdAt, updatedAt,
		completedAtPtr,
	), nil
}

func (r *PostgresGuestOrderRepository) scanOrders(rows *sql.Rows, err error) ([]*model.GuestOrder, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.GuestOrder
	for rows.Next() {
		var id, orderNumber, tableID, customerName, customerPhone, sessionID string
		var tableNumber int
		var itemsJSON []byte
		var subtotal, taxAmount, taxPercent, discountAmount, discountPercent, totalAmount float64
		var paymentMethod, paymentStatus, status string
		var paymentAmount, changeAmount float64
		var notes string
		var createdAt, updatedAt time.Time
		var completedAt sql.NullTime

		if err := rows.Scan(
			&id, &orderNumber, &tableID, &tableNumber, &customerName, &customerPhone,
			&itemsJSON, &subtotal, &taxAmount, &taxPercent, &discountAmount, &discountPercent,
			&totalAmount, &paymentMethod, &paymentStatus, &paymentAmount, &changeAmount,
			&status, &notes, &sessionID, &createdAt, &updatedAt, &completedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan guest order: %w", err)
		}

		var items []model.GuestOrderItem
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			return nil, fmt.Errorf("failed to parse order items: %w", err)
		}

		reconstructedItems := make([]model.GuestOrderItem, len(items))
		for i, item := range items {
			reconstructedItems[i] = *model.ReconstructGuestOrderItem(
				item.ProductID(), item.ProductName(), item.Quantity(),
				item.UnitPrice(), item.Subtotal(), item.Notes(),
			)
		}

		var completedAtPtr *time.Time
		if completedAt.Valid {
			completedAtPtr = &completedAt.Time
		}

		orders = append(orders, model.ReconstructGuestOrder(
			id, orderNumber, tableID, customerName, customerPhone, sessionID,
			tableNumber, reconstructedItems,
			subtotal, taxAmount, taxPercent, discountAmount, discountPercent, totalAmount,
			model.PaymentMethod(paymentMethod),
			model.GuestOrderPaymentStatus(paymentStatus),
			paymentAmount, changeAmount,
			model.GuestOrderStatus(status),
			notes,
			createdAt, updatedAt,
			completedAtPtr,
		))
	}

	return orders, rows.Err()
}

var _ repository.GuestOrderRepository = (*PostgresGuestOrderRepository)(nil)
