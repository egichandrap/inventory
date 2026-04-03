package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	"github.com/example/jwt-ddd-clean/internal/handler"
	inventoryhttp "github.com/example/jwt-ddd-clean/internal/http/inventory"
	httpmiddleware "github.com/example/jwt-ddd-clean/internal/http/middleware"
	infrastructurejwt "github.com/example/jwt-ddd-clean/internal/infrastructure/jwt"
	infrarepo "github.com/example/jwt-ddd-clean/internal/infrastructure/repository"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	router     *mux.Router
	config     ServerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string
	Port            string
	SecretKey       string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	DatabasePath    string
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:            "localhost",
		Port:            "8080",
		SecretKey:       "your-super-secret-key-change-in-production",
		Issuer:          "jwt-ddd-clean",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		DatabasePath:    "./data/inventory.db",
	}
}

// setupRoutes configures all API routes using gorilla/mux
func setupRoutes(
	r *mux.Router,
	tokenHTTPHandler *TokenHTTPHandler,
	inventoryHTTPHandler *inventoryhttp.InventoryHTTPHandler,
	authMiddleware *httpmiddleware.AuthMiddleware,
) {
	// Public routes (no authentication required)
	r.HandleFunc("/api/token/generate", tokenHTTPHandler.GenerateToken).Methods(http.MethodPost)
	r.HandleFunc("/api/token/refresh", tokenHTTPHandler.RefreshToken).Methods(http.MethodPost)
	r.HandleFunc("/api/token/validate", tokenHTTPHandler.ValidateToken).Methods(http.MethodPost)
	r.HandleFunc("/api/token/revoke", tokenHTTPHandler.RevokeToken).Methods(http.MethodPost)
	r.HandleFunc("/api/health", tokenHTTPHandler.Health).Methods(http.MethodGet)

	// Protected routes (authentication required) - Inventory
	inventoryRouter := r.PathPrefix("/api/inventory").Subrouter()
	inventoryRouter.Use(authMiddleware.Authenticate)

	inventoryRouter.HandleFunc("", inventoryHTTPHandler.ListInventory).Methods(http.MethodGet)
	inventoryRouter.HandleFunc("", inventoryHTTPHandler.CreateInventory).Methods(http.MethodPost)
	inventoryRouter.HandleFunc("/{id}", inventoryHTTPHandler.GetInventory).Methods(http.MethodGet)
	inventoryRouter.HandleFunc("/{id}", inventoryHTTPHandler.UpdateInventory).Methods(http.MethodPut)
	inventoryRouter.HandleFunc("/{id}", inventoryHTTPHandler.DeleteInventory).Methods(http.MethodDelete)
	inventoryRouter.HandleFunc("/{id}/stock", inventoryHTTPHandler.UpdateStock).Methods(http.MethodPut)
	inventoryRouter.HandleFunc("/{id}/stock/adjust", inventoryHTTPHandler.AdjustStock).Methods(http.MethodPost)

	// Root endpoint - API info
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"service": "jwt-ddd-clean",
			"version": "1.0.0",
			"endpoints": map[string][]string{
				"public": {
					"POST /api/token/generate",
					"POST /api/token/refresh",
					"POST /api/token/validate",
					"POST /api/token/revoke",
					"GET  /api/health",
				},
				"protected": {
					"GET    /api/inventory",
					"POST   /api/inventory",
					"GET    /api/inventory/{id}",
					"PUT    /api/inventory/{id}",
					"DELETE /api/inventory/{id}",
					"PUT    /api/inventory/{id}/stock",
					"POST   /api/inventory/{id}/stock/adjust",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}).Methods(http.MethodGet)
}

// NewServer creates a new HTTP server with gorilla/mux
func NewServer(config ServerConfig) *Server {
	// Infrastructure layer - JWT
	jwtProvider := infrastructurejwt.NewProvider(infrastructurejwt.Config{
		SecretKey: config.SecretKey,
		Issuer:    config.Issuer,
		Algorithm: "HS256",
	})

	// Infrastructure layer - Repositories (In-Memory)
	var tokenRepo repository.TokenRepository = infrarepo.NewMemoryTokenRepository()
	var inventoryRepo repository.InventoryRepository = infrarepo.NewMemoryInventoryRepository()

	// Domain layer - Services
	tokenService := service.NewTokenService(
		tokenRepo,
		jwtProvider,
		config.AccessTokenTTL,
		config.RefreshTokenTTL,
	)

	inventoryService := service.NewInventoryService(inventoryRepo)

	// Handler layer - Token
	tokenHandler := handler.NewTokenHandler(tokenService, &handler.UserService{})
	tokenHTTPHandler := NewTokenHTTPHandler(tokenHandler)

	// Handler layer - Inventory
	inventoryHTTPHandler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

	// Middleware
	authMiddleware := httpmiddleware.NewAuthMiddleware(tokenService)

	// Create mux router
	r := mux.NewRouter()

	// Setup routes
	setupRoutes(r, tokenHTTPHandler, inventoryHTTPHandler, authMiddleware)

	server := &http.Server{
		Addr:         config.Host + ":" + config.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: server,
		router:     r,
		config:     config,
	}
}

// NewServerWithDatabase creates a new HTTP server with PostgreSQL database connection using gorilla/mux
func NewServerWithDatabase(config ServerConfig, db *sql.DB) *Server {
	// Infrastructure layer - JWT
	jwtProvider := infrastructurejwt.NewProvider(infrastructurejwt.Config{
		SecretKey: config.SecretKey,
		Issuer:    config.Issuer,
		Algorithm: "HS256",
	})

	// Infrastructure layer - Repositories (PostgreSQL)
	var tokenRepo repository.TokenRepository = infrarepo.NewMemoryTokenRepository()
	var inventoryRepo repository.InventoryRepository = infrarepo.NewPostgresInventoryRepository(db)

	// Domain layer - Services
	tokenService := service.NewTokenService(
		tokenRepo,
		jwtProvider,
		config.AccessTokenTTL,
		config.RefreshTokenTTL,
	)

	inventoryService := service.NewInventoryService(inventoryRepo)

	// Handler layer - Token
	tokenHandler := handler.NewTokenHandler(tokenService, &handler.UserService{})
	tokenHTTPHandler := NewTokenHTTPHandler(tokenHandler)

	// Handler layer - Inventory
	inventoryHTTPHandler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

	// Middleware
	authMiddleware := httpmiddleware.NewAuthMiddleware(tokenService)

	// Create mux router
	r := mux.NewRouter()

	// Setup routes
	setupRoutes(r, tokenHTTPHandler, inventoryHTTPHandler, authMiddleware)

	server := &http.Server{
		Addr:         config.Host + ":" + config.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: server,
		router:     r,
		config:     config,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("🚀 Starting JWT API Server on http://%s:%s\n", s.config.Host, s.config.Port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("   Public Endpoints:\n")
	fmt.Printf("     POST /api/token/generate - Generate new tokens\n")
	fmt.Printf("     POST /api/token/refresh  - Refresh access token\n")
	fmt.Printf("     POST /api/token/validate - Validate token\n")
	fmt.Printf("     POST /api/token/revoke   - Revoke token\n")
	fmt.Printf("     GET  /api/health         - Health check\n")
	fmt.Printf("   Protected Endpoints (require authentication):\n")
	fmt.Printf("     GET    /api/inventory           - List inventory items\n")
	fmt.Printf("     POST   /api/inventory           - Create inventory item\n")
	fmt.Printf("     GET    /api/inventory/{id}      - Get inventory item\n")
	fmt.Printf("     PUT    /api/inventory/{id}      - Update inventory item\n")
	fmt.Printf("     DELETE /api/inventory/{id}      - Delete inventory item\n")
	fmt.Printf("     PUT    /api/inventory/{id}/stock - Update stock quantity\n")
	fmt.Printf("     POST   /api/inventory/{id}/stock/adjust - Adjust stock\n")
	fmt.Println()

	return s.httpServer.ListenAndServe()
}

// StartWithGracefulShutdown starts the server with graceful shutdown
func (s *Server) StartWithGracefulShutdown() error {
	// Start server in goroutine
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server failed to start: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("❌ Server forced to shutdown: %v\n", err)
		return err
	}

	fmt.Println("✅ Server stopped gracefully")
	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
