package persistence

import (
	"context"
	"database/sql"
	"fmt"
)

// UnitOfWork manages database transactions
type UnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

// NewUnitOfWork creates a new UnitOfWork
func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Begin starts a new transaction
func (uow *UnitOfWork) Begin(ctx context.Context) error {
	if uow.tx != nil {
		return fmt.Errorf("transaction already in progress")
	}

	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	uow.tx = tx
	return nil
}

// Commit commits the transaction
func (uow *UnitOfWork) Commit() error {
	if uow.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	if err := uow.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	uow.tx = nil
	return nil
}

// Rollback rolls back the transaction
func (uow *UnitOfWork) Rollback() error {
	if uow.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	if err := uow.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	uow.tx = nil
	return nil
}

// Tx returns the current transaction
func (uow *UnitOfWork) Tx() *sql.Tx {
	return uow.tx
}

// IsInTransaction checks if a transaction is in progress
func (uow *UnitOfWork) IsInTransaction() bool {
	return uow.tx != nil
}

// ExecuteInTransaction executes a function within a transaction
// Automatically commits on success or rolls back on error
func (uow *UnitOfWork) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	defer func() {
		if uow.IsInTransaction() {
			_ = uow.Rollback()
		}
	}()

	if err := fn(ctx); err != nil {
		return err
	}

	return uow.Commit()
}
