package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// MemoryGuestOrderRepository is an in-memory implementation for testing
type MemoryGuestOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*model.GuestOrder
	seq    int
}

// NewMemoryGuestOrderRepository creates a new in-memory guest order repository
func NewMemoryGuestOrderRepository() repository.GuestOrderRepository {
	return &MemoryGuestOrderRepository{
		orders: make(map[string]*model.GuestOrder),
		seq:    0,
	}
}

func (r *MemoryGuestOrderRepository) Create(ctx context.Context, order *model.GuestOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID()] = order
	return nil
}

func (r *MemoryGuestOrderRepository) GetByID(ctx context.Context, id string) (*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if order, ok := r.orders[id]; ok {
		return order, nil
	}
	return nil, nil
}

func (r *MemoryGuestOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, order := range r.orders {
		if order.OrderNumber() == orderNumber {
			return order, nil
		}
	}
	return nil, nil
}

func (r *MemoryGuestOrderRepository) Update(ctx context.Context, order *model.GuestOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID()] = order
	return nil
}

func (r *MemoryGuestOrderRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.orders, id)
	return nil
}

func (r *MemoryGuestOrderRepository) List(ctx context.Context, filter repository.GuestOrderFilter) ([]*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.GuestOrder
	for _, order := range r.orders {
		if filter.Status != "" && order.Status() != filter.Status {
			continue
		}
		if filter.PaymentStatus != "" && order.PaymentStatus() != filter.PaymentStatus {
			continue
		}
		if filter.TableID != "" && order.TableID() != filter.TableID {
			continue
		}
		if !filter.StartDate.IsZero() && order.CreatedAt().Before(filter.StartDate) {
			continue
		}
		if !filter.EndDate.IsZero() && order.CreatedAt().After(filter.EndDate) {
			continue
		}
		if filter.Search != "" && !strings.Contains(order.OrderNumber(), filter.Search) && !strings.Contains(order.CustomerName(), filter.Search) {
			continue
		}
		result = append(result, order)
	}

	// Apply pagination
	if filter.Limit > 0 {
		start := filter.Offset
		end := start + filter.Limit

		if start >= len(result) {
			return []*model.GuestOrder{}, nil
		}
		if end > len(result) {
			end = len(result)
		}

		result = result[start:end]
	}

	return result, nil
}

func (r *MemoryGuestOrderRepository) Count(ctx context.Context, filter repository.GuestOrderFilter) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, order := range r.orders {
		if filter.Status != "" && order.Status() != filter.Status {
			continue
		}
		if filter.TableID != "" && order.TableID() != filter.TableID {
			continue
		}
		if !filter.StartDate.IsZero() && order.CreatedAt().Before(filter.StartDate) {
			continue
		}
		if !filter.EndDate.IsZero() && order.CreatedAt().After(filter.EndDate) {
			continue
		}
		if filter.Search != "" && !strings.Contains(order.OrderNumber(), filter.Search) && !strings.Contains(order.CustomerName(), filter.Search) {
			continue
		}
		count++
	}
	return int64(count), nil
}

func (r *MemoryGuestOrderRepository) ListWithPagination(ctx context.Context, filter repository.GuestOrderFilter) (*repository.PaginatedGuestOrders, error) {
	orders, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.Limit
	if filter.Limit > 0 && int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &repository.PaginatedGuestOrders{
		Orders:     orders,
		Total:      total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		TotalPages: totalPages,
	}, nil
}

func (r *MemoryGuestOrderRepository) GetByTableID(ctx context.Context, tableID string, limit int) ([]*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.GuestOrder
	for _, order := range r.orders {
		if order.TableID() == tableID {
			result = append(result, order)
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func (r *MemoryGuestOrderRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.GuestOrder
	for _, order := range r.orders {
		if !order.CreatedAt().Before(startDate) && !order.CreatedAt().After(endDate) {
			result = append(result, order)
		}
	}

	return result, nil
}

func (r *MemoryGuestOrderRepository) GetByStatus(ctx context.Context, status model.GuestOrderStatus, limit int) ([]*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.GuestOrder
	for _, order := range r.orders {
		if order.Status() == status {
			result = append(result, order)
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func (r *MemoryGuestOrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%04d", now.Format("20060102"), r.seq), nil
}

func (r *MemoryGuestOrderRepository) GetPendingOrders(ctx context.Context) ([]*model.GuestOrder, error) {
	return r.GetByStatus(ctx, model.OrderPending, 0)
}

func (r *MemoryGuestOrderRepository) GetActiveOrders(ctx context.Context) ([]*model.GuestOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.GuestOrder
	for _, order := range r.orders {
		if order.Status() != model.OrderServed && order.Status() != model.OrderCancelled {
			result = append(result, order)
		}
	}

	return result, nil
}

// Ensure it implements the interface
var _ repository.GuestOrderRepository = (*MemoryGuestOrderRepository)(nil)
