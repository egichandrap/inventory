package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TableStatus represents the status of a table
type TableStatus string

const (
	TableAvailable   TableStatus = "AVAILABLE"
	TableOccupied    TableStatus = "OCCUPIED"
	TableReserved    TableStatus = "RESERVED"
	TableMaintenance TableStatus = "MAINTENANCE"
)

// TableLocation represents the location/area of a table
type TableLocation string

const (
	LocationIndoor  TableLocation = "INDOOR"
	LocationOutdoor TableLocation = "OUTDOOR"
	LocationVIP     TableLocation = "VIP"
	LocationPatio   TableLocation = "PATIO"
)

// Table represents a restaurant table entity
type Table struct {
	id          string
	number      int
	location    TableLocation
	capacity    int
	status      TableStatus
	qrCode      string // URL or token for QR code
	qrGenerated bool
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTable creates a new table entity
func NewTable(number int, location TableLocation, capacity int, description string) (*Table, error) {
	if number <= 0 {
		return nil, fmt.Errorf("nomor meja harus lebih dari 0")
	}
	if capacity <= 0 {
		return nil, fmt.Errorf("kapasitas meja harus lebih dari 0")
	}
	if capacity > 50 {
		return nil, fmt.Errorf("kapasitas meja tidak boleh lebih dari 50")
	}

	now := time.Now()
	return &Table{
		id:          uuid.New().String(),
		number:      number,
		location:    location,
		capacity:    capacity,
		status:      TableAvailable,
		qrCode:      "",
		qrGenerated: false,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructTable recreates a table entity from database
func ReconstructTable(
	id string,
	number int,
	location TableLocation,
	capacity int,
	status TableStatus,
	qrCode string,
	qrGenerated bool,
	description string,
	createdAt, updatedAt time.Time,
) *Table {
	return &Table{
		id:          id,
		number:      number,
		location:    location,
		capacity:    capacity,
		status:      status,
		qrCode:      qrCode,
		qrGenerated: qrGenerated,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Accessors
func (t *Table) ID() string                    { return t.id }
func (t *Table) Number() int                   { return t.number }
func (t *Table) Location() TableLocation       { return t.location }
func (t *Table) Capacity() int                 { return t.capacity }
func (t *Table) Status() TableStatus           { return t.status }
func (t *Table) QRCode() string                { return t.qrCode }
func (t *Table) IsQRGenerated() bool           { return t.qrGenerated }
func (t *Table) Description() string           { return t.description }
func (t *Table) CreatedAt() time.Time          { return t.createdAt }
func (t *Table) UpdatedAt() time.Time          { return t.updatedAt }

// UpdateStatus changes the table status
func (t *Table) UpdateStatus(status TableStatus) error {
	validStatuses := map[TableStatus]bool{
		TableAvailable:   true,
		TableOccupied:    true,
		TableReserved:    true,
		TableMaintenance: true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("status meja tidak valid: %s", status)
	}

	t.status = status
	t.updatedAt = time.Now()
	return nil
}

// UpdateDetails updates table details
func (t *Table) UpdateDetails(location TableLocation, capacity int, description string) error {
	if capacity <= 0 {
		return fmt.Errorf("kapasitas meja harus lebih dari 0")
	}
	if capacity > 50 {
		return fmt.Errorf("kapasitas meja tidak boleh lebih dari 50")
	}

	t.location = location
	t.capacity = capacity
	t.description = description
	t.updatedAt = time.Now()
	return nil
}

// GenerateQRCode marks QR code as generated and stores the code
func (t *Table) GenerateQRCode(qrCode string) {
	t.qrCode = qrCode
	t.qrGenerated = true
	t.updatedAt = time.Now()
}

// ClearQRCode clears the QR code
func (t *Table) ClearQRCode() {
	t.qrCode = ""
	t.qrGenerated = false
	t.updatedAt = time.Now()
}

// IsAvailable checks if table is available
func (t *Table) IsAvailable() bool {
	return t.status == TableAvailable
}

// IsOccupied checks if table is occupied
func (t *Table) IsOccupied() bool {
	return t.status == TableOccupied
}

// MarkOccupied marks table as occupied
func (t *Table) MarkOccupied() error {
	return t.UpdateStatus(TableOccupied)
}

// MarkAvailable marks table as available
func (t *Table) MarkAvailable() error {
	return t.UpdateStatus(TableAvailable)
}

// MarkReserved marks table as reserved
func (t *Table) MarkReserved() error {
	return t.UpdateStatus(TableReserved)
}

// MarkMaintenance marks table as maintenance
func (t *Table) MarkMaintenance() error {
	return t.UpdateStatus(TableMaintenance)
}
