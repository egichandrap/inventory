package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer aggregate root
type Customer struct {
	id           string
	name         string
	email        string
	phone        string
	address      string
	loyaltyPoints int
	totalPurchases float64
	createdAt    time.Time
	updatedAt    time.Time
}

// NewCustomer creates a new customer
func NewCustomer(name, email, phone string) (*Customer, error) {
	if name == "" {
		return nil, fmt.Errorf("nama customer tidak boleh kosong")
	}

	now := time.Now()
	return &Customer{
		id:           uuid.New().String(),
		name:         name,
		email:        email,
		phone:        phone,
		loyaltyPoints: 0,
		totalPurchases: 0,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructCustomer recreates a customer from database
func ReconstructCustomer(
	id, name, email, phone, address string,
	loyaltyPoints int,
	totalPurchases float64,
	createdAt, updatedAt time.Time,
) *Customer {
	return &Customer{
		id:           id,
		name:         name,
		email:        email,
		phone:        phone,
		address:      address,
		loyaltyPoints: loyaltyPoints,
		totalPurchases: totalPurchases,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Accessors
func (c *Customer) ID() string              { return c.id }
func (c *Customer) Name() string            { return c.name }
func (c *Customer) Email() string           { return c.email }
func (c *Customer) Phone() string           { return c.phone }
func (c *Customer) Address() string         { return c.address }
func (c *Customer) LoyaltyPoints() int      { return c.loyaltyPoints }
func (c *Customer) TotalPurchases() float64 { return c.totalPurchases }
func (c *Customer) CreatedAt() time.Time    { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time    { return c.updatedAt }

// UpdateContactInfo updates customer contact information
func (c *Customer) UpdateContactInfo(email, phone, address string) {
	c.email = email
	c.phone = phone
	c.address = address
	c.updatedAt = time.Now()
}

// AddLoyaltyPoints adds loyalty points to the customer
func (c *Customer) AddLoyaltyPoints(points int) {
	if points > 0 {
		c.loyaltyPoints += points
		c.updatedAt = time.Now()
	}
}

// RedeemLoyaltyPoints redeems loyalty points
func (c *Customer) RedeemLoyaltyPoints(points int) error {
	if points < 0 {
		return fmt.Errorf("poin yang diredeem tidak boleh negatif")
	}

	if c.loyaltyPoints < points {
		return fmt.Errorf("poin tidak mencukupi")
	}

	c.loyaltyPoints -= points
	c.updatedAt = time.Now()
	return nil
}

// RecordPurchase records a purchase and adds loyalty points
func (c *Customer) RecordPurchase(amount float64) {
	if amount > 0 {
		c.totalPurchases += amount
		// Add 1 point for every 10000 spent
		points := int(amount / 10000)
		c.AddLoyaltyPoints(points)
	}
}

// CanRedeemPoints checks if customer can redeem points
func (c *Customer) CanRedeemPoints(points int) bool {
	return c.loyaltyPoints >= points
}
