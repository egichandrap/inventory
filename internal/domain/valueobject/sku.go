package valueobject

import (
	"fmt"
	"regexp"
)

// SKU represents a Stock Keeping Unit value object
type SKU string

// NewSKU creates a validated SKU
func NewSKU(value string) (SKU, error) {
	if value == "" {
		return "", fmt.Errorf("SKU tidak boleh kosong")
	}

	// SKU format: ABC-123 or ABC123 (alphanumeric with optional dash)
	if matched, _ := regexp.MatchString(`^[A-Z0-9\-]+$`, value); !matched {
		return "", fmt.Errorf("format SKU tidak valid, hanya boleh huruf kapital, angka, dan dash")
	}

	return SKU(value), nil
}

// String returns the SKU as a string
func (s SKU) String() string {
	return string(s)
}

// Validate checks if the SKU is valid
func (s SKU) Validate() error {
	if s == "" {
		return fmt.Errorf("SKU tidak boleh kosong")
	}
	return nil
}
