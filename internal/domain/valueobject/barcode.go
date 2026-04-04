package valueobject

import (
	"fmt"
	"regexp"
)

// BarcodeType represents the type of barcode
type BarcodeType string

const (
	BarcodeEAN13   BarcodeType = "EAN13"
	BarcodeEAN8    BarcodeType = "EAN8"
	BarcodeCode128 BarcodeType = "CODE128"
	BarcodeQR      BarcodeType = "QR"
	BarcodeUPC     BarcodeType = "UPC"
)

// Barcode represents a barcode value object
type Barcode struct {
	value     string
	barcodeType BarcodeType
}

// NewBarcode creates a validated Barcode
func NewBarcode(value string, barcodeType BarcodeType) (*Barcode, error) {
	if value == "" {
		return nil, fmt.Errorf("nilai barcode tidak boleh kosong")
	}

	if !barcodeType.IsValid() {
		return nil, fmt.Errorf("tipe barcode tidak valid: %s", barcodeType)
	}

	// Validate format based on type
	if err := validateBarcodeFormat(value, barcodeType); err != nil {
		return nil, err
	}

	return &Barcode{
		value:     value,
		barcodeType: barcodeType,
	}, nil
}

// MustNewBarcode creates a Barcode or panics on error
func MustNewBarcode(value string, barcodeType BarcodeType) *Barcode {
	b, err := NewBarcode(value, barcodeType)
	if err != nil {
		panic(err)
	}
	return b
}

// Value returns the barcode value
func (b *Barcode) Value() string {
	return b.value
}

// Type returns the barcode type
func (b *Barcode) Type() BarcodeType {
	return b.barcodeType
}

// IsValid checks if the barcode is valid
func (b *Barcode) IsValid() bool {
	return b.barcodeType.IsValid()
}

// IsValid checks if the barcode type is valid
func (t BarcodeType) IsValid() bool {
	switch t {
	case BarcodeEAN13, BarcodeEAN8, BarcodeCode128, BarcodeQR, BarcodeUPC:
		return true
	default:
		return false
	}
}

// String returns string representation
func (b *Barcode) String() string {
	return fmt.Sprintf("%s:%s", b.barcodeType, b.value)
}

// Equals checks if two barcodes are equal
func (b *Barcode) Equals(other *Barcode) bool {
	return b.value == other.value && b.barcodeType == other.barcodeType
}

// validateBarcodeFormat validates the format based on barcode type
func validateBarcodeFormat(value string, barcodeType BarcodeType) error {
	switch barcodeType {
	case BarcodeEAN13:
		// EAN13: 13 digits
		if matched, _ := regexp.MatchString(`^\d{13}$`, value); !matched {
			return fmt.Errorf("EAN13 harus 13 digit angka")
		}
	case BarcodeEAN8:
		// EAN8: 8 digits
		if matched, _ := regexp.MatchString(`^\d{8}$`, value); !matched {
			return fmt.Errorf("EAN8 harus 8 digit angka")
		}
	case BarcodeUPC:
		// UPC: 12 digits
		if matched, _ := regexp.MatchString(`^\d{12}$`, value); !matched {
			return fmt.Errorf("UPC harus 12 digit angka")
		}
	case BarcodeCode128:
		// Code128: alphanumeric
		if len(value) == 0 {
			return fmt.Errorf("Code128 tidak boleh kosong")
		}
	case BarcodeQR:
		// QR: any text
		if len(value) == 0 {
			return fmt.Errorf("QR code tidak boleh kosong")
		}
	}

	return nil
}
