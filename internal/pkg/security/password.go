package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher handles secure password hashing
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

// HashPassword hashes a password securely
func (h *PasswordHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a password with a hash
func (h *PasswordHasher) ComparePassword(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return fmt.Errorf("password and hashed password cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

// GenerateSecureToken generates a cryptographically secure token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// MustGenerateSecureToken generates a token or panics
func MustGenerateSecureToken(length int) string {
	token, err := GenerateSecureToken(length)
	if err != nil {
		panic(err)
	}
	return token
}
