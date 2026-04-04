package model

import (
	"time"
)

// SalesReport represents a sales report
type SalesReport struct {
	PeriodStart       time.Time
	PeriodEnd         time.Time
	TotalSales        float64
	TotalTransactions int
	TotalItems        int
	AverageTransaction float64
	TopProducts       []ProductSales
	TopCashiers       []CashierSales
	PaymentMethodBreakdown []PaymentMethodReport
	DailyBreakdown    []DailySales
}

// ProductSales represents product sales data
type ProductSales struct {
	ProductID   string
	ProductName string
	SKU         string
	QuantitySold int
	TotalRevenue float64
}

// CashierSales represents cashier sales data
type CashierSales struct {
	CashierID   string
	CashierName string
	TransactionCount int
	TotalSales  float64
	AverageTransaction float64
}

// PaymentMethodReport represents payment method breakdown
type PaymentMethodReport struct {
	PaymentMethod PaymentMethod
	TransactionCount int
	TotalAmount   float64
	Percentage    float64
}

// DailySales represents daily sales breakdown
type DailySales struct {
	Date         time.Time
	TotalSales   float64
	TransactionCount int
	TotalItems   int
}

// InventoryReport represents an inventory report
type InventoryReport struct {
	TotalItems      int
	TotalValue      float64
	LowStockItems   []LowStockItem
	OutOfStockItems []OutOfStockItem
	TopMovingItems  []ProductSales
	SlowMovingItems []ProductSales
}

// LowStockItem represents a low stock item
type LowStockItem struct {
	ProductID   string
	ProductName string
	SKU         string
	CurrentQty  int
	MinStock    int
}

// OutOfStockItem represents an out of stock item
type OutOfStockItem struct {
	ProductID   string
	ProductName string
	SKU         string
	LastSold    time.Time
}

// ReportFilter defines filter options for reports
type ReportFilter struct {
	StartDate   time.Time
	EndDate     time.Time
	StoreID     string
	CategoryID  string
	CashierID   string
	PaymentMethod PaymentMethod
	GroupBy     string // daily, weekly, monthly
}
