package http

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CSRFConfig holds CSRF configuration
type CSRFConfig struct {
	CookieName   string
	HeaderName   string
	TokenLength  int
	CookieSecure bool
	CookieTTL    time.Duration
}

// DefaultCSRFConfig returns default CSRF config
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		CookieName:   "csrf_token",
		HeaderName:   "X-CSRF-Token",
		TokenLength:  32,
		CookieSecure: true,
		CookieTTL:    24 * time.Hour,
	}
}

// CSRFMiddleware provides CSRF protection
type CSRFMiddleware struct {
	config CSRFConfig
	tokens sync.Map // In-memory token store (use Redis in production)
}

// NewCSRFMiddleware creates a new CSRF middleware
func NewCSRFMiddleware(config CSRFConfig) *CSRFMiddleware {
	return &CSRFMiddleware{
		config: config,
	}
}

// GenerateToken generates a new CSRF token
func (m *CSRFMiddleware) GenerateToken() (string, error) {
	bytes := make([]byte, m.config.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// Middleware handles CSRF validation
func (m *CSRFMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF for safe methods
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Get token from header
			token := r.Header.Get(m.config.HeaderName)
			if token == "" {
				http.Error(w, "CSRF token missing", http.StatusForbidden)
				return
			}

			// Validate token (in production, validate against session/store)
			if !m.validateToken(token) {
				http.Error(w, "CSRF token invalid", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// validateToken validates the CSRF token
func (m *CSRFMiddleware) validateToken(token string) bool {
	// In production, validate against stored tokens
	// For now, just validate format
	if len(token) != m.config.TokenLength*2 {
		return false
	}

	// Additional validation: check against blacklist/store
	_, exists := m.tokens.Load(token)
	if exists {
		return false // Token already used (one-time use)
	}

	// Mark token as used
	m.tokens.Store(token, time.Now())
	return true
}

// GetToken returns CSRF token for forms
func (m *CSRFMiddleware) GetToken(w http.ResponseWriter, r *http.Request) (string, error) {
	token, err := m.GenerateToken()
	if err != nil {
		return "", err
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     m.config.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   m.config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(m.config.CookieTTL),
	})

	return token, nil
}

// CleanupTokens removes expired tokens from memory
func (m *CSRFMiddleware) CleanupTokens() {
	m.tokens.Range(func(key, value interface{}) bool {
		tokenTime := value.(time.Time)
		if time.Since(tokenTime) > m.config.CookieTTL {
			m.tokens.Delete(key)
		}
		return true
	})
}

// SafeMethods returns list of HTTP methods that don't require CSRF
func SafeMethods() []string {
	return []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
	}
}

// IsSafeMethod checks if method is safe (doesn't require CSRF)
func IsSafeMethod(method string) bool {
	for _, safe := range SafeMethods() {
		if method == safe {
			return true
		}
	}
	return false
}

// ValidateCSRFToken validates token from request
func ValidateCSRFToken(r *http.Request, cookieName, headerName string) error {
	if IsSafeMethod(r.Method) {
		return nil
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return fmt.Errorf("CSRF cookie missing: %w", err)
	}

	token := r.Header.Get(headerName)
	if token == "" {
		// Fallback to form value
		token = r.FormValue("_csrf")
	}

	if token == "" {
		return fmt.Errorf("CSRF token missing")
	}

	if token != cookie.Value {
		return fmt.Errorf("CSRF token mismatch")
	}

	return nil
}
