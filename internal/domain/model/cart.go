package model

import "time"

// CartItem represents an item in the shopping cart
type CartItem struct {
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	SKU         string    `json:"sku"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Subtotal    float64   `json:"subtotal"`
}

// Cart represents a shopping cart
type Cart struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"` // Cashier who created the cart
	CustomerName string    `json:"customer_name,omitempty"`
	Items       []CartItem `json:"items"`
	TotalAmount float64    `json:"total_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewCart creates a new cart
func NewCart(userID string, customerName string) *Cart {
	now := time.Now()
	return &Cart{
		UserID:       userID,
		CustomerName: customerName,
		Items:        make([]CartItem, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddItem adds or updates an item in the cart
func (c *Cart) AddItem(productID, productName, sku string, quantity int, unitPrice float64) {
	// Check if item already exists
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity += quantity
			c.Items[i].Subtotal = float64(c.Items[i].Quantity) * c.Items[i].UnitPrice
			c.RecalculateTotal()
			c.UpdatedAt = time.Now()
			return
		}
	}

	// Add new item
	newItem := CartItem{
		ProductID:   productID,
		ProductName: productName,
		SKU:         sku,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    float64(quantity) * unitPrice,
	}
	c.Items = append(c.Items, newItem)
	c.RecalculateTotal()
	c.UpdatedAt = time.Now()
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID string) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.RecalculateTotal()
			c.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateItemQuantity updates the quantity of an item
func (c *Cart) UpdateItemQuantity(productID string, quantity int) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				c.RemoveItem(productID)
				return
			}
			c.Items[i].Quantity = quantity
			c.Items[i].Subtotal = float64(quantity) * c.Items[i].UnitPrice
			c.RecalculateTotal()
			c.UpdatedAt = time.Now()
			return
		}
	}
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = make([]CartItem, 0)
	c.TotalAmount = 0
	c.UpdatedAt = time.Now()
}

// RecalculateTotal recalculates the total amount
func (c *Cart) RecalculateTotal() {
	c.TotalAmount = 0
	for _, item := range c.Items {
		c.TotalAmount += item.Subtotal
	}
}
