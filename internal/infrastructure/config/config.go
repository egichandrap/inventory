package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	DSN             string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey       string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "jwt_ddd"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", "5m"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			Issuer:          getEnv("JWT_ISSUER", "jwt-ddd-clean"),
			AccessTokenTTL:  getEnvDuration("JWT_ACCESS_TOKEN_TTL", "15m"),
			RefreshTokenTTL: getEnvDuration("JWT_REFRESH_TOKEN_TTL", "168h"), // 7 days
		},
	}

	// Build PostgreSQL DSN
	cfg.Database.DSN = buildPostgresDSN(cfg.Database)

	return cfg, nil
}

// buildPostgresDSN builds PostgreSQL connection string
func buildPostgresDSN(db DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode,
	)
}

// getEnv returns environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvInt returns environment variable as int with fallback
func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}

// getEnvDuration returns environment variable as time.Duration with fallback
func getEnvDuration(key, fallback string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		d, _ := time.ParseDuration(fallback)
		return d
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		d, _ := time.ParseDuration(fallback)
		return d
	}
	return d
}
