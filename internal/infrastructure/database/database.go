package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config holds PostgreSQL database configuration
type Config struct {
	Driver          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
	config Config
}

// New creates a new PostgreSQL database connection
func New(config Config) (*DB, error) {
	config.Driver = "postgres"

	db, err := sql.Open(config.Driver, config.DSN)
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

// DefaultConfig returns default PostgreSQL configuration
func DefaultConfig() Config {
	return Config{
		Driver:          "postgres",
		DSN:             "host=localhost port=5432 user=postgres password=postgres dbname=jwt_ddd sslmode=disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Driver returns the database driver name
func (db *DB) Driver() string {
	return db.config.Driver
}
