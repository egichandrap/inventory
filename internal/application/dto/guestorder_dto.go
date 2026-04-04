package dto

import (
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// GuestOrderItemResponse represents item in order response
type GuestOrderItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	Notes       string  `json:"notes,omitempty"`
}

// GuestOrderResponse represents guest order response
type GuestOrderResponse struct {
	ID              string                        `json:"id"`
	OrderNumber     string                        `json:"order_number"`
	TableID         string                        `json:"table_id"`
	TableNumber     int                           `json:"table_number"`
	CustomerName    string                        `json:"customer_name"`
	CustomerPhone   string                        `json:"customer_phone,omitempty"`
	Items           []GuestOrderItemResponse      `json:"items"`
	Subtotal        float64                       `json:"subtotal"`
	TaxAmount       float64                       `json:"tax_amount"`
	TaxPercent      float64                       `json:"tax_percent"`
	DiscountAmount  float64                       `json:"discount_amount"`
	DiscountPercent float64                       `json:"discount_percent"`
	TotalAmount     float64                       `json:"total_amount"`
	PaymentMethod   model.PaymentMethod           `json:"payment_method"`
	PaymentStatus   model.GuestOrderPaymentStatus `json:"payment_status"`
	PaymentAmount   float64                       `json:"payment_amount"`
	ChangeAmount    float64                       `json:"change_amount"`
	Status          model.GuestOrderStatus        `json:"status"`
	Notes           string                        `json:"notes,omitempty"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
	CompletedAt     *time.Time                    `json:"completed_at,omitempty"`
}

// CreateGuestOrderRequest represents create guest order request
type CreateGuestOrderRequest struct {
	TableID       string `json:"table_id" validate:"required"`
	CustomerName  string `json:"customer_name" validate:"required"`
	CustomerPhone string `json:"customer_phone"`
	SessionID     string `json:"session_id"`
}

// AddOrderItemRequest represents add item request
type AddOrderItemRequest struct {
	ProductID   string  `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required,min=0"`
	Notes       string  `json:"notes"`
}

// GuestCheckoutRequest represents checkout request for guest orders
type GuestCheckoutRequest struct {
	PaymentMethod model.PaymentMethod `json:"payment_method" validate:"required"`
	PaymentAmount float64             `json:"payment_amount" validate:"required,min=0"`
	Notes         string              `json:"notes"`
}

// GuestOrderListResponse represents paginated order list
type GuestOrderListResponse struct {
	Orders     []GuestOrderResponse `json:"orders"`
	Total      int64                `json:"total"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
	TotalPages int                  `json:"total_pages"`
}

// GuestOrderSalesResponse represents sales summary response
type GuestOrderSalesResponse struct {
	TotalSales  float64 `json:"total_sales"`
	TotalOrders int     `json:"total_orders"`
	TotalItems  int     `json:"total_items"`
	Date        string  `json:"date"`
}

// ToGuestOrderItemResponse converts domain item to DTO
func ToGuestOrderItemResponse(item model.GuestOrderItem) GuestOrderItemResponse {
	return GuestOrderItemResponse{
		ProductID:   item.ProductID(),
		ProductName: item.ProductName(),
		Quantity:    item.Quantity(),
		UnitPrice:   item.UnitPrice(),
		Subtotal:    item.Subtotal(),
		Notes:       item.Notes(),
	}
}

// ToGuestOrderResponse converts domain order to DTO
func ToGuestOrderResponse(order *model.GuestOrder) GuestOrderResponse {
	items := make([]GuestOrderItemResponse, len(order.Items()))
	for i, item := range order.Items() {
		items[i] = ToGuestOrderItemResponse(item)
	}

	return GuestOrderResponse{
		ID:              order.ID(),
		OrderNumber:     order.OrderNumber(),
		TableID:         order.TableID(),
		TableNumber:     order.TableNumber(),
		CustomerName:    order.CustomerName(),
		CustomerPhone:   order.CustomerPhone(),
		Items:           items,
		Subtotal:        order.Subtotal(),
		TaxAmount:       order.TaxAmount(),
		TaxPercent:      order.TaxPercent(),
		DiscountAmount:  order.DiscountAmount(),
		DiscountPercent: order.DiscountPercent(),
		TotalAmount:     order.TotalAmount(),
		PaymentMethod:   order.PaymentMethod(),
		PaymentStatus:   order.PaymentStatus(),
		PaymentAmount:   order.PaymentAmount(),
		ChangeAmount:    order.ChangeAmount(),
		Status:          order.Status(),
		Notes:           order.Notes(),
		CreatedAt:       order.CreatedAt(),
		UpdatedAt:       order.UpdatedAt(),
		CompletedAt:     order.CompletedAt(),
	}
}

// ToGuestOrderListResponse converts domain orders to DTO
func ToGuestOrderListResponse(orders []*model.GuestOrder, total int64, limit, offset int) GuestOrderListResponse {
	responses := make([]GuestOrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = ToGuestOrderResponse(o)
	}

	totalPages := int(total) / limit
	if limit > 0 && int(total)%limit > 0 {
		totalPages++
	}

	return GuestOrderListResponse{
		Orders:     responses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}
}

// GuestOrderSalesSummary represents today's sales summary
type GuestOrderSalesSummary struct {
	TotalSales        float64 `json:"total_sales"`
	TotalOrders       int     `json:"total_orders"`
	TotalItems        int     `json:"total_items"`
	Date              string  `json:"date"`
}

// ToGuestOrderSalesResponse converts sales summary to DTO
func ToGuestOrderSalesResponse(summary *GuestOrderSalesSummary) GuestOrderSalesResponse {
	return GuestOrderSalesResponse{
		TotalSales:  summary.TotalSales,
		TotalOrders: summary.TotalOrders,
		TotalItems:  summary.TotalItems,
		Date:        summary.Date,
	}
}
