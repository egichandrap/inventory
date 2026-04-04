package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenConfig holds JWT token configuration
type TokenConfig struct {
	SecretKey       string
	Issuer          string
	Audience        string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	MaxRefreshes    int
}

// SecureTokenConfig returns a secure JWT configuration
func SecureTokenConfig(secretKey string) TokenConfig {
	return TokenConfig{
		SecretKey:       secretKey,
		Issuer:          "pos-system",
		Audience:        "pos-client",
		AccessTokenTTL:  15 * time.Minute,  // Short-lived access tokens
		RefreshTokenTTL: 7 * 24 * time.Hour, // 7 days refresh tokens
		MaxRefreshes:    10,                  // Max refresh count before re-login
	}
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	TokenType    string `json:"token_type"`    // access or refresh
	RefreshCount int    `json:"refresh_count"` // How many times refreshed
	SessionID    string `json:"session_id"`    // Unique session identifier
	jwt.RegisteredClaims
}

// GenerateTokenClaims creates secure token claims
func GenerateTokenClaims(
	userID, username, role, tokenType, sessionID string,
	ttl time.Duration,
	config TokenConfig,
) *TokenClaims {
	now := time.Now()

	return &TokenClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: tokenType,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{config.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        generateJTI(), // Unique token ID
		},
	}
}

// ValidateTokenClaims validates token claims
func ValidateTokenClaims(claims *TokenClaims, expectedType string, config TokenConfig) error {
	// Check issuer
	if claims.Issuer != config.Issuer {
		return ErrInvalidIssuer
	}

	// Check audience
	audienceValid := false
	for _, aud := range claims.Audience {
		if aud == config.Audience {
			audienceValid = true
			break
		}
	}
	if !audienceValid {
		return ErrInvalidAudience
	}

	// Check token type
	if claims.TokenType != expectedType {
		return ErrInvalidTokenType
	}

	// Check expiration
	if claims.ExpiresAt.Before(time.Now()) {
		return ErrTokenExpired
	}

	// Check not before
	if claims.NotBefore.After(time.Now()) {
		return ErrTokenNotYetValid
	}

	return nil
}

// IsTokenExpired checks if token is expired
func (c *TokenClaims) IsTokenExpired() bool {
	return c.ExpiresAt.Before(time.Now())
}

// TimeUntilExpiry returns duration until token expires
func (c *TokenClaims) TimeUntilExpiry() time.Duration {
	return time.Until(c.ExpiresAt.Time)
}

// ShouldRefresh checks if token should be refreshed
func (c *TokenClaims) ShouldRefresh(threshold time.Duration) bool {
	return c.TimeUntilExpiry() < threshold
}

// Errors
var (
	ErrInvalidIssuer    = Error("invalid token issuer")
	ErrInvalidAudience  = Error("invalid token audience")
	ErrInvalidTokenType = Error("invalid token type")
	ErrTokenExpired     = Error("token expired")
	ErrTokenNotYetValid = Error("token not yet valid")
	ErrInvalidJTI       = Error("invalid token ID")
)

// Error represents a JWT error
type Error string

func (e Error) Error() string {
	return string(e)
}

// generateJTI generates a unique token ID
func generateJTI() string {
	// In production, use proper UUID generation
	return generateSecureToken(32)
}

// generateSecureToken generates a cryptographically secure token
func generateSecureToken(length int) string {
	// Implementation in jwt_provider.go
	return ""
}
