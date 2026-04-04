package http

import (
	"net/http"
	"sync"
	"time"
)

// BruteForceConfig holds brute force protection configuration
type BruteForceConfig struct {
	MaxAttempts   int           // Maximum attempts before lockout
	LockoutDuration time.Duration // How long to lock out
	WindowDuration  time.Duration // Time window for counting attempts
}

// DefaultBruteForceConfig returns default configuration
func DefaultBruteForceConfig() BruteForceConfig {
	return BruteForceConfig{
		MaxAttempts:     5,
		LockoutDuration: 15 * time.Minute,
		WindowDuration:  10 * time.Minute,
	}
}

// BruteForceMiddleware protects against brute force attacks
type BruteForceMiddleware struct {
	config BruteForceConfig
	mu     sync.RWMutex
	attempts map[string]*attemptInfo // key: IP or username
}

type attemptInfo struct {
	count      int
	firstAttempt time.Time
	lastAttempt  time.Time
	lockedUntil  time.Time
}

// NewBruteForceMiddleware creates new brute force middleware
func NewBruteForceMiddleware(config BruteForceConfig) *BruteForceMiddleware {
	return &BruteForceMiddleware{
		config:   config,
		attempts: make(map[string]*attemptInfo),
	}
}

// Middleware returns the HTTP middleware
func (bf *BruteForceMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only protect POST requests (login, etc.)
			if r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}

			// Get client IP
			ip := getClientIP(r)

			// Check if IP is locked
			if bf.isLocked(ip) {
				http.Error(w, "Too many attempts. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// Record attempt
			bf.recordAttempt(ip)

			next.ServeHTTP(w, r)
		})
	}
}

// isLocked checks if an IP is currently locked
func (bf *BruteForceMiddleware) isLocked(key string) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	info, exists := bf.attempts[key]
	if !exists {
		return false
	}

	return time.Now().Before(info.lockedUntil)
}

// recordAttempt records a login attempt
func (bf *BruteForceMiddleware) recordAttempt(key string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	now := time.Now()
	info, exists := bf.attempts[key]

	if !exists {
		bf.attempts[key] = &attemptInfo{
			count:      1,
			firstAttempt: now,
			lastAttempt:  now,
		}
		return
	}

	// Reset if window has passed
	if now.Sub(info.firstAttempt) > bf.config.WindowDuration {
		info.count = 1
		info.firstAttempt = now
		info.lastAttempt = now
		return
	}

	info.count++
	info.lastAttempt = now

	// Lock if max attempts exceeded
	if info.count >= bf.config.MaxAttempts {
		info.lockedUntil = now.Add(bf.config.LockoutDuration)
	}
}

// Reset clears all attempt data
func (bf *BruteForceMiddleware) Reset() {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	bf.attempts = make(map[string]*attemptInfo)
}

// Cleanup removes expired entries
func (bf *BruteForceMiddleware) Cleanup() {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	now := time.Now()
	for key, info := range bf.attempts {
		// Remove if lockout has expired and no recent attempts
		if now.After(info.lockedUntil) && now.Sub(info.lastAttempt) > bf.config.WindowDuration*2 {
			delete(bf.attempts, key)
		}
	}
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take first IP in the list
		if idx := indexOf(xff, ','); idx > 0 {
			return xff[:idx]
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

// indexOf returns index of character in string, or -1
func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
