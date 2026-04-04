package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// PostgresTableRepository implements repository.TableRepository with PostgreSQL
type PostgresTableRepository struct {
	db *sql.DB
}

// NewPostgresTableRepository creates a new PostgresTableRepository
func NewPostgresTableRepository(db *sql.DB) repository.TableRepository {
	return &PostgresTableRepository{db: db}
}

func (r *PostgresTableRepository) Create(ctx context.Context, table *model.Table) error {
	query := `
		INSERT INTO tables (id, number, location, capacity, status, qr_code, qr_generated, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		table.ID(),
		table.Number(),
		table.Location(),
		table.Capacity(),
		table.Status(),
		table.QRCode(),
		table.IsQRGenerated(),
		table.Description(),
		table.CreatedAt(),
		table.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (r *PostgresTableRepository) GetByID(ctx context.Context, id string) (*model.Table, error) {
	query := `
		SELECT id, number, location, capacity, status, qr_code, qr_generated, description, created_at, updated_at
		FROM tables
		WHERE id = $1
	`

	var location, status, qrCode, description string
	var number, capacity int
	var qrGenerated bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&id, &number, &location, &capacity, &status, &qrCode, &qrGenerated, &description, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	return model.ReconstructTable(
		id,
		number,
		model.TableLocation(location),
		capacity,
		model.TableStatus(status),
		qrCode,
		qrGenerated,
		description,
		createdAt,
		updatedAt,
	), nil
}

func (r *PostgresTableRepository) GetByNumber(ctx context.Context, number int) (*model.Table, error) {
	query := `
		SELECT id, number, location, capacity, status, qr_code, qr_generated, description, created_at, updated_at
		FROM tables
		WHERE number = $1
	`

	var id, location, status, qrCode, description string
	var capacity int
	var qrGenerated bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&id, &number, &location, &capacity, &status, &qrCode, &qrGenerated, &description, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get table by number: %w", err)
	}

	return model.ReconstructTable(
		id,
		number,
		model.TableLocation(location),
		capacity,
		model.TableStatus(status),
		qrCode,
		qrGenerated,
		description,
		createdAt,
		updatedAt,
	), nil
}

func (r *PostgresTableRepository) Update(ctx context.Context, table *model.Table) error {
	query := `
		UPDATE tables
		SET number = $1, location = $2, capacity = $3, status = $4, 
		    qr_code = $5, qr_generated = $6, description = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.ExecContext(ctx, query,
		table.Number(),
		table.Location(),
		table.Capacity(),
		table.Status(),
		table.QRCode(),
		table.IsQRGenerated(),
		table.Description(),
		table.UpdatedAt(),
		table.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update table: %w", err)
	}
	return nil
}

func (r *PostgresTableRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tables WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}
	return nil
}

func (r *PostgresTableRepository) List(ctx context.Context, filter *repository.TableFilter) ([]*model.Table, error) {
	query := `
		SELECT id, number, location, capacity, status, qr_code, qr_generated, description, created_at, updated_at
		FROM tables
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.Location != nil {
			query += fmt.Sprintf(" AND location = $%d", argCount)
			args = append(args, string(*filter.Location))
			argCount++
		}
		if filter.Status != nil {
			query += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, string(*filter.Status))
			argCount++
		}
		if filter.MinCapacity > 0 {
			query += fmt.Sprintf(" AND capacity >= $%d", argCount)
			args = append(args, filter.MinCapacity)
			argCount++
		}
		if filter.MaxCapacity > 0 {
			query += fmt.Sprintf(" AND capacity <= $%d", argCount)
			args = append(args, filter.MaxCapacity)
			argCount++
		}
		if filter.Search != "" {
			query += fmt.Sprintf(" AND description ILIKE $%d", argCount)
			args = append(args, "%"+filter.Search+"%")
			argCount++
		}
	}

	query += " ORDER BY number ASC"

	if filter != nil && filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++

		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []*model.Table
	for rows.Next() {
		var id, location, status, qrCode, description string
		var number, capacity int
		var qrGenerated bool
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &number, &location, &capacity, &status, &qrCode, &qrGenerated, &description, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}

		tables = append(tables, model.ReconstructTable(
			id, number, model.TableLocation(location), capacity,
			model.TableStatus(status), qrCode, qrGenerated, description,
			createdAt, updatedAt,
		))
	}

	return tables, rows.Err()
}

func (r *PostgresTableRepository) Count(ctx context.Context, filter *repository.TableFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM tables WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.Location != nil {
			query += fmt.Sprintf(" AND location = $%d", argCount)
			args = append(args, string(*filter.Location))
			argCount++
		}
		if filter.Status != nil {
			query += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, string(*filter.Status))
			argCount++
		}
		if filter.MinCapacity > 0 {
			query += fmt.Sprintf(" AND capacity >= $%d", argCount)
			args = append(args, filter.MinCapacity)
			argCount++
		}
		if filter.MaxCapacity > 0 {
			query += fmt.Sprintf(" AND capacity <= $%d", argCount)
			args = append(args, filter.MaxCapacity)
			argCount++
		}
		if filter.Search != "" {
			query += fmt.Sprintf(" AND description ILIKE $%d", argCount)
			args = append(args, "%"+filter.Search+"%")
			argCount++
		}
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *PostgresTableRepository) GetAvailableTables(ctx context.Context, location *model.TableLocation) ([]*model.Table, error) {
	query := `
		SELECT id, number, location, capacity, status, qr_code, qr_generated, description, created_at, updated_at
		FROM tables
		WHERE status = 'AVAILABLE'
	`

	args := []interface{}{}
	argCount := 1

	if location != nil {
		query += fmt.Sprintf(" AND location = $%d", argCount)
		args = append(args, string(*location))
	}

	query += " ORDER BY number ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tables: %w", err)
	}
	defer rows.Close()

	var tables []*model.Table
	for rows.Next() {
		var id, loc, status, qrCode, description string
		var number, capacity int
		var qrGenerated bool
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &number, &loc, &capacity, &status, &qrCode, &qrGenerated, &description, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}

		tables = append(tables, model.ReconstructTable(
			id, number, model.TableLocation(loc), capacity,
			model.TableStatus(status), qrCode, qrGenerated, description,
			createdAt, updatedAt,
		))
	}

	return tables, rows.Err()
}

func (r *PostgresTableRepository) ExistsByNumber(ctx context.Context, number int, excludeID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tables WHERE number = $1`
	args := []interface{}{number}

	if excludeID != "" {
		query += " AND id != $2"
		args = append(args, excludeID)
	}
	query += ")"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	return exists, err
}

// Ensure it implements the interface
var _ repository.TableRepository = (*PostgresTableRepository)(nil)
