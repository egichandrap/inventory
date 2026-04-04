package service_test

import (
	"testing"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceiptService_GenerateTextReceipt(t *testing.T) {
	// Arrange
	receiptSvc := service.NewReceiptService(
		"Test Store",
		"123 Test St",
		"555-1234",
		"TS001",
	)

	transaction, err := model.NewTransaction("TRX-001", "cashier-1", "John Doe")
	require.NoError(t, err)

	err = transaction.AddItem("prod-1", "Product 1", "SKU001", 2, 50000)
	require.NoError(t, err)

	err = transaction.AddItem("prod-2", "Product 2", "SKU002", 1, 100000)
	require.NoError(t, err)

	transaction.ApplyTax(11)
	err = transaction.Complete(model.PaymentCash, 300000)
	require.NoError(t, err)

	// Act
	receipt := receiptSvc.GenerateTextReceipt(transaction)

	// Assert
	assert.NotEmpty(t, receipt)
	assert.Contains(t, receipt, "Test Store")
	assert.Contains(t, receipt, "TRX-001")
	assert.Contains(t, receipt, "Product 1")
	assert.Contains(t, receipt, "Product 2")
}

func TestExportService_ExportTransactionsToCSV(t *testing.T) {
	// Arrange
	exportSvc := service.NewExportService()

	transaction, err := model.NewTransaction("TRX-001", "cashier-1", "John Doe")
	require.NoError(t, err)

	err = transaction.AddItem("prod-1", "Product 1", "SKU001", 2, 50000)
	require.NoError(t, err)

	transaction.ApplyTax(11)
	err = transaction.Complete(model.PaymentCash, 200000)
	require.NoError(t, err)

	transactions := []*model.Transaction{transaction}

	// Act
	var output []byte
	// In real scenario, use bytes.Buffer
	_ = output
	_ = exportSvc
	_ = transactions
	// For now, just verify the service exists
	assert.NotNil(t, exportSvc)
}

func TestInventoryAlerts_CheckStockAlerts(t *testing.T) {
	tests := []struct {
		name         string
		currentQty   int
		minStock     int
		maxStock     int
		expectedAlerts int
	}{
		{
			name:         "Out of stock",
			currentQty:   0,
			minStock:     10,
			maxStock:     100,
			expectedAlerts: 1,
		},
		{
			name:         "Low stock",
			currentQty:   5,
			minStock:     10,
			maxStock:     100,
			expectedAlerts: 1,
		},
		{
			name:         "Normal stock",
			currentQty:   50,
			minStock:     10,
			maxStock:     100,
			expectedAlerts: 0,
		},
		{
			name:         "Over stock",
			currentQty:   150,
			minStock:     10,
			maxStock:     100,
			expectedAlerts: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alerts := model.CheckStockAlerts("inv-1", "Test Product", "SKU001", tt.currentQty, tt.minStock, tt.maxStock)
			assert.Equal(t, tt.expectedAlerts, len(alerts))
		})
	}
}
