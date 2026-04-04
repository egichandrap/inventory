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

// TableHandler handles table management HTTP requests
type TableHandler struct {
	tableUsecase usecase.TableUsecase
}

// NewTableHandler creates a new TableHandler
func NewTableHandler(tableUsecase usecase.TableUsecase) *TableHandler {
	return &TableHandler{
		tableUsecase: tableUsecase,
	}
}

// CreateTable handles POST /api/tables
func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	// Validate required fields
	if req.Number <= 0 || req.Capacity <= 0 {
		h.sendError(w, apperrors.NewValidationError("nomor meja dan capacity harus lebih dari 0"))
		return
	}

	table, err := h.tableUsecase.CreateTable(r.Context(), req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusCreated, true, "Meja berhasil dibuat", table)
}

// GetTable handles GET /api/tables/{id}
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["id"]

	table, err := h.tableUsecase.GetTable(r.Context(), tableID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Detail meja", table)
}

// UpdateTable handles PUT /api/tables/{id}
func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["id"]

	var req dto.UpdateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	table, err := h.tableUsecase.UpdateTable(r.Context(), tableID, req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Meja berhasil diupdate", table)
}

// DeleteTable handles DELETE /api/tables/{id}
func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["id"]

	if err := h.tableUsecase.DeleteTable(r.Context(), tableID); err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Meja berhasil dihapus", nil)
}

// ListTables handles GET /api/tables
func (h *TableHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	filter := repository.TableFilter{
		Limit: 50,
	}

	// Parse query parameters
	location := r.URL.Query().Get("location")
	if location != "" {
		loc := model.TableLocation(location)
		filter.Location = &loc
	}

	status := r.URL.Query().Get("status")
	if status != "" {
		st := model.TableStatus(status)
		filter.Status = &st
	}

	minCapacity := r.URL.Query().Get("min_capacity")
	if minCapacity != "" {
		if val, err := strconv.Atoi(minCapacity); err == nil {
			filter.MinCapacity = val
		}
	}

	maxCapacity := r.URL.Query().Get("max_capacity")
	if maxCapacity != "" {
		if val, err := strconv.Atoi(maxCapacity); err == nil {
			filter.MaxCapacity = val
		}
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

	tables, err := h.tableUsecase.ListTables(r.Context(), filter)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Daftar meja", tables)
}

// UpdateTableStatus handles POST /api/tables/{id}/status
func (h *TableHandler) UpdateTableStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["id"]

	var req struct {
		Status model.TableStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, apperrors.NewValidationError("request body tidak valid"))
		return
	}

	if req.Status == "" {
		h.sendError(w, apperrors.NewValidationError("status harus diisi"))
		return
	}

	table, err := h.tableUsecase.UpdateTableStatus(r.Context(), tableID, req.Status)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Status meja berhasil diupdate", table)
}

// GenerateQRCode handles POST /api/tables/{id}/qr
func (h *TableHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["id"]

	qrCode, err := h.tableUsecase.GenerateQRCode(r.Context(), tableID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	response := map[string]string{
		"qr_code": qrCode,
	}

	h.sendJSON(w, http.StatusOK, true, "QR code berhasil digenerate", response)
}

// GetAvailableTables handles GET /api/tables/available
func (h *TableHandler) GetAvailableTables(w http.ResponseWriter, r *http.Request) {
	var location *model.TableLocation
	loc := r.URL.Query().Get("location")
	if loc != "" {
		tableLoc := model.TableLocation(loc)
		location = &tableLoc
	}

	tables, err := h.tableUsecase.GetAvailableTables(r.Context(), location)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, true, "Meja tersedia", tables)
}

// Helper methods for JSON response
func (h *TableHandler) sendJSON(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *TableHandler) sendError(w http.ResponseWriter, err error) {
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
