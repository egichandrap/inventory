package model

import (
	"time"
)

// Inventory represents a warehouse inventory item entity
type Inventory struct {
	ID          string    `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	Location    string    `json:"location"`
	MinStock    int       `json:"min_stock"`
	MaxStock    int       `json:"max_stock"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InventoryFilter represents filtering options for inventory queries
type InventoryFilter struct {
	SKU      *string
	Name     *string
	Location *string
	MinQty   *int
	MaxQty   *int
	Limit    int
	Offset   int
}

// PaginatedInventory represents a paginated inventory response
type PaginatedInventory struct {
	Items      []*Inventory `json:"items"`
	Total      int64        `json:"total"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
	TotalPages int          `json:"total_pages"`
}
