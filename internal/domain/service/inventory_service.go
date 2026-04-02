package service

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
)

// InventoryService handles business logic for inventory operations
type InventoryService struct {
	inventoryRepo repository.InventoryRepository
}

// NewInventoryService creates a new InventoryService
func NewInventoryService(inventoryRepo repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

// CreateInventory creates a new inventory item
func (s *InventoryService) CreateInventory(ctx context.Context, inv *model.Inventory) (*model.Inventory, error) {
	// Validate required fields
	if inv.SKU == "" {
		return nil, apperrors.NewValidationError("sku", "is required")
	}
	if inv.Name == "" {
		return nil, apperrors.NewValidationError("name", "is required")
	}
	if inv.Unit == "" {
		return nil, apperrors.NewValidationError("unit", "is required")
	}

	// Check if SKU already exists
	exists, err := s.inventoryRepo.ExistsBySKU(ctx, inv.SKU, "")
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to check SKU existence", apperrors.ErrInternalErr.GetHTTPStatus())
	}
	if exists {
		return nil, apperrors.New(apperrors.ErrConflict, "SKU already exists", apperrors.ErrConflictErr.GetHTTPStatus()).WithDetails(inv.SKU)
	}

	// Validate quantity
	if inv.Quantity < 0 {
		return nil, apperrors.NewValidationError("quantity", "cannot be negative")
	}

	// Validate min/max stock
	if inv.MinStock < 0 {
		return nil, apperrors.NewValidationError("min_stock", "cannot be negative")
	}
	if inv.MaxStock < 0 {
		return nil, apperrors.NewValidationError("max_stock", "cannot be negative")
	}
	if inv.MinStock > inv.MaxStock && inv.MaxStock > 0 {
		return nil, apperrors.NewValidationError("min_stock", "cannot be greater than max_stock")
	}

	// Validate price
	if inv.Price < 0 {
		return nil, apperrors.NewValidationError("price", "cannot be negative")
	}

	if err := s.inventoryRepo.Create(ctx, inv); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to create inventory item", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	return inv, nil
}

// GetInventory retrieves an inventory item by ID
func (s *InventoryService) GetInventory(ctx context.Context, id string) (*model.Inventory, error) {
	if id == "" {
		return nil, apperrors.NewValidationError("id", "is required")
	}

	inv, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}
	if inv == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus()).WithDetails(id)
	}

	return inv, nil
}

// UpdateInventory updates an existing inventory item
func (s *InventoryService) UpdateInventory(ctx context.Context, inv *model.Inventory) (*model.Inventory, error) {
	if inv.ID == "" {
		return nil, apperrors.NewValidationError("id", "is required")
	}

	// Check if inventory exists
	existing, err := s.inventoryRepo.GetByID(ctx, inv.ID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}
	if existing == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus()).WithDetails(inv.ID)
	}

	// Check if SKU is being changed and if it already exists
	if inv.SKU != existing.SKU {
		exists, err := s.inventoryRepo.ExistsBySKU(ctx, inv.SKU, inv.ID)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to check SKU existence", apperrors.ErrInternalErr.GetHTTPStatus())
		}
		if exists {
			return nil, apperrors.New(apperrors.ErrConflict, "SKU already exists", apperrors.ErrConflictErr.GetHTTPStatus()).WithDetails(inv.SKU)
		}
	}

	// Validate quantity if changed
	if inv.Quantity < 0 {
		return nil, apperrors.NewValidationError("quantity", "cannot be negative")
	}

	if err := s.inventoryRepo.Update(ctx, inv); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to update inventory item", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	return inv, nil
}

// DeleteInventory deletes an inventory item
func (s *InventoryService) DeleteInventory(ctx context.Context, id string) error {
	if id == "" {
		return apperrors.NewValidationError("id", "is required")
	}

	// Check if inventory exists
	existing, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}
	if existing == nil {
		return apperrors.New(apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus()).WithDetails(id)
	}

	if err := s.inventoryRepo.Delete(ctx, id); err != nil {
		return apperrors.Wrap(err, apperrors.ErrInternal, "Failed to delete inventory item", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	return nil
}

// ListInventory retrieves a paginated list of inventory items
func (s *InventoryService) ListInventory(ctx context.Context, filter *model.InventoryFilter) (*model.PaginatedInventory, error) {
	if filter == nil {
		filter = &model.InventoryFilter{}
	}

	// Set default pagination values
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	items, err := s.inventoryRepo.List(ctx, filter)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to list inventory items", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	total, err := s.inventoryRepo.Count(ctx, filter)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to count inventory items", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &model.PaginatedInventory{
		Items:      items,
		Total:      total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		TotalPages: totalPages,
	}, nil
}

// UpdateStock updates the quantity of an inventory item
func (s *InventoryService) UpdateStock(ctx context.Context, id string, quantity int) (*model.Inventory, error) {
	if id == "" {
		return nil, apperrors.NewValidationError("id", "is required")
	}

	// Check if inventory exists
	existing, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}
	if existing == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus()).WithDetails(id)
	}

	if quantity < 0 {
		return nil, apperrors.NewValidationError("quantity", "cannot be negative")
	}

	if err := s.inventoryRepo.UpdateQuantity(ctx, id, quantity); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to update stock quantity", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	// Fetch updated inventory
	updated, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrNotFound, "Failed to fetch updated inventory", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}

	return updated, nil
}

// AdjustStock adjusts the quantity by a given amount (positive or negative)
func (s *InventoryService) AdjustStock(ctx context.Context, id string, adjustment int) (*model.Inventory, error) {
	if id == "" {
		return nil, apperrors.NewValidationError("id", "is required")
	}

	// Check if inventory exists
	existing, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}
	if existing == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "Inventory item not found", apperrors.ErrNotFoundErr.GetHTTPStatus()).WithDetails(id)
	}

	newQuantity := existing.Quantity + adjustment
	if newQuantity < 0 {
		return nil, apperrors.NewValidationError("quantity", "adjustment would result in negative stock")
	}

	if err := s.inventoryRepo.UpdateQuantity(ctx, id, newQuantity); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrInternal, "Failed to adjust stock quantity", apperrors.ErrInternalErr.GetHTTPStatus())
	}

	// Fetch updated inventory
	updated, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrNotFound, "Failed to fetch updated inventory", apperrors.ErrNotFoundErr.GetHTTPStatus())
	}

	return updated, nil
}
