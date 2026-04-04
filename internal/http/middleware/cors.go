package http

import (
	"net/http"
	"strings"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns secure default CORS config
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{}, // Empty = deny all by default
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           3600,
	}
}

// CORSMiddleware handles CORS preflight and requests
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	// Normalize origins to lowercase
	for i, origin := range config.AllowedOrigins {
		config.AllowedOrigins[i] = strings.ToLower(origin)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			originLower := strings.ToLower(origin)

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == originLower {
					allowed = true
					break
				}
			}

			// If origin not in allowed list, deny request
			if origin != "" && !allowed {
				http.Error(w, "CORS: Origin not allowed", http.StatusForbidden)
				return
			}

			// Set CORS headers
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ","))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", string(rune(config.MaxAge)))

			// Handle preflight
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
