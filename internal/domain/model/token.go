package model

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrInvalidToken = errors.New("invalid token")
)

// Token represents a JWT token entity
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TokenClaims represents the claims within a JWT token
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	Access  *Token
	Refresh *Token
}
