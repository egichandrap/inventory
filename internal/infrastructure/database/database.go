package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds database configuration
type Config struct {
	Driver         string
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	ConnMaxLifetime time.Duration
}

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
	config Config
}

// New creates a new database connection
func New(config Config) (*DB, error) {
	db, err := sql.Open(config.Driver, config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		DB:     db,
		config: config,
	}, nil
}

// DefaultConfig returns default database configuration for SQLite
func DefaultConfig() Config {
	return Config{
		Driver:          "sqlite3",
		DataSourceName:  "./data/inventory.db",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
