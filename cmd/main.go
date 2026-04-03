package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/jwt-ddd-clean/internal/infrastructure/config"
	"github.com/example/jwt-ddd-clean/internal/infrastructure/database"
	"github.com/example/jwt-ddd-clean/internal/infrastructure/http"
)

func main() {
	// Command line flags (override env config)
	serverMode := flag.Bool("server", false, "Run as HTTP server")
	host := flag.String("host", "", "Server host (overrides .env)")
	port := flag.String("port", "", "Server port (overrides .env)")
	flag.Parse()

	// Load configuration from .env and environment variables
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Override config with CLI flags if provided
	if *host != "" {
		cfg.Server.Host = *host
	}
	if *port != "" {
		cfg.Server.Port = *port
	}

	if *serverMode {
		// Run as HTTP server with PostgreSQL
		runServer(cfg)
	} else {
		// Run demo mode
		runDemo(cfg)
	}
}

func runServer(cfg *config.Config) {
	fmt.Println("🚀 Starting server with PostgreSQL...")

	// Initialize PostgreSQL connection
	dbConfig := database.Config{
		Driver:          "postgres",
		DSN:             cfg.Database.DSN,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	db, err := database.New(dbConfig)
	if err != nil {
		fmt.Printf("❌ Failed to connect to PostgreSQL: %v\n", err)
		fmt.Println("💡 Make sure PostgreSQL is running and credentials are correct")
		fmt.Println("💡 Check your .env file for database configuration")
		os.Exit(1)
	}
	defer db.Close()

	fmt.Printf("✅ Connected to PostgreSQL: %s@%s:%s/%s\n",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Run migrations automatically
	fmt.Println("🔄 Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		fmt.Printf("❌ Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP server with database
	serverConfig := http.ServerConfig{
		Host:            cfg.Server.Host,
		Port:            cfg.Server.Port,
		SecretKey:       cfg.JWT.SecretKey,
		Issuer:          cfg.JWT.Issuer,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
	}

	server := http.NewServerWithDatabase(serverConfig, db.DB)

	if err := server.StartWithGracefulShutdown(); err != nil {
		fmt.Printf("❌ Server failed to start: %v\n", err)
		os.Exit(1)
	}
}

func runDemo(cfg *config.Config) {
	fmt.Println("=== JWT Token Generator (DDD + Clean Architecture) ===")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Server Host: %s\n", cfg.Server.Host)
	fmt.Printf("  Server Port: %s\n", cfg.Server.Port)
	fmt.Printf("  Database:    PostgreSQL (%s@%s:%s/%s)\n",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	fmt.Printf("  JWT Issuer:  %s\n", cfg.JWT.Issuer)
	fmt.Printf("  Access TTL:  %v\n", cfg.JWT.AccessTokenTTL)
	fmt.Println()
	fmt.Println("Run with -server flag to start HTTP API server")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Start server with PostgreSQL (default)")
	fmt.Println("  go run cmd/main.go -server")
	fmt.Println()
	fmt.Println("  # Custom host and port")
	fmt.Println("  go run cmd/main.go -server -host 0.0.0.0 -port 3000")
	fmt.Println()
	fmt.Println("  # Override database connection")
	fmt.Println("  DB_HOST=192.168.1.100 DB_NAME=mydb go run cmd/main.go -server")
}
