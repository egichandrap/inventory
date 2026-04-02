package service_test

import (
	"context"
	"testing"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	apperrors "github.com/example/jwt-ddd-clean/internal/pkg/errors"
	"github.com/example/jwt-ddd-clean/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestInventoryService_CreateInventory(t *testing.T) {
	t.Run("should create inventory item successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
			Price:    99.99,
		}

		// Act
		result, err := inventoryService.CreateInventory(context.Background(), inv)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "SKU-001", result.SKU)
		assert.Equal(t, "Test Product", result.Name)
	})

	t.Run("should return error when SKU is empty", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:    "inv-001",
			Name:  "Test Product",
			Unit:  "unit",
		}

		// Act
		result, err := inventoryService.CreateInventory(context.Background(), inv)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, apperrors.ErrInvalidFieldErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:  "inv-001",
			SKU: "SKU-001",
			Unit: "unit",
		}

		// Act
		result, err := inventoryService.CreateInventory(context.Background(), inv)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when unit is empty", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:   "inv-001",
			SKU:  "SKU-001",
			Name: "Test Product",
		}

		// Act
		result, err := inventoryService.CreateInventory(context.Background(), inv)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when quantity is negative", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: -10,
			Unit:     "unit",
		}

		// Act
		result, err := inventoryService.CreateInventory(context.Background(), inv)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when SKU already exists", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Create first item
		inv1 := &model.Inventory{
			ID:       "inv-001",
			SKU:      "SKU-001",
			Name:     "Test Product 1",
			Quantity: 100,
			Unit:     "unit",
		}
		_, err := inventoryService.CreateInventory(context.Background(), inv1)
		assert.NoError(t, err)

		// Try to create duplicate
		inv2 := &model.Inventory{
			ID:       "inv-002",
			SKU:      "SKU-001",
			Name:     "Test Product 2",
			Quantity: 50,
			Unit:     "unit",
		}

		// Act
		result, err := inventoryService.CreateInventory(context.Background(), inv2)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, apperrors.ErrConflictErr)
	})
}

func TestInventoryService_GetInventory(t *testing.T) {
	t.Run("should get inventory item successfully", func(t *testing.T) {
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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.GetInventory(context.Background(), "inv-001")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "inv-001", result.ID)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Act
		result, err := inventoryService.GetInventory(context.Background(), "non-existent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, apperrors.ErrNotFoundErr)
	})

	t.Run("should return error when ID is empty", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Act
		result, err := inventoryService.GetInventory(context.Background(), "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInventoryService_UpdateInventory(t *testing.T) {
	t.Run("should update inventory item successfully", func(t *testing.T) {
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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Update
		inv.Name = "Updated Product"
		inv.Quantity = 200

		// Act
		result, err := inventoryService.UpdateInventory(context.Background(), inv)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Product", result.Name)
		assert.Equal(t, 200, result.Quantity)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			ID:       "non-existent",
			SKU:      "SKU-001",
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}

		// Act
		result, err := inventoryService.UpdateInventory(context.Background(), inv)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, apperrors.ErrNotFoundErr)
	})

	t.Run("should return error when ID is empty", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		inv := &model.Inventory{
			SKU:  "SKU-001",
			Name: "Test Product",
			Unit: "unit",
		}

		// Act
		result, err := inventoryService.UpdateInventory(context.Background(), inv)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInventoryService_DeleteInventory(t *testing.T) {
	t.Run("should delete inventory item successfully", func(t *testing.T) {
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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		err = inventoryService.DeleteInventory(context.Background(), "inv-001")

		// Assert
		assert.NoError(t, err)

		// Verify deletion
		result, _ := inventoryService.GetInventory(context.Background(), "inv-001")
		assert.Nil(t, result)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Act
		err := inventoryService.DeleteInventory(context.Background(), "non-existent")

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperrors.ErrNotFoundErr)
	})
}

func TestInventoryService_ListInventory(t *testing.T) {
	t.Run("should list inventory items successfully", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Create test data
		for i := 1; i <= 5; i++ {
			inv := &model.Inventory{
				ID:       "inv-00" + string(rune('0'+i)),
				SKU:      "SKU-00" + string(rune('0'+i)),
				Name:     "Product " + string(rune('0'+i)),
				Quantity: i * 10,
				Unit:     "unit",
			}
			_, err := inventoryService.CreateInventory(context.Background(), inv)
			assert.NoError(t, err)
		}

		// Act
		result, err := inventoryService.ListInventory(context.Background(), &model.InventoryFilter{
			Limit:  10,
			Offset: 0,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, len(result.Items))
		assert.Equal(t, int64(5), result.Total)
	})

	t.Run("should filter by SKU", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		sku := "SKU-TEST"
		inv := &model.Inventory{
			ID:       "inv-001",
			SKU:      sku,
			Name:     "Test Product",
			Quantity: 100,
			Unit:     "unit",
		}
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.ListInventory(context.Background(), &model.InventoryFilter{
			SKU: &sku,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, sku, result.Items[0].SKU)
	})

	t.Run("should handle empty inventory", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Act
		result, err := inventoryService.ListInventory(context.Background(), nil)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
	})
}

func TestInventoryService_UpdateStock(t *testing.T) {
	t.Run("should update stock quantity successfully", func(t *testing.T) {
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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.UpdateStock(context.Background(), "inv-001", 50)

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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.UpdateStock(context.Background(), "inv-001", -10)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Act
		result, err := inventoryService.UpdateStock(context.Background(), "non-existent", 50)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInventoryService_AdjustStock(t *testing.T) {
	t.Run("should adjust stock quantity positively", func(t *testing.T) {
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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.AdjustStock(context.Background(), "inv-001", 50)

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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.AdjustStock(context.Background(), "inv-001", -30)

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
		_, err := inventoryService.CreateInventory(context.Background(), inv)
		assert.NoError(t, err)

		// Act
		result, err := inventoryService.AdjustStock(context.Background(), "inv-001", -20)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when inventory not found", func(t *testing.T) {
		// Arrange
		repo := repository.NewMemoryInventoryRepository()
		inventoryService := service.NewInventoryService(repo)

		// Act
		result, err := inventoryService.AdjustStock(context.Background(), "non-existent", 50)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
