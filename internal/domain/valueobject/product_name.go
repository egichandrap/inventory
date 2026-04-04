package valueobject

import (
	"fmt"
	"strings"
)

// ProductName represents a product name value object
type ProductName string

// NewProductName creates a validated ProductName
func NewProductName(name string) (ProductName, error) {
	if name == "" {
		return "", fmt.Errorf("nama produk tidak boleh kosong")
	}

	if len(name) > 255 {
		return "", fmt.Errorf("nama produk maksimal 255 karakter")
	}

	return ProductName(strings.TrimSpace(name)), nil
}

// MustNewProductName creates a ProductName or panics on error
func MustNewProductName(name string) ProductName {
	pn, err := NewProductName(name)
	if err != nil {
		panic(err)
	}
	return pn
}

// String returns the product name as string
func (pn ProductName) String() string {
	return string(pn)
}

// IsEmpty checks if the product name is empty
func (pn ProductName) IsEmpty() bool {
	return strings.TrimSpace(string(pn)) == ""
}
