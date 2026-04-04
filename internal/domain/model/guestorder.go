package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GuestOrderStatus represents the status of a guest order
type GuestOrderStatus string

const (
	OrderPending   GuestOrderStatus = "PENDING"
	OrderConfirmed GuestOrderStatus = "CONFIRMED"
	OrderPreparing GuestOrderStatus = "PREPARING"
	OrderReady     GuestOrderStatus = "READY"
	OrderServed    GuestOrderStatus = "SERVED"
	OrderCancelled GuestOrderStatus = "CANCELLED"
)

// GuestOrderPaymentStatus represents payment status
type GuestOrderPaymentStatus string

const (
	PaymentPending   GuestOrderPaymentStatus = "PENDING"
	PaymentPaid      GuestOrderPaymentStatus = "PAID"
	PaymentRefunded  GuestOrderPaymentStatus = "REFUNDED"
)

// GuestOrderItem represents an item in a guest order
type GuestOrderItem struct {
	productID   string
	productName string
	quantity    int
	unitPrice   float64
	subtotal    float64
	notes       string
}

// NewGuestOrderItem creates a guest order item
func NewGuestOrderItem(productID, productName string, quantity int, unitPrice float64, notes string) (*GuestOrderItem, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID tidak boleh kosong")
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity harus lebih dari 0")
	}
	if unitPrice < 0 {
		return nil, fmt.Errorf("harga tidak boleh negatif")
	}

	return &GuestOrderItem{
		productID:   productID,
		productName: productName,
		quantity:    quantity,
		unitPrice:   unitPrice,
		subtotal:    float64(quantity) * unitPrice,
		notes:       notes,
	}, nil
}

// ReconstructGuestOrderItem recreates from database
func ReconstructGuestOrderItem(
	productID, productName string,
	quantity int,
	unitPrice, subtotal float64,
	notes string,
) *GuestOrderItem {
	return &GuestOrderItem{
		productID:   productID,
		productName: productName,
		quantity:    quantity,
		unitPrice:   unitPrice,
		subtotal:    subtotal,
		notes:       notes,
	}
}

// Accessors
func (i *GuestOrderItem) ProductID() string   { return i.productID }
func (i *GuestOrderItem) ProductName() string { return i.productName }
func (i *GuestOrderItem) Quantity() int       { return i.quantity }
func (i *GuestOrderItem) UnitPrice() float64  { return i.unitPrice }
func (i *GuestOrderItem) Subtotal() float64   { return i.subtotal }
func (i *GuestOrderItem) Notes() string       { return i.notes }

// GuestOrder represents an order from a guest customer (no login required)
type GuestOrder struct {
	id              string
	orderNumber     string
	tableID         string
	tableNumber     int
	customerName    string
	customerPhone   string
	items           []GuestOrderItem
	subtotal        float64
	taxAmount       float64
	taxPercent      float64
	discountAmount  float64
	discountPercent float64
	totalAmount     float64
	paymentMethod   PaymentMethod
	paymentStatus   GuestOrderPaymentStatus
	paymentAmount   float64
	changeAmount    float64
	status          GuestOrderStatus
	notes           string
	sessionID       string // For tracking guest session
	createdAt       time.Time
	updatedAt       time.Time
	completedAt     *time.Time
}

// NewGuestOrder creates a new guest order
func NewGuestOrder(tableID string, tableNumber int, customerName string, customerPhone string, sessionID string) (*GuestOrder, error) {
	if tableID == "" {
		return nil, fmt.Errorf("table ID tidak boleh kosong")
	}
	if customerName == "" {
		return nil, fmt.Errorf("nama customer tidak boleh kosong")
	}

	now := time.Now()
	orderNumber := fmt.Sprintf("ORD-%s-%04d", now.Format("20060102"), generateOrderSequence(now))

	return &GuestOrder{
		id:              uuid.New().String(),
		orderNumber:     orderNumber,
		tableID:         tableID,
		tableNumber:     tableNumber,
		customerName:    customerName,
		customerPhone:   customerPhone,
		items:           make([]GuestOrderItem, 0),
		subtotal:        0,
		taxAmount:       0,
		taxPercent:      0,
		discountAmount:  0,
		discountPercent: 0,
		totalAmount:     0,
		paymentMethod:   PaymentCash,
		paymentStatus:   PaymentPending,
		paymentAmount:   0,
		changeAmount:    0,
		status:          OrderPending,
		notes:           "",
		sessionID:       sessionID,
		createdAt:       now,
		updatedAt:       now,
		completedAt:     nil,
	}, nil
}

// ReconstructGuestOrder recreates from database
func ReconstructGuestOrder(
	id, orderNumber, tableID, customerName, customerPhone, sessionID string,
	tableNumber int,
	items []GuestOrderItem,
	subtotal, taxAmount, taxPercent, discountAmount, discountPercent, totalAmount float64,
	paymentMethod PaymentMethod,
	paymentStatus GuestOrderPaymentStatus,
	paymentAmount, changeAmount float64,
	status GuestOrderStatus,
	notes string,
	createdAt, updatedAt time.Time,
	completedAt *time.Time,
) *GuestOrder {
	return &GuestOrder{
		id:              id,
		orderNumber:     orderNumber,
		tableID:         tableID,
		tableNumber:     tableNumber,
		customerName:    customerName,
		customerPhone:   customerPhone,
		items:           items,
		subtotal:        subtotal,
		taxAmount:       taxAmount,
		taxPercent:      taxPercent,
		discountAmount:  discountAmount,
		discountPercent: discountPercent,
		totalAmount:     totalAmount,
		paymentMethod:   paymentMethod,
		paymentStatus:   paymentStatus,
		paymentAmount:   paymentAmount,
		changeAmount:    changeAmount,
		status:          status,
		notes:           notes,
		sessionID:       sessionID,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
		completedAt:     completedAt,
	}
}

// Accessors
func (o *GuestOrder) ID() string                    { return o.id }
func (o *GuestOrder) OrderNumber() string           { return o.orderNumber }
func (o *GuestOrder) TableID() string               { return o.tableID }
func (o *GuestOrder) TableNumber() int              { return o.tableNumber }
func (o *GuestOrder) CustomerName() string          { return o.customerName }
func (o *GuestOrder) CustomerPhone() string         { return o.customerPhone }
func (o *GuestOrder) Items() []GuestOrderItem       { return o.items }
func (o *GuestOrder) Subtotal() float64             { return o.subtotal }
func (o *GuestOrder) TaxAmount() float64            { return o.taxAmount }
func (o *GuestOrder) TaxPercent() float64           { return o.taxPercent }
func (o *GuestOrder) DiscountAmount() float64       { return o.discountAmount }
func (o *GuestOrder) DiscountPercent() float64      { return o.discountPercent }
func (o *GuestOrder) TotalAmount() float64          { return o.totalAmount }
func (o *GuestOrder) PaymentMethod() PaymentMethod  { return o.paymentMethod }
func (o *GuestOrder) PaymentStatus() GuestOrderPaymentStatus { return o.paymentStatus }
func (o *GuestOrder) PaymentAmount() float64        { return o.paymentAmount }
func (o *GuestOrder) ChangeAmount() float64         { return o.changeAmount }
func (o *GuestOrder) Status() GuestOrderStatus      { return o.status }
func (o *GuestOrder) Notes() string                 { return o.notes }
func (o *GuestOrder) SessionID() string             { return o.sessionID }
func (o *GuestOrder) CreatedAt() time.Time          { return o.createdAt }
func (o *GuestOrder) UpdatedAt() time.Time          { return o.updatedAt }
func (o *GuestOrder) CompletedAt() *time.Time       { return o.completedAt }

// AddItem adds an item to the order
func (o *GuestOrder) AddItem(productID, productName string, quantity int, unitPrice float64, notes string) error {
	item, err := NewGuestOrderItem(productID, productName, quantity, unitPrice, notes)
	if err != nil {
		return err
	}

	o.items = append(o.items, *item)
	o.recalculateTotals()
	o.updatedAt = time.Now()
	return nil
}

// RemoveItem removes an item from the order
func (o *GuestOrder) RemoveItem(productID string) {
	for i, item := range o.items {
		if item.productID == productID {
			o.items = append(o.items[:i], o.items[i+1:]...)
			o.recalculateTotals()
			o.updatedAt = time.Now()
			return
		}
	}
}

// UpdateItemQuantity updates the quantity of an item
func (o *GuestOrder) UpdateItemQuantity(productID string, quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("quantity tidak boleh negatif")
	}

	for i, item := range o.items {
		if item.productID == productID {
			if quantity == 0 {
				o.RemoveItem(productID)
				return nil
			}
			o.items[i].quantity = quantity
			o.items[i].subtotal = float64(quantity) * item.unitPrice
			o.recalculateTotals()
			o.updatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("item tidak ditemukan")
}

// ClearItems removes all items from the order
func (o *GuestOrder) ClearItems() {
	o.items = make([]GuestOrderItem, 0)
	o.recalculateTotals()
	o.updatedAt = time.Now()
}

// ApplyTax applies tax to the order
func (o *GuestOrder) ApplyTax(percent float64) {
	o.taxPercent = percent
	o.recalculateTotals()
}

// ApplyDiscount applies discount to the order
func (o *GuestOrder) ApplyDiscount(percent float64) {
	o.discountPercent = percent
	o.recalculateTotals()
}

// ProcessPayment processes payment for the order
func (o *GuestOrder) ProcessPayment(method PaymentMethod, amount float64) error {
	o.paymentMethod = method
	o.paymentAmount = amount

	if amount < o.totalAmount {
		return fmt.Errorf("jumlah pembayaran tidak mencukupi: %.2f < %.2f", amount, o.totalAmount)
	}

	o.changeAmount = amount - o.totalAmount
	o.paymentStatus = PaymentPaid
	o.updatedAt = time.Now()
	return nil
}

// Confirm confirms the order (staff action)
func (o *GuestOrder) Confirm() error {
	if o.status != OrderPending {
		return fmt.Errorf("hanya order pending yang bisa dikonfirmasi")
	}
	o.status = OrderConfirmed
	o.updatedAt = time.Now()
	return nil
}

// MarkPreparing marks order as being prepared
func (o *GuestOrder) MarkPreparing() error {
	if o.status != OrderConfirmed {
		return fmt.Errorf("hanya order confirmed yang bisa diproses")
	}
	o.status = OrderPreparing
	o.updatedAt = time.Now()
	return nil
}

// MarkReady marks order as ready to serve
func (o *GuestOrder) MarkReady() error {
	if o.status != OrderPreparing {
		return fmt.Errorf("hanya order preparing yang bisa ditandai ready")
	}
	o.status = OrderReady
	o.updatedAt = time.Now()
	return nil
}

// MarkServed marks order as served
func (o *GuestOrder) MarkServed() error {
	if o.status != OrderReady {
		return fmt.Errorf("hanya order ready yang bisa ditandai served")
	}
	o.status = OrderServed
	now := time.Now()
	o.completedAt = &now
	o.updatedAt = now
	return nil
}

// Cancel cancels the order
func (o *GuestOrder) Cancel() error {
	if o.status == OrderServed || o.status == OrderCancelled {
		return fmt.Errorf("order sudah selesai atau dibatalkan")
	}
	o.status = OrderCancelled
	now := time.Now()
	o.completedAt = &now
	o.updatedAt = now
	return nil
}

// IsPending checks if order is pending
func (o *GuestOrder) IsPending() bool {
	return o.status == OrderPending
}

// IsCompleted checks if order is completed
func (o *GuestOrder) IsCompleted() bool {
	return o.status == OrderServed || o.status == OrderCancelled
}

// recalculateTotals recalculates order totals
func (o *GuestOrder) recalculateTotals() {
	o.subtotal = 0
	for _, item := range o.items {
		o.subtotal += item.subtotal
	}

	// Apply discount
	o.discountAmount = o.subtotal * (o.discountPercent / 100)
	afterDiscount := o.subtotal - o.discountAmount

	// Apply tax
	o.taxAmount = afterDiscount * (o.taxPercent / 100)
	o.totalAmount = afterDiscount + o.taxAmount
}

// Helper function to generate order sequence number
func generateOrderSequence(t time.Time) int {
	// Simple sequence based on time (not production-ready for high volume)
	// In production, use a proper sequence generator
	return (t.Hour() * 3600 + t.Minute() * 60 + t.Second()) % 10000
}
