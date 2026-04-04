package model

import (
	"fmt"
	"time"
)

// AlertType represents the type of inventory alert
type AlertType string

const (
	AlertLowStock     AlertType = "LOW_STOCK"
	AlertOutOfStock   AlertType = "OUT_OF_STOCK"
	AlertOverStock    AlertType = "OVER_STOCK"
	AlertExpiringSoon AlertType = "EXPIRING_SOON"
	AlertExpired      AlertType = "EXPIRED"
)

// AlertSeverity represents the severity level
type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "LOW"
	SeverityMedium   AlertSeverity = "MEDIUM"
	SeverityHigh     AlertSeverity = "HIGH"
	SeverityCritical AlertSeverity = "CRITICAL"
)

// InventoryAlert represents an inventory alert
type InventoryAlert struct {
	id          string
	inventoryID string
	itemName    string
	itemSKU     string
	alertType   AlertType
	severity    AlertSeverity
	message     string
	currentQty int
	thresholdQty int
	createdAt   time.Time
	acknowledged bool
	acknowledgedAt *time.Time
	acknowledgedBy string
}

// NewInventoryAlert creates a new inventory alert
func NewInventoryAlert(
	inventoryID, itemName, itemSKU string,
	alertType AlertType,
	severity AlertSeverity,
	message string,
	currentQty, thresholdQty int,
) *InventoryAlert {
	return &InventoryAlert{
		id:           generateAlertID(),
		inventoryID:  inventoryID,
		itemName:     itemName,
		itemSKU:      itemSKU,
		alertType:    alertType,
		severity:     severity,
		message:      message,
		currentQty:   currentQty,
		thresholdQty: thresholdQty,
		createdAt:    time.Now(),
		acknowledged: false,
	}
}

// Acknowledge acknowledges the alert
func (a *InventoryAlert) Acknowledge(userID string) {
	a.acknowledged = true
	now := time.Now()
	a.acknowledgedAt = &now
	a.acknowledgedBy = userID
}

// Accessors
func (a *InventoryAlert) ID() string              { return a.id }
func (a *InventoryAlert) InventoryID() string     { return a.inventoryID }
func (a *InventoryAlert) ItemName() string        { return a.itemName }
func (a *InventoryAlert) ItemSKU() string         { return a.itemSKU }
func (a *InventoryAlert) AlertType() AlertType    { return a.alertType }
func (a *InventoryAlert) Severity() AlertSeverity { return a.severity }
func (a *InventoryAlert) Message() string         { return a.message }
func (a *InventoryAlert) CurrentQty() int         { return a.currentQty }
func (a *InventoryAlert) ThresholdQty() int       { return a.thresholdQty }
func (a *InventoryAlert) CreatedAt() time.Time    { return a.createdAt }
func (a *InventoryAlert) IsAcknowledged() bool    { return a.acknowledged }
func (a *InventoryAlert) AcknowledgedAt() *time.Time { return a.acknowledgedAt }
func (a *InventoryAlert) AcknowledgedBy() string  { return a.acknowledgedBy }

// CheckStockAlerts checks inventory stock levels and returns alerts
func CheckStockAlerts(inventoryID, itemName, itemSKU string, currentQty, minStock, maxStock int) []*InventoryAlert {
	var alerts []*InventoryAlert

	// Check out of stock
	if currentQty == 0 {
		alert := NewInventoryAlert(
			inventoryID, itemName, itemSKU,
			AlertOutOfStock,
			SeverityCritical,
			fmt.Sprintf("Stok %s habis total", itemName),
			currentQty,
			minStock,
		)
		alerts = append(alerts, alert)
	} else if currentQty <= minStock {
		// Check low stock
		alert := NewInventoryAlert(
			inventoryID, itemName, itemSKU,
			AlertLowStock,
			SeverityHigh,
			fmt.Sprintf("Stok %s rendah: %d (minimum: %d)", itemName, currentQty, minStock),
			currentQty,
			minStock,
		)
		alerts = append(alerts, alert)
	}

	// Check over stock
	if maxStock > 0 && currentQty > maxStock {
		alert := NewInventoryAlert(
			inventoryID, itemName, itemSKU,
			AlertOverStock,
			SeverityMedium,
			fmt.Sprintf("Stok %s berlebih: %d (maksimum: %d)", itemName, currentQty, maxStock),
			currentQty,
			maxStock,
		)
		alerts = append(alerts, alert)
	}

	return alerts
}

// generateAlertID generates a unique alert ID
func generateAlertID() string {
	return fmt.Sprintf("ALERT-%d", time.Now().UnixNano())
}
