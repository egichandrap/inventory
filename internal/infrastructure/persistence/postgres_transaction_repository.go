package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// PostgresTransactionRepository implements repository.TransactionRepository with PostgreSQL
type PostgresTransactionRepository struct {
	db *sql.DB
}

// NewPostgresTransactionRepository creates a new PostgresTransactionRepository
func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

// Create creates a new transaction
func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	itemsJSON, err := json.Marshal(transaction.Items())
	if err != nil {
		return fmt.Errorf("failed to marshal transaction items: %w", err)
	}

	query := `
		INSERT INTO transactions (
			id, transaction_no, cashier_id, cashier_name, customer_name,
			items, subtotal, discount_amount, discount_percent,
			tax_amount, tax_percent, total_amount, payment_method,
			payment_amount, change_amount, status, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err = r.db.ExecContext(ctx, query,
		transaction.ID(),
		transaction.TransactionNo(),
		transaction.CashierID(),
		transaction.CashierName(),
		transaction.CustomerName(),
		itemsJSON,
		transaction.Subtotal(),
		transaction.DiscountAmount(),
		transaction.DiscountPercent(),
		transaction.TaxAmount(),
		transaction.TaxPercent(),
		transaction.TotalAmount(),
		transaction.PaymentMethod(),
		transaction.PaymentAmount(),
		transaction.ChangeAmount(),
		transaction.Status(),
		transaction.Notes(),
		transaction.CreatedAt(),
		transaction.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *PostgresTransactionRepository) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	query := `
		SELECT id, transaction_no, cashier_id, cashier_name, customer_name,
			   items, subtotal, discount_amount, discount_percent,
			   tax_amount, tax_percent, total_amount, payment_method,
			   payment_amount, change_amount, status, notes, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var transaction model.Transaction
	var itemsJSON []byte
	var cashierID sql.NullString
	var customerName, cashierName, transactionNo, notes, statusStr, paymentMethodStr sql.NullString
	var subtotal, discountAmount, discountPercent, taxAmount, taxPercent float64
	var totalAmount, paymentAmount, changeAmount float64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction,
		&transactionNo,
		&cashierID,
		&cashierName,
		&customerName,
		&itemsJSON,
		&subtotal,
		&discountAmount,
		&discountPercent,
		&taxAmount,
		&taxPercent,
		&totalAmount,
		&paymentAmount,
		&changeAmount,
		&statusStr,
		&notes,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Parse items
	var items []model.TransactionItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to parse transaction items: %w", err)
	}

	// Convert strings
	cID := ""
	if cashierID.Valid {
		cID = cashierID.String
	}
	cName := ""
	if cashierName.Valid {
		cName = cashierName.String
	}
	custName := ""
	if customerName.Valid {
		custName = customerName.String
	}
	tNo := ""
	if transactionNo.Valid {
		tNo = transactionNo.String
	}
	n := ""
	if notes.Valid {
		n = notes.String
	}

	paymentMethod := model.PaymentMethod(paymentMethodStr.String)
	status := model.TransactionStatus(statusStr.String)

	transaction = *model.ReconstructTransaction(
		id,
		tNo,
		cID,
		cName,
		custName,
		items,
		subtotal,
		discountAmount,
		discountPercent,
		taxAmount,
		taxPercent,
		totalAmount,
		paymentAmount,
		changeAmount,
		paymentMethod,
		status,
		n,
		createdAt,
		updatedAt,
	)

	return &transaction, nil
}

// GetByTransactionNo retrieves a transaction by transaction number
func (r *PostgresTransactionRepository) GetByTransactionNo(ctx context.Context, transactionNo string) (*model.Transaction, error) {
	query := `
		SELECT id, transaction_no, cashier_id, cashier_name, customer_name,
			   items, subtotal, discount_amount, discount_percent,
			   tax_amount, tax_percent, total_amount, payment_method,
			   payment_amount, change_amount, status, notes, created_at, updated_at
		FROM transactions
		WHERE transaction_no = $1
	`

	var transaction model.Transaction
	var itemsJSON []byte
	var cashierID, customerName, cashierName, notes, statusStr, paymentMethodStr sql.NullString
	var subtotal, discountAmount, discountPercent, taxAmount, taxPercent float64
	var totalAmount, paymentAmount, changeAmount float64
	var createdAt, updatedAt time.Time
	var id string
	var tNo string

	err := r.db.QueryRowContext(ctx, query, transactionNo).Scan(
		&id,
		&tNo,
		&cashierID,
		&cashierName,
		&customerName,
		&itemsJSON,
		&subtotal,
		&discountAmount,
		&discountPercent,
		&taxAmount,
		&taxPercent,
		&totalAmount,
		&paymentAmount,
		&changeAmount,
		&statusStr,
		&notes,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Parse items
	var items []model.TransactionItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to parse transaction items: %w", err)
	}

	cID := ""
	if cashierID.Valid {
		cID = cashierID.String
	}
	cName := ""
	if cashierName.Valid {
		cName = cashierName.String
	}
	custName := ""
	if customerName.Valid {
		custName = customerName.String
	}
	n := ""
	if notes.Valid {
		n = notes.String
	}

	paymentMethod := model.PaymentMethod(paymentMethodStr.String)
	status := model.TransactionStatus(statusStr.String)

	transaction = *model.ReconstructTransaction(
		id,
		tNo,
		cID,
		cName,
		custName,
		items,
		subtotal,
		discountAmount,
		discountPercent,
		taxAmount,
		taxPercent,
		totalAmount,
		paymentAmount,
		changeAmount,
		paymentMethod,
		status,
		n,
		createdAt,
		updatedAt,
	)

	return &transaction, nil
}

// Update updates an existing transaction
func (r *PostgresTransactionRepository) Update(ctx context.Context, transaction *model.Transaction) error {
	itemsJSON, err := json.Marshal(transaction.Items())
	if err != nil {
		return fmt.Errorf("failed to marshal transaction items: %w", err)
	}

	query := `
		UPDATE transactions
		SET customer_name = $2, items = $3, subtotal = $4, discount_amount = $5,
			discount_percent = $6, tax_amount = $7, tax_percent = $8, total_amount = $9,
			payment_method = $10, payment_amount = $11, change_amount = $12,
			status = $13, notes = $14, updated_at = $15
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		transaction.ID(),
		transaction.CustomerName(),
		itemsJSON,
		transaction.Subtotal(),
		transaction.DiscountAmount(),
		transaction.DiscountPercent(),
		transaction.TaxAmount(),
		transaction.TaxPercent(),
		transaction.TotalAmount(),
		transaction.PaymentMethod(),
		transaction.PaymentAmount(),
		transaction.ChangeAmount(),
		transaction.Status(),
		transaction.Notes(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// Delete removes a transaction
func (r *PostgresTransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// List retrieves transactions with filtering
func (r *PostgresTransactionRepository) List(ctx context.Context, filter repository.TransactionFilter) ([]*model.Transaction, error) {
	paginated, err := r.ListWithPagination(ctx, filter)
	if err != nil {
		return nil, err
	}
	return paginated.Transactions, nil
}

// Count returns total number of transactions
func (r *PostgresTransactionRepository) Count(ctx context.Context, filter repository.TransactionFilter) (int64, error) {
	query, args := buildTransactionCountQuery(filter)
	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}
	return count, nil
}

// ListWithPagination retrieves transactions with pagination
func (r *PostgresTransactionRepository) ListWithPagination(ctx context.Context, filter repository.TransactionFilter) (*repository.PaginatedTransactions, error) {
	// Get total count
	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Set defaults
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	query, args := buildTransactionListQuery(filter)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		var itemsJSON []byte
		var cashierID, customerName, cashierName, notes, statusStr, paymentMethodStr sql.NullString
		var subtotal, discountAmount, discountPercent, taxAmount, taxPercent float64
		var totalAmount, paymentAmount, changeAmount float64
		var createdAt, updatedAt time.Time
		var id, transactionNo string

		err := rows.Scan(
			&id,
			&transactionNo,
			&cashierID,
			&cashierName,
			&customerName,
			&itemsJSON,
			&subtotal,
			&discountAmount,
			&discountPercent,
			&taxAmount,
			&taxPercent,
			&totalAmount,
			&paymentAmount,
			&changeAmount,
			&statusStr,
			&notes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Parse items
		var items []model.TransactionItem
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			return nil, fmt.Errorf("failed to parse transaction items: %w", err)
		}

		cID := ""
		if cashierID.Valid {
			cID = cashierID.String
		}
		cName := ""
		if cashierName.Valid {
			cName = cashierName.String
		}
		custName := ""
		if customerName.Valid {
			custName = customerName.String
		}
		n := ""
		if notes.Valid {
			n = notes.String
		}

		paymentMethod := model.PaymentMethod(paymentMethodStr.String)
		status := model.TransactionStatus(statusStr.String)

		transaction = *model.ReconstructTransaction(
			id,
			transactionNo,
			cID,
			cName,
			custName,
			items,
			subtotal,
			discountAmount,
			discountPercent,
			taxAmount,
			taxPercent,
			totalAmount,
			paymentAmount,
			changeAmount,
			paymentMethod,
			status,
			n,
			createdAt,
			updatedAt,
		)

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit != 0 {
		totalPages++
	}

	return &repository.PaginatedTransactions{
		Transactions: transactions,
		Total:        total,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
		TotalPages:   totalPages,
	}, nil
}

// GetByDateRange retrieves transactions within a date range
func (r *PostgresTransactionRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.Transaction, error) {
	query := `
		SELECT id, transaction_no, cashier_id, cashier_name, customer_name,
			   items, subtotal, discount_amount, discount_percent,
			   tax_amount, tax_percent, total_amount, payment_method,
			   payment_amount, change_amount, status, notes, created_at, updated_at
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by date range: %w", err)
	}
	defer rows.Close()

	transactions, err := scanTransactions(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan transactions: %w", err)
	}

	return transactions, nil
}

// GetByCashierID retrieves transactions by cashier
func (r *PostgresTransactionRepository) GetByCashierID(ctx context.Context, cashierID string, limit int) ([]*model.Transaction, error) {
	query := `
		SELECT id, transaction_no, cashier_id, cashier_name, customer_name,
			   items, subtotal, discount_amount, discount_percent,
			   tax_amount, tax_percent, total_amount, payment_method,
			   payment_amount, change_amount, status, notes, created_at, updated_at
		FROM transactions
		WHERE cashier_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, cashierID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by cashier: %w", err)
	}
	defer rows.Close()

	transactions, err := scanTransactions(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan transactions: %w", err)
	}

	return transactions, nil
}

// GenerateTransactionNo generates a unique transaction number
func (r *PostgresTransactionRepository) GenerateTransactionNo(ctx context.Context) (string, error) {
	query := `
		SELECT COALESCE(MAX(CAST(SPLIT_PART(transaction_no, '-', 3) AS INTEGER)), 0) + 1
		FROM transactions
		WHERE transaction_no LIKE $1
	`

	today := time.Now().Format("20060102")
	pattern := fmt.Sprintf("TRX-%s-%%", today)

	var nextNum int
	err := r.db.QueryRowContext(ctx, query, pattern).Scan(&nextNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate transaction number: %w", err)
	}

	return fmt.Sprintf("TRX-%s-%04d", today, nextNum), nil
}

// Helper functions

func buildTransactionCountQuery(filter repository.TransactionFilter) (string, []interface{}) {
	query := `SELECT COUNT(*) FROM transactions WHERE 1=1`
	var args []interface{}
	argNum := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}
	if filter.PaymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argNum)
		args = append(args, filter.PaymentMethod)
		argNum++
	}
	if filter.CashierID != "" {
		query += fmt.Sprintf(" AND cashier_id = $%d", argNum)
		args = append(args, filter.CashierID)
		argNum++
	}
	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argNum)
		args = append(args, filter.StartDate)
		argNum++
	}
	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND created_at < $%d", argNum)
		args = append(args, filter.EndDate)
		argNum++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (transaction_no ILIKE $%d OR customer_name ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}

	return query, args
}

func buildTransactionListQuery(filter repository.TransactionFilter) (string, []interface{}) {
	query := `
		SELECT id, transaction_no, cashier_id, cashier_name, customer_name,
			   items, subtotal, discount_amount, discount_percent,
			   tax_amount, tax_percent, total_amount, payment_method,
			   payment_amount, change_amount, status, notes, created_at, updated_at
		FROM transactions
		WHERE 1=1
	`
	var args []interface{}
	argNum := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}
	if filter.PaymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argNum)
		args = append(args, filter.PaymentMethod)
		argNum++
	}
	if filter.CashierID != "" {
		query += fmt.Sprintf(" AND cashier_id = $%d", argNum)
		args = append(args, filter.CashierID)
		argNum++
	}
	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argNum)
		args = append(args, filter.StartDate)
		argNum++
	}
	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND created_at < $%d", argNum)
		args = append(args, filter.EndDate)
		argNum++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (transaction_no ILIKE $%d OR customer_name ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filter.Limit, filter.Offset)

	return query, args
}

func scanTransactions(rows *sql.Rows) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		var itemsJSON []byte
		var cashierID, customerName, cashierName, notes, statusStr, paymentMethodStr sql.NullString
		var subtotal, discountAmount, discountPercent, taxAmount, taxPercent float64
		var totalAmount, paymentAmount, changeAmount float64
		var createdAt, updatedAt time.Time
		var id, transactionNo string

		err := rows.Scan(
			&id,
			&transactionNo,
			&cashierID,
			&cashierName,
			&customerName,
			&itemsJSON,
			&subtotal,
			&discountAmount,
			&discountPercent,
			&taxAmount,
			&taxPercent,
			&totalAmount,
			&paymentAmount,
			&changeAmount,
			&statusStr,
			&notes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Parse items
		var items []model.TransactionItem
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			return nil, fmt.Errorf("failed to parse transaction items: %w", err)
		}

		cID := ""
		if cashierID.Valid {
			cID = cashierID.String
		}
		cName := ""
		if cashierName.Valid {
			cName = cashierName.String
		}
		custName := ""
		if customerName.Valid {
			custName = customerName.String
		}
		n := ""
		if notes.Valid {
			n = notes.String
		}

		paymentMethod := model.PaymentMethod(paymentMethodStr.String)
		status := model.TransactionStatus(statusStr.String)

		transaction = *model.ReconstructTransaction(
			id,
			transactionNo,
			cID,
			cName,
			custName,
			items,
			subtotal,
			discountAmount,
			discountPercent,
			taxAmount,
			taxPercent,
			totalAmount,
			paymentAmount,
			changeAmount,
			paymentMethod,
			status,
			n,
			createdAt,
			updatedAt,
		)

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// Ensure interface implementation
var _ repository.TransactionRepository = (*PostgresTransactionRepository)(nil)
