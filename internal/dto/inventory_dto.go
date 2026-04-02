package dto

// CreateInventoryRequest represents a request to create an inventory item
type CreateInventoryRequest struct {
	SKU         string  `json:"sku" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Unit        string  `json:"unit" binding:"required"`
	Location    string  `json:"location"`
	MinStock    int     `json:"min_stock"`
	MaxStock    int     `json:"max_stock"`
	Price       float64 `json:"price"`
}

// UpdateInventoryRequest represents a request to update an inventory item
type UpdateInventoryRequest struct {
	ID          string  `json:"id" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Unit        string  `json:"unit" binding:"required"`
	Location    string  `json:"location"`
	MinStock    int     `json:"min_stock"`
	MaxStock    int     `json:"max_stock"`
	Price       float64 `json:"price"`
}

// UpdateStockRequest represents a request to update stock quantity
type UpdateStockRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

// AdjustStockRequest represents a request to adjust stock quantity
type AdjustStockRequest struct {
	Adjustment int `json:"adjustment" binding:"required"`
}

// InventoryResponse represents an inventory item response
type InventoryResponse struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Unit        string  `json:"unit"`
	Location    string  `json:"location"`
	MinStock    int     `json:"min_stock"`
	MaxStock    int     `json:"max_stock"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// InventoryListResponse represents a paginated inventory list response
type InventoryListResponse struct {
	Items      []*InventoryResponse `json:"items"`
	Total      int64                `json:"total"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
	TotalPages int                  `json:"total_pages"`
}

// StockUpdateResponse represents a stock update response
type StockUpdateResponse struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	PreviousQty int     `json:"previous_quantity"`
	UpdatedAt   string  `json:"updated_at"`
}
