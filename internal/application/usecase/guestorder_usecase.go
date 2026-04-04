package usecase

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/application/dto"
	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
)

// GuestOrderUsecase defines the guest order usecase interface
type GuestOrderUsecase interface {
	CreateOrder(ctx context.Context, req dto.CreateGuestOrderRequest) (*dto.GuestOrderResponse, error)
	GetOrder(ctx context.Context, id string) (*dto.GuestOrderResponse, error)
	AddItem(ctx context.Context, orderID string, req dto.AddOrderItemRequest) (*dto.GuestOrderResponse, error)
	UpdateItemQuantity(ctx context.Context, orderID string, productID string, quantity int) (*dto.GuestOrderResponse, error)
	RemoveItem(ctx context.Context, orderID, productID string) (*dto.GuestOrderResponse, error)
	Checkout(ctx context.Context, orderID string, req dto.GuestCheckoutRequest) (*dto.GuestOrderResponse, error)
	CancelOrder(ctx context.Context, id string) (*dto.GuestOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, id string, status model.GuestOrderStatus) (*dto.GuestOrderResponse, error)
	ListOrders(ctx context.Context, filter repository.GuestOrderFilter) (*dto.GuestOrderListResponse, error)
	GetPendingOrders(ctx context.Context) ([]dto.GuestOrderResponse, error)
	GetActiveOrders(ctx context.Context) ([]dto.GuestOrderResponse, error)
	GetOrdersByTable(ctx context.Context, tableID string) ([]dto.GuestOrderResponse, error)
	GetTodaySales(ctx context.Context) (*dto.GuestOrderSalesResponse, error)
}

type guestOrderUsecase struct {
	guestOrderService *service.GuestOrderService
}

// NewGuestOrderUsecase creates a new GuestOrderUsecase
func NewGuestOrderUsecase(guestOrderService *service.GuestOrderService) GuestOrderUsecase {
	return &guestOrderUsecase{
		guestOrderService: guestOrderService,
	}
}

func (u *guestOrderUsecase) CreateOrder(ctx context.Context, req dto.CreateGuestOrderRequest) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.CreateOrder(ctx, req.TableID, req.CustomerName, req.CustomerPhone, req.SessionID)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) GetOrder(ctx context.Context, id string) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) AddItem(ctx context.Context, orderID string, req dto.AddOrderItemRequest) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.AddItem(ctx, orderID, req.ProductID, req.ProductName, req.Quantity, req.UnitPrice, req.Notes)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) UpdateItemQuantity(ctx context.Context, orderID string, productID string, quantity int) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.UpdateItemQuantity(ctx, orderID, productID, quantity)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) RemoveItem(ctx context.Context, orderID, productID string) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.RemoveItem(ctx, orderID, productID)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) Checkout(ctx context.Context, orderID string, req dto.GuestCheckoutRequest) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.ProcessCheckout(ctx, orderID, req.PaymentMethod, req.PaymentAmount)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) CancelOrder(ctx context.Context, id string) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.CancelOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) UpdateOrderStatus(ctx context.Context, id string, status model.GuestOrderStatus) (*dto.GuestOrderResponse, error) {
	order, err := u.guestOrderService.UpdateOrderStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderResponse(order)
	return &resp, nil
}

func (u *guestOrderUsecase) ListOrders(ctx context.Context, filter repository.GuestOrderFilter) (*dto.GuestOrderListResponse, error) {
	paginated, err := u.guestOrderService.ListOrdersWithPagination(ctx, filter)
	if err != nil {
		return nil, err
	}

	resp := dto.ToGuestOrderListResponse(paginated.Orders, paginated.Total, paginated.Limit, paginated.Offset)
	return &resp, nil
}

func (u *guestOrderUsecase) GetPendingOrders(ctx context.Context) ([]dto.GuestOrderResponse, error) {
	orders, err := u.guestOrderService.GetPendingOrders(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GuestOrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = dto.ToGuestOrderResponse(o)
	}

	return responses, nil
}

func (u *guestOrderUsecase) GetActiveOrders(ctx context.Context) ([]dto.GuestOrderResponse, error) {
	orders, err := u.guestOrderService.GetActiveOrders(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GuestOrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = dto.ToGuestOrderResponse(o)
	}

	return responses, nil
}

func (u *guestOrderUsecase) GetOrdersByTable(ctx context.Context, tableID string) ([]dto.GuestOrderResponse, error) {
	orders, err := u.guestOrderService.GetOrdersByTable(ctx, tableID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GuestOrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = dto.ToGuestOrderResponse(o)
	}

	return responses, nil
}

func (u *guestOrderUsecase) GetTodaySales(ctx context.Context) (*dto.GuestOrderSalesResponse, error) {
	summary, err := u.guestOrderService.GetTodaySales(ctx)
	if err != nil {
		return nil, err
	}

	// Convert service summary to dto summary
	dtoSummary := &dto.GuestOrderSalesSummary{
		TotalSales:  summary.TotalSales,
		TotalOrders: summary.TotalOrders,
		TotalItems:  summary.TotalItems,
		Date:        summary.Date,
	}

	resp := dto.ToGuestOrderSalesResponse(dtoSummary)
	return &resp, nil
}
