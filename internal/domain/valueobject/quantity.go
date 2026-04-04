package valueobject

import (
	"fmt"
)

// Quantity represents a product quantity value object
type Quantity int

// NewQuantity creates a validated Quantity
func NewQuantity(value int) (Quantity, error) {
	if value < 0 {
		return 0, fmt.Errorf("quantity tidak boleh negatif")
	}

	return Quantity(value), nil
}

// MustNewQuantity creates a Quantity or panics on error
func MustNewQuantity(value int) Quantity {
	q, err := NewQuantity(value)
	if err != nil {
		panic(err)
	}
	return q
}

// Int returns the quantity as an int
func (q Quantity) Int() int {
	return int(q)
}

// Add adds two quantities
func (q Quantity) Add(other Quantity) Quantity {
	return q + other
}

// Subtract subtracts one quantity from another
func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	result := q - other
	if result < 0 {
		return 0, fmt.Errorf("quantity tidak boleh negatif")
	}
	return result, nil
}

// IsZero checks if quantity is zero
func (q Quantity) IsZero() bool {
	return q == 0
}

// IsPositive checks if quantity is positive
func (q Quantity) IsPositive() bool {
	return q > 0
}
