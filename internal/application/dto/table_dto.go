package dto

import (
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// TableResponse represents table response
type TableResponse struct {
	ID          string              `json:"id"`
	Number      int                 `json:"number"`
	Location    model.TableLocation `json:"location"`
	Capacity    int                 `json:"capacity"`
	Status      model.TableStatus   `json:"status"`
	QRCode      string              `json:"qr_code,omitempty"`
	QRGenerated bool                `json:"qr_generated"`
	Description string              `json:"description,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// CreateTableRequest represents create table request
type CreateTableRequest struct {
	Number      int                 `json:"number" validate:"required,min=1"`
	Location    model.TableLocation `json:"location" validate:"required"`
	Capacity    int                 `json:"capacity" validate:"required,min=1,max=50"`
	Description string              `json:"description"`
}

// UpdateTableRequest represents update table request
type UpdateTableRequest struct {
	Location    model.TableLocation `json:"location"`
	Capacity    int                 `json:"capacity" validate:"min=1,max=50"`
	Description string              `json:"description"`
}

// TableListResponse represents paginated table list
type TableListResponse struct {
	Tables     []TableResponse `json:"tables"`
	Total      int64           `json:"total"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
	TotalPages int             `json:"total_pages"`
}

// ToTableResponse converts domain Table to DTO
func ToTableResponse(table *model.Table) TableResponse {
	return TableResponse{
		ID:          table.ID(),
		Number:      table.Number(),
		Location:    table.Location(),
		Capacity:    table.Capacity(),
		Status:      table.Status(),
		QRCode:      table.QRCode(),
		QRGenerated: table.IsQRGenerated(),
		Description: table.Description(),
		CreatedAt:   table.CreatedAt(),
		UpdatedAt:   table.UpdatedAt(),
	}
}

// ToTableListResponse converts domain Tables to DTO
func ToTableListResponse(tables []*model.Table, total int64, limit, offset int) TableListResponse {
	responses := make([]TableResponse, len(tables))
	for i, t := range tables {
		responses[i] = ToTableResponse(t)
	}

	totalPages := int(total) / limit
	if limit > 0 && int(total)%limit > 0 {
		totalPages++
	}

	return TableListResponse{
		Tables:     responses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}
}
