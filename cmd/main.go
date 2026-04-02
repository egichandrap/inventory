package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/example/jwt-ddd-clean/internal/infrastructure/http"
)

func main() {
	// Command line flags
	serverMode := flag.Bool("server", false, "Run as HTTP server")
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	if *serverMode {
		// Run as HTTP server
		config := http.ServerConfig{
			Host:            *host,
			Port:            *port,
			SecretKey:       "your-super-secret-key-change-in-production",
			Issuer:          "jwt-ddd-clean",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		}

		server := http.NewServer(config)
		if err := server.StartWithGracefulShutdown(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Run demo mode
		runDemo()
	}
}

func runDemo() {
	fmt.Println("=== JWT Token Generator (DDD + Clean Architecture) ===\n")
	fmt.Println("Run with -server flag to start HTTP API server\n")
}
