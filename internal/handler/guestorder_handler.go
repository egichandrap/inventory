package handler

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"strconv"

	"github.com/example/jwt-ddd-clean/internal/application/dto"
	"github.com/example/jwt-ddd-clean/internal/application/usecase"
	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/gorilla/mux"
)

// GuestOrderHandler handles guest order HTTP requests
type GuestOrderHandler struct {
	guestOrderUsecase usecase.GuestOrderUsecase
}

// NewGuestOrderHandler creates a new GuestOrderHandler
func NewGuestOrderHandler(guestOrderUsecase usecase.GuestOrderUsecase) *GuestOrderHandler {
	return &GuestOrderHandler{
		guestOrderUsecase: guestOrderUsecase,
	}
}

// CreateOrder handles POST /api/guest/orders
func (h *GuestOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGuestOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	if req.TableID == "" || req.CustomerName == "" {
		h.sendError(w, apperrors.NewValidationError("table_id dan customer_name harus diisi"))
		return
	}

	order, err := h.guestOrderUsecase.CreateOrder(r.Context(), req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusCreated, true, "Order berhasil dibuat", order)
}

// GetOrder handles GET /api/guest/orders/{id}
func (h *GuestOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, err := h.guestOrderUsecase.GetOrder(r.Context(), orderID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Detail order", order)
}

// AddItem handles POST /api/guest/orders/{id}/items
func (h *GuestOrderHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var req dto.AddOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	if req.ProductID == "" || req.ProductName == "" || req.Quantity < 1 {
		h.sendError(w, apperrors.NewValidationError("product_id, product_name, dan quantity harus diisi"))
		return
	}

	order, err := h.guestOrderUsecase.AddItem(r.Context(), orderID, req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Item berhasil ditambahkan", order)
}

// UpdateItemQuantity handles PUT /api/guest/orders/{id}/items/{productID}
func (h *GuestOrderHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	productID := vars["productID"]

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	if req.Quantity < 0 {
		h.sendError(w, apperrors.NewValidationError("quantity tidak boleh negatif"))
		return
	}

	order, err := h.guestOrderUsecase.UpdateItemQuantity(r.Context(), orderID, productID, req.Quantity)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Quantity item berhasil diupdate", order)
}

// RemoveItem handles DELETE /api/guest/orders/{id}/items/{productID}
func (h *GuestOrderHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	productID := vars["productID"]

	order, err := h.guestOrderUsecase.RemoveItem(r.Context(), orderID, productID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Item berhasil dihapus", order)
}

// Checkout handles POST /api/guest/orders/{id}/checkout
func (h *GuestOrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var req dto.GuestCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	if req.PaymentMethod == "" || req.PaymentAmount <= 0 {
		h.sendError(w, apperrors.NewValidationError("payment_method dan payment_amount harus diisi"))
		return
	}

	order, err := h.guestOrderUsecase.Checkout(r.Context(), orderID, req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Checkout berhasil", order)
}

// CancelOrder handles POST /api/guest/orders/{id}/cancel
func (h *GuestOrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, err := h.guestOrderUsecase.CancelOrder(r.Context(), orderID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Order berhasil dibatalkan", order)
}

// UpdateOrderStatus handles POST /api/orders/{id}/status
func (h *GuestOrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var req struct {
		Status model.GuestOrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	if req.Status == "" {
		h.sendError(w, apperrors.NewValidationError("status harus diisi"))
		return
	}

	order, err := h.guestOrderUsecase.UpdateOrderStatus(r.Context(), orderID, req.Status)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Status order berhasil diupdate", order)
}

// ListOrders handles GET /api/orders
func (h *GuestOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	filter := repository.GuestOrderFilter{
		Limit: 50,
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	if status != "" {
		filter.Status = model.GuestOrderStatus(status)
	}

	paymentStatus := r.URL.Query().Get("payment_status")
	if paymentStatus != "" {
		filter.PaymentStatus = model.GuestOrderPaymentStatus(paymentStatus)
	}

	tableID := r.URL.Query().Get("table_id")
	if tableID != "" {
		filter.TableID = tableID
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		filter.Search = search
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			filter.Limit = val
		}
	}

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			filter.Offset = val
		}
	}

	orders, err := h.guestOrderUsecase.ListOrders(r.Context(), filter)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Daftar order", orders)
}

// GetPendingOrders handles GET /api/orders/pending
func (h *GuestOrderHandler) GetPendingOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.guestOrderUsecase.GetPendingOrders(r.Context())
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Order pending", orders)
}

// GetActiveOrders handles GET /api/orders/active
func (h *GuestOrderHandler) GetActiveOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.guestOrderUsecase.GetActiveOrders(r.Context())
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Order aktif", orders)
}

// GetOrdersByTable handles GET /api/orders/table/{tableID}
func (h *GuestOrderHandler) GetOrdersByTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["tableID"]

	orders, err := h.guestOrderUsecase.GetOrdersByTable(r.Context(), tableID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Order untuk meja", orders)
}

// GetTodaySales handles GET /api/reports/sales/today
func (h *GuestOrderHandler) GetTodaySales(w http.ResponseWriter, r *http.Request) {
	sales, err := h.guestOrderUsecase.GetTodaySales(r.Context())
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Penjualan hari ini", sales)
}

// Helper methods
func (h *GuestOrderHandler) sendJSON(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *GuestOrderHandler) sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		w.WriteHeader(appErr.GetHTTPStatus())
		json.NewEncoder(w).Encode(appErr.ToResponse())
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    "ERR_INTERNAL",
			"message": "An unexpected error occurred",
			"details": err.Error(),
		},
	})
}
