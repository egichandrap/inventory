package service

import (
	"context"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	"github.com/example/jwt-ddd-clean/internal/pkg/errors"
)

// GuestOrderService handles guest order business logic
type GuestOrderService struct {
	orderRepo     repository.GuestOrderRepository
	tableRepo     repository.TableRepository
	inventoryRepo repository.InventoryRepository
	taxRate       float64
}

// NewGuestOrderService creates a new GuestOrderService
func NewGuestOrderService(
	orderRepo repository.GuestOrderRepository,
	tableRepo repository.TableRepository,
	inventoryRepo repository.InventoryRepository,
	taxRate float64,
) *GuestOrderService {
	return &GuestOrderService{
		orderRepo:     orderRepo,
		tableRepo:     tableRepo,
		inventoryRepo: inventoryRepo,
		taxRate:       taxRate,
	}
}

// CreateOrder creates a new guest order
func (s *GuestOrderService) CreateOrder(ctx context.Context, tableID, customerName, customerPhone, sessionID string) (*model.GuestOrder, error) {
	// Validate table exists and is available
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil || table == nil {
		return nil, errors.NewNotFoundError("meja", "id", tableID)
	}

	if !table.IsAvailable() && !table.IsOccupied() {
		return nil, errors.NewValidationError("meja tidak tersedia untuk order")
	}

	// Create order
	order, err := model.NewGuestOrder(tableID, table.Number(), customerName, customerPhone, sessionID)
	if err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	// Apply default tax
	order.ApplyTax(s.taxRate)

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal membuat order: %v", err)
	}

	// Mark table as occupied
	if table.IsAvailable() {
		table.MarkOccupied()
		s.tableRepo.Update(ctx, table)
	}

	return order, nil
}

// AddItem adds an item to the order
func (s *GuestOrderService) AddItem(ctx context.Context, orderID, productID, productName string, quantity int, unitPrice float64, notes string) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}

	if !order.IsPending() {
		return nil, errors.NewValidationError("hanya order pending yang bisa ditambahkan")
	}

	// Validate stock
	product, err := s.inventoryRepo.GetByID(ctx, productID)
	if err != nil || product == nil {
		return nil, errors.NewNotFoundError("produk", "id", productID)
	}

	if product.Quantity() < quantity {
		return nil, errors.NewValidationError("stok tidak mencukupi untuk %s", product.Name())
	}

	if err := order.AddItem(productID, productName, quantity, unitPrice, notes); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal menambahkan item ke order")
	}

	return order, nil
}

// RemoveItem removes an item from the order
func (s *GuestOrderService) RemoveItem(ctx context.Context, orderID, productID string) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}

	if !order.IsPending() {
		return nil, errors.NewValidationError("hanya order pending yang bisa diubah")
	}

	order.RemoveItem(productID)

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal menghapus item dari order")
	}

	return order, nil
}

// UpdateItemQuantity updates item quantity
func (s *GuestOrderService) UpdateItemQuantity(ctx context.Context, orderID, productID string, quantity int) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}

	if !order.IsPending() {
		return nil, errors.NewValidationError("hanya order pending yang bisa diubah")
	}

	if err := order.UpdateItemQuantity(productID, quantity); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal mengupdate quantity")
	}

	return order, nil
}

// ProcessCheckout processes the order checkout
func (s *GuestOrderService) ProcessCheckout(ctx context.Context, orderID string, paymentMethod model.PaymentMethod, paymentAmount float64) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}

	if !order.IsPending() {
		return nil, errors.NewValidationError("hanya order pending yang bisa di-checkout")
	}

	if len(order.Items()) == 0 {
		return nil, errors.NewValidationError("order kosong")
	}

	// Deduct inventory
	for _, item := range order.Items() {
		product, err := s.inventoryRepo.GetByID(ctx, item.ProductID())
		if err != nil || product == nil {
			return nil, errors.NewNotFoundError("produk", "id", item.ProductID())
		}

		newQty := product.Quantity() - item.Quantity()
		if newQty < 0 {
			return nil, errors.NewValidationError("stok tidak mencukupi untuk %s", product.Name())
		}

		if err := s.inventoryRepo.UpdateQuantity(ctx, product.ID(), newQty); err != nil {
			return nil, errors.NewInternalError("gagal mengupdate stok")
		}
	}

	// Process payment
	if err := order.ProcessPayment(paymentMethod, paymentAmount); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	// Update order status to confirmed
	order.Confirm()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal menyimpan order")
	}

	return order, nil
}

// CancelOrder cancels an order and restores inventory
func (s *GuestOrderService) CancelOrder(ctx context.Context, orderID string) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}

	if order.IsCompleted() {
		return nil, errors.NewValidationError("order sudah selesai, tidak bisa dibatalkan")
	}

	// Restore inventory
	for _, item := range order.Items() {
		product, err := s.inventoryRepo.GetByID(ctx, item.ProductID())
		if err == nil && product != nil {
			newQty := product.Quantity() + item.Quantity()
			s.inventoryRepo.UpdateQuantity(ctx, product.ID(), newQty)
		}
	}

	order.Cancel()
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal membatalkan order")
	}

	// Mark table as available
	table, err := s.tableRepo.GetByID(ctx, order.TableID())
	if err == nil && table != nil && table.IsOccupied() {
		table.MarkAvailable()
		s.tableRepo.Update(ctx, table)
	}

	return order, nil
}

// UpdateOrderStatus updates order status (for staff)
func (s *GuestOrderService) UpdateOrderStatus(ctx context.Context, orderID string, status model.GuestOrderStatus) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}

	var updateErr error
	switch status {
	case model.OrderConfirmed:
		updateErr = order.Confirm()
	case model.OrderPreparing:
		updateErr = order.MarkPreparing()
	case model.OrderReady:
		updateErr = order.MarkReady()
	case model.OrderServed:
		updateErr = order.MarkServed()
	case model.OrderCancelled:
		updateErr = order.Cancel()
	default:
		return nil, errors.NewValidationError("status order tidak valid")
	}

	if updateErr != nil {
		return nil, errors.NewValidationError(updateErr.Error())
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, errors.NewInternalError("gagal mengupdate status order")
	}

	// If served or cancelled, mark table as available
	if status == model.OrderServed || status == model.OrderCancelled {
		table, err := s.tableRepo.GetByID(ctx, order.TableID())
		if err == nil && table != nil && table.IsOccupied() {
			table.MarkAvailable()
			s.tableRepo.Update(ctx, table)
		}
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (s *GuestOrderService) GetOrderByID(ctx context.Context, orderID string) (*model.GuestOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, errors.NewNotFoundError("order", "id", orderID)
	}
	return order, nil
}

// GetPendingOrders retrieves all pending orders
func (s *GuestOrderService) GetPendingOrders(ctx context.Context) ([]*model.GuestOrder, error) {
	orders, err := s.orderRepo.GetPendingOrders(ctx)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil order pending")
	}
	return orders, nil
}

// GetActiveOrders retrieves all active orders
func (s *GuestOrderService) GetActiveOrders(ctx context.Context) ([]*model.GuestOrder, error) {
	orders, err := s.orderRepo.GetActiveOrders(ctx)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil order aktif")
	}
	return orders, nil
}

// GetOrdersByTable retrieves orders for a table
func (s *GuestOrderService) GetOrdersByTable(ctx context.Context, tableID string) ([]*model.GuestOrder, error) {
	orders, err := s.orderRepo.GetByTableID(ctx, tableID, 50)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil order untuk meja")
	}
	return orders, nil
}

// GetTodaySales retrieves today's sales summary
func (s *GuestOrderService) GetTodaySales(ctx context.Context) (*GuestOrderSalesSummary, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	orders, err := s.orderRepo.GetByDateRange(ctx, startOfDay, endOfDay)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil data penjualan")
	}

	totalSales := 0.0
	totalOrders := 0
	totalItems := 0

	for _, o := range orders {
		if o.Status() == model.OrderServed {
			totalSales += o.TotalAmount()
			totalOrders++
			totalItems += len(o.Items())
		}
	}

	return &GuestOrderSalesSummary{
		TotalSales:        totalSales,
		TotalOrders:       totalOrders,
		TotalItems:        totalItems,
		Date:              startOfDay.Format("2006-01-02"),
	}, nil
}

// ListOrdersWithPagination retrieves orders with pagination
func (s *GuestOrderService) ListOrdersWithPagination(ctx context.Context, filter repository.GuestOrderFilter) (*repository.PaginatedGuestOrders, error) {
	return s.orderRepo.ListWithPagination(ctx, filter)
}

// GuestOrderSalesSummary represents today's sales summary
type GuestOrderSalesSummary struct {
	TotalSales        float64 `json:"total_sales"`
	TotalOrders       int     `json:"total_orders"`
	TotalItems        int     `json:"total_items"`
	Date              string  `json:"date"`
}
