package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	"github.com/example/jwt-ddd-clean/internal/dto"
	inventoryhttp "github.com/example/jwt-ddd-clean/internal/http/inventory"
	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/example/jwt-ddd-clean/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestInventoryHTTPHandler_CreateInventory(t *testing.T) {
	t.Run("should create inventory item successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		reqBody := dto.CreateInventoryRequest{
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
			Price:    99.99,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response apperrors.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("should return error when SKU is missing", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		reqBody := dto.CreateInventoryRequest{
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response apperrors.ErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
	})

	t.Run("should return error when name is missing", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		reqBody := dto.CreateInventoryRequest{
			SKU:      "SKU-001",
			Quantity: 100,
			Unit:     "unit",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when unit is missing", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		reqBody := dto.CreateInventoryRequest{
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestInventoryHTTPHandler_GetInventory(t *testing.T) {
	t.Run("should get inventory item successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		// Create item first
		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/inv-001", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response apperrors.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/non-existent", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response apperrors.ErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response.Success)
		assert.Equal(t, "ERR_NOT_FOUND", response.Error.Code)
	})
}

func TestInventoryHTTPHandler_ListInventory(t *testing.T) {
	t.Run("should list inventory items successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		// Create test data
		for i := 1; i <= 3; i++ {
			inv := &model.Inventory{
				ID:       "inv-00" + string(rune('0'+i)),
				SKU:      "SKU-00" + string(rune('0'+i)),
				Name:     "Product " + string(rune('0'+i)),
				Quantity: i * 10,
				Unit:     "unit",
			}
			_, _ = inventoryService.CreateInventory(t.Context(), inv)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/inventory?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		// Act
		handler.ListInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response apperrors.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("should filter by SKU", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-TEST",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory?sku=TEST", nil)
		w := httptest.NewRecorder()

		// Act
		handler.ListInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestInventoryHTTPHandler_UpdateInventory(t *testing.T) {
	t.Run("should update inventory item successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		// Create item first
		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		reqBody := dto.UpdateInventoryRequest{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Updated Product",
			Quantity: 200,
			Unit:     "unit",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/inventory/inv-001", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.UpdateInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response apperrors.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		reqBody := dto.UpdateInventoryRequest{
			ID:   "non-existent",
			SKU:  "SKU-001",
			Name: "Test Product",
			Unit: "unit",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/inventory/non-existent", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.UpdateInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestInventoryHTTPHandler_DeleteInventory(t *testing.T) {
	t.Run("should delete inventory item successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		// Create item first
		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		req := httptest.NewRequest(http.MethodDelete, "/api/inventory/inv-001", nil)
		w := httptest.NewRecorder()

		// Act
		handler.DeleteInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response apperrors.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)
		handler := inventoryhttp.NewInventoryHTTPHandler(inventoryService)

		req := httptest.NewRequest(http.MethodDelete, "/api/inventory/non-existent", nil)
		w := httptest.NewRecorder()

		// Act
		handler.DeleteInventory(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestInventoryHTTPHandler_UpdateStock(t *testing.T) {
	t.Run("should update stock quantity successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Create item first
		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		// Act - call service directly for unit test
		result, err := inventoryService.UpdateStock(t.Context(), "inv-001", 50)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 50, result.Quantity)
	})

	t.Run("should return error when quantity is negative", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		// Act - call service directly for unit test
		result, err := inventoryService.UpdateStock(t.Context(), "inv-001", -10)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInventoryHTTPHandler_AdjustStock(t *testing.T) {
	t.Run("should adjust stock quantity positively", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Create item first
		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		// Act - call service directly for unit test
		result, err := inventoryService.AdjustStock(t.Context(), "inv-001", 50)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 150, result.Quantity)
	})

	t.Run("should adjust stock quantity negatively", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		// Act - call service directly for unit test
		result, err := inventoryService.AdjustStock(t.Context(), "inv-001", -30)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 70, result.Quantity)
	})

	t.Run("should return error when adjustment results in negative stock", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 10,
			Unit:     "unit",
		}
		_, _ = inventoryService.CreateInventory(t.Context(), inv)

		// Act - call service directly for unit test
		result, err := inventoryService.AdjustStock(t.Context(), "inv-001", -20)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInventoryResponse_DTO(t *testing.T) {
	t.Run("should convert inventory to response correctly", func(t *testing.T) {
		// Arrange
		now := time.Now()
		inv := &model.Inventory{
			ID:          "inv-001",
			SKU:         "SKU-001",
			Name:        "Test Product",
			Description: "Test Description",
			Quantity:    100,
			Unit:        "unit",
			Location:    "Warehouse A",
			MinStock:    10,
			MaxStock:    200,
			Price:       99.99,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Act
		response := &dto.InventoryResponse{
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

		// Assert
		assert.Equal(t, "inv-001", response.ID)
		assert.Equal(t, "SKU-001", response.SKU)
		assert.Equal(t, "Test Product", response.Name)
		assert.Equal(t, 100, response.Quantity)
		assert.Equal(t, 99.99, response.Price)
	})

	t.Run("should marshal inventory list response correctly", func(t *testing.T) {
		// Arrange
		response := &dto.InventoryListResponse{
			Items: []*dto.InventoryResponse{
				{
					ID:       "inv-001",
					SKU:      "SKU-001",
					Name:     "Product 1",
					Quantity: 100,
					Unit:     "unit",
				},
			},
			Total:      1,
			Limit:      10,
			Offset:     0,
			TotalPages: 1,
		}

		// Act
		data, err := json.Marshal(response)

		// Assert
		assert.NoError(t, err)
		assert.Contains(t, string(data), "inv-001")
		assert.Contains(t, string(data), "SKU-001")
	})
}
