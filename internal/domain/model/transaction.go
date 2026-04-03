package model

import "time"

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentCash      PaymentMethod = "CASH"
	PaymentCard      PaymentMethod = "CARD"
	PaymentQRIS      PaymentMethod = "QRIS"
	PaymentEWallet   PaymentMethod = "E_WALLET"
	PaymentTransfer  PaymentMethod = "TRANSFER"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "PENDING"
	TransactionCompleted TransactionStatus = "COMPLETED"
	TransactionCancelled TransactionStatus = "CANCELLED"
	TransactionRefunded  TransactionStatus = "REFUNDED"
)

// TransactionItem represents an item in a transaction
type TransactionItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

// Transaction represents a sales transaction
type Transaction struct {
	ID              string            `json:"id"`
	TransactionNo   string            `json:"transaction_no"` // e.g., TRX-20260403-0001
	CashierID       string            `json:"cashier_id"`
	CashierName     string            `json:"cashier_name"`
	CustomerName    string            `json:"customer_name,omitempty"`
	Items           []TransactionItem `json:"items"`
	Subtotal        float64           `json:"subtotal"`
	DiscountAmount  float64           `json:"discount_amount"`
	DiscountPercent float64           `json:"discount_percent"`
	TaxAmount       float64           `json:"tax_amount"`
	TaxPercent      float64           `json:"tax_percent"`
	TotalAmount     float64           `json:"total_amount"`
	PaymentMethod   PaymentMethod     `json:"payment_method"`
	PaymentAmount   float64           `json:"payment_amount"` // Amount paid by customer
	ChangeAmount    float64           `json:"change_amount"`  // Change returned
	Status          TransactionStatus `json:"status"`
	Notes           string            `json:"notes,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// NewTransaction creates a new transaction
func NewTransaction(transactionNo string, cashierID, cashierName string) *Transaction {
	now := time.Now()
	return &Transaction{
		TransactionNo: transactionNo,
		CashierID:     cashierID,
		CashierName:   cashierName,
		Items:         make([]TransactionItem, 0),
		Status:        TransactionPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// AddItem adds an item to the transaction
func (t *Transaction) AddItem(productID, productName, sku string, quantity int, unitPrice float64) {
	item := TransactionItem{
		ProductID:   productID,
		ProductName: productName,
		SKU:         sku,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    float64(quantity) * unitPrice,
	}
	t.Items = append(t.Items, item)
	t.RecalculateTotals()
	t.UpdatedAt = time.Now()
}

// RecalculateTotals recalculates all totals
func (t *Transaction) RecalculateTotals() {
	t.Subtotal = 0
	for _, item := range t.Items {
		t.Subtotal += item.Subtotal
	}

	t.TotalAmount = t.Subtotal - t.DiscountAmount + t.TaxAmount
}

// ApplyDiscount applies a discount
func (t *Transaction) ApplyDiscount(amount float64, percent float64) {
	t.DiscountAmount = amount
	t.DiscountPercent = percent
	t.RecalculateTotals()
}

// ApplyTax applies tax
func (t *Transaction) ApplyTax(percent float64) {
	t.TaxPercent = percent
	t.TaxAmount = t.Subtotal * (percent / 100)
	t.RecalculateTotals()
}

// ProcessPayment processes the payment
func (t *Transaction) ProcessPayment(method PaymentMethod, amount float64) error {
	t.PaymentMethod = method
	t.PaymentAmount = amount
	t.ChangeAmount = amount - t.TotalAmount

	if t.ChangeAmount < 0 {
		return nil // Payment insufficient, but don't fail yet
	}

	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the transaction
func (t *Transaction) Cancel() {
	t.Status = TransactionCancelled
	t.UpdatedAt = time.Now()
}
