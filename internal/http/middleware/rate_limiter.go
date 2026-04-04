package http

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	clients map[string]*client
	mu      sync.RWMutex
	limit   int
	window  time.Duration
}

type client struct {
	count     int
	lastAccess time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}

	// Start cleanup goroutine
	go rl.cleanup()
	return rl
}

// RateLimitMiddleware creates a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		clientIP := r.RemoteAddr

		// Check rate limit
		if !rl.allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow checks if a request is allowed
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[key]
	now := time.Now()

	if !exists {
		rl.clients[key] = &client{
			count:      1,
			lastAccess: now,
		}
		return true
	}

	// Reset if window has passed
	if now.Sub(c.lastAccess) > rl.window {
		c.count = 1
		c.lastAccess = now
		return true
	}

	// Increment count
	c.count++
	if c.count > rl.limit {
		return false
	}

	return true
}

// cleanup removes expired entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, c := range rl.clients {
			if now.Sub(c.lastAccess) > rl.window*2 {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}
