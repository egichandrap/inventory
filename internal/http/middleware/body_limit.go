package http

import (
	"net/http"
)

const (
	// DefaultMaxBodySize is the default maximum request body size (1MB)
	DefaultMaxBodySize = 1 << 20 // 1MB
)

// MaxBodySizeMiddleware limits request body size
func MaxBodySizeMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBodySize
	}

	return func(next http.Handler) http.Handler {
		return http.MaxBytesHandler(next, maxBytes)
	}
}

// StrictMaxBodyMiddleware for sensitive endpoints (100KB)
func StrictMaxBodyMiddleware() func(http.Handler) http.Handler {
	return MaxBodySizeMiddleware(100 << 10) // 100KB
}

// LoginMaxBodyMiddleware for login endpoint (10KB)
func LoginMaxBodyMiddleware() func(http.Handler) http.Handler {
	return MaxBodySizeMiddleware(10 << 10) // 10KB
}
