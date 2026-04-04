package valueobject

import (
	"fmt"
)

// Money represents a monetary value with currency
type Money struct {
	amount   float64
	currency Currency
}

// Currency represents the currency type
type Currency string

const (
	IDR Currency = "IDR"
	USD Currency = "USD"
)

// NewMoney creates a validated Money value object
func NewMoney(amount float64, currency Currency) (*Money, error) {
	if amount < 0 {
		return nil, fmt.Errorf("jumlah uang tidak boleh negatif")
	}

	if !currency.IsValid() {
		return nil, fmt.Errorf("mata uang tidak valid: %s", currency)
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// MustNewMoney creates a Money value object or panics on error
func MustNewMoney(amount float64, currency Currency) *Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// Amount returns the monetary amount
func (m *Money) Amount() float64 {
	return m.amount
}

// Currency returns the currency type
func (m *Money) Currency() Currency {
	return m.currency
}

// IsValid checks if currency is valid
func (c Currency) IsValid() bool {
	switch c {
	case IDR, USD:
		return true
	default:
		return false
	}
}

// Add adds two Money values
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("tidak dapat menambahkan mata uang yang berbeda: %s dan %s", m.currency, other.currency)
	}

	return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract subtracts one Money from another
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("tidak dapat mengurangkan mata uang yang berbeda: %s dan %s", m.currency, other.currency)
	}

	result := m.amount - other.amount
	if result < 0 {
		return nil, fmt.Errorf("hasil pengurangan tidak boleh negatif")
	}

	return NewMoney(result, m.currency)
}

// Multiply multiplies the money by a factor
func (m *Money) Multiply(factor float64) (*Money, error) {
	if factor < 0 {
		return nil, fmt.Errorf("faktor perkalian tidak boleh negatif")
	}

	return NewMoney(m.amount*factor, m.currency)
}

// Percentage calculates a percentage of the money
func (m *Money) Percentage(percent float64) (*Money, error) {
	if percent < 0 || percent > 100 {
		return nil, fmt.Errorf("persentase harus antara 0 dan 100")
	}

	return NewMoney(m.amount*(percent/100), m.currency)
}

// Equals checks if two Money values are equal
func (m *Money) Equals(other *Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns string representation
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}

// Zero returns a zero-value Money
func ZeroMoney(currency Currency) *Money {
	if !currency.IsValid() {
		currency = IDR
	}
	return MustNewMoney(0, currency)
}
