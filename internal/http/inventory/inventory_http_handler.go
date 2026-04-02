package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"

	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	"github.com/example/jwt-ddd-clean/internal/dto"
)

// InventoryHTTPHandler handles HTTP requests for inventory operations
type InventoryHTTPHandler struct {
	inventoryService *service.InventoryService
}

// NewInventoryHTTPHandler creates a new InventoryHTTPHandler
func NewInventoryHTTPHandler(inventoryService *service.InventoryService) *InventoryHTTPHandler {
	return &InventoryHTTPHandler{
		inventoryService: inventoryService,
	}
}

// CreateInventory handles POST /api/inventory
func (h *InventoryHTTPHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inv := &model.Inventory{
		ID:          generateID(),
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Location:    req.Location,
		MinStock:    req.MinStock,
		MaxStock:    req.MaxStock,
		Price:       req.Price,
	}

	result, err := h.inventoryService.CreateInventory(r.Context(), inv)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := h.toInventoryResponse(result)
	h.sendSuccess(w, "Inventory item created successfully", response, http.StatusCreated)
}

// GetInventory handles GET /api/inventory/{id}
func (h *InventoryHTTPHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		h.sendError(w, apperrors.NewValidationError("id", "is required"))
		return
	}

	result, err := h.inventoryService.GetInventory(r.Context(), id)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := h.toInventoryResponse(result)
	h.sendSuccess(w, "Inventory item retrieved successfully", response, http.StatusOK)
}

// UpdateInventory handles PUT /api/inventory/{id}
func (h *InventoryHTTPHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id := extractIDFromPath(r.URL.Path)

	var req dto.UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
		return
	}

	// Override ID from path if not provided in body
	if req.ID == "" {
		req.ID = id
	}

	inv := &model.Inventory{
		ID:          req.ID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Location:    req.Location,
		MinStock:    req.MinStock,
		MaxStock:    req.MaxStock,
		Price:       req.Price,
	}

	result, err := h.inventoryService.UpdateInventory(r.Context(), inv)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := h.toInventoryResponse(result)
	h.sendSuccess(w, "Inventory item updated successfully", response, http.StatusOK)
}

// DeleteInventory handles DELETE /api/inventory/{id}
func (h *InventoryHTTPHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		h.sendError(w, apperrors.NewValidationError("id", "is required"))
		return
	}

	if err := h.inventoryService.DeleteInventory(r.Context(), id); err != nil {
		h.sendError(w, err)
		return
	}

	h.sendSuccess(w, "Inventory item deleted successfully", nil, http.StatusOK)
}

// ListInventory handles GET /api/inventory
func (h *InventoryHTTPHandler) ListInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	filter := &model.InventoryFilter{}

	if sku := r.URL.Query().Get("sku"); sku != "" {
		filter.SKU = &sku
	}
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = &name
	}
	if location := r.URL.Query().Get("location"); location != "" {
		filter.Location = &location
	}
	if minQty := r.URL.Query().Get("min_qty"); minQty != "" {
		if val, err := strconv.Atoi(minQty); err == nil {
			filter.MinQty = &val
		}
	}
	if maxQty := r.URL.Query().Get("max_qty"); maxQty != "" {
		if val, err := strconv.Atoi(maxQty); err == nil {
			filter.MaxQty = &val
		}
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			filter.Limit = val
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			filter.Offset = val
		}
	}

	result, err := h.inventoryService.ListInventory(r.Context(), filter)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := &dto.InventoryListResponse{
		Items:      h.toInventoryResponseList(result.Items),
		Total:      result.Total,
		Limit:      result.Limit,
		Offset:     result.Offset,
		TotalPages: result.TotalPages,
	}

	h.sendSuccess(w, "Inventory items retrieved successfully", response, http.StatusOK)
}

// UpdateStock handles PUT /api/inventory/{id}/stock
func (h *InventoryHTTPHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		h.sendError(w, apperrors.NewValidationError("id", "is required"))
		return
	}

	var req dto.UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
		return
	}

	// Get current inventory for response
	current, _ := h.inventoryService.GetInventory(r.Context(), id)
	previousQty := 0
	if current != nil {
		previousQty = current.Quantity
	}

	result, err := h.inventoryService.UpdateStock(r.Context(), id, req.Quantity)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := &dto.StockUpdateResponse{
		ID:          result.ID,
		SKU:         result.SKU,
		Name:        result.Name,
		Quantity:    result.Quantity,
		PreviousQty: previousQty,
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	h.sendSuccess(w, "Stock quantity updated successfully", response, http.StatusOK)
}

// AdjustStock handles POST /api/inventory/{id}/stock/adjust
func (h *InventoryHTTPHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		h.sendError(w, apperrors.NewValidationError("id", "is required"))
		return
	}

	var req dto.AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.ErrValidationErr.WithDetails("Invalid request body"))
		return
	}

	// Get current inventory for response
	current, _ := h.inventoryService.GetInventory(r.Context(), id)
	previousQty := 0
	if current != nil {
		previousQty = current.Quantity
	}

	result, err := h.inventoryService.AdjustStock(r.Context(), id, req.Adjustment)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := &dto.StockUpdateResponse{
		ID:          result.ID,
		SKU:         result.SKU,
		Name:        result.Name,
		Quantity:    result.Quantity,
		PreviousQty: previousQty,
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	h.sendSuccess(w, "Stock quantity adjusted successfully", response, http.StatusOK)
}

func (h *InventoryHTTPHandler) toInventoryResponse(inv *model.Inventory) *dto.InventoryResponse {
	return &dto.InventoryResponse{
		ID:          inv.ID,
		SKU:         inv.SKU,
		Name:        inv.Name,
		Description: inv.Description,
		Quantity:    inv.Quantity,
		Unit:        inv.Unit,
		Location:    inv.Location,
		MinStock:    inv.MinStock,
		MaxStock:    inv.MaxStock,
		Price:       inv.Price,
		CreatedAt:   inv.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   inv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *InventoryHTTPHandler) toInventoryResponseList(inventories []*model.Inventory) []*dto.InventoryResponse {
	result := make([]*dto.InventoryResponse, len(inventories))
	for i, inv := range inventories {
		result[i] = h.toInventoryResponse(inv)
	}
	return result
}

func (h *InventoryHTTPHandler) sendSuccess(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := apperrors.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *InventoryHTTPHandler) sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		w.WriteHeader(appErr.GetHTTPStatus())
		json.NewEncoder(w).Encode(appErr.ToResponse())
		return
	}

	// Fallback for non-AppError
	w.WriteHeader(http.StatusInternalServerError)
	response := apperrors.ErrorResponse{
		Success: false,
		Error: apperrors.ErrorDetail{
			Code:    string(apperrors.ErrInternal),
			Message: "An unexpected error occurred",
			Details: err.Error(),
		},
	}
	json.NewEncoder(w).Encode(response)
}

// generateID generates a new unique ID
func generateID() string {
	return uuid.New().String()
}

// extractIDFromPath extracts the ID from a URL path like /api/inventory/{id}
func extractIDFromPath(path string) string {
	// Remove trailing slash
	path = strings.TrimSuffix(path, "/")
	// Split by /
	parts := strings.Split(path, "/")
	// Return the last part
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
