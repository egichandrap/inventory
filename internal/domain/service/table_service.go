package service

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	"github.com/example/jwt-ddd-clean/internal/pkg/errors"
)

// TableService handles table management business logic
type TableService struct {
	tableRepo repository.TableRepository
	qrService *QRCodeService
}

// NewTableService creates a new TableService
func NewTableService(tableRepo repository.TableRepository, qrService *QRCodeService) *TableService {
	return &TableService{
		tableRepo: tableRepo,
		qrService: qrService,
	}
}

// Create creates a new table
func (s *TableService) Create(ctx context.Context, number int, location model.TableLocation, capacity int, description string) (*model.Table, error) {
	// Check if number already exists
	exists, err := s.tableRepo.ExistsByNumber(ctx, number, "")
	if err != nil {
		return nil, errors.NewInternalError("gagal memeriksa nomor meja")
	}
	if exists {
		return nil, errors.NewValidationError("nomor meja sudah digunakan")
	}

	// Create table
	table, err := model.NewTable(number, location, capacity, description)
	if err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	if err := s.tableRepo.Create(ctx, table); err != nil {
		return nil, errors.NewInternalError("gagal membuat meja: %v", err)
	}

	return table, nil
}

// GetByID retrieves a table by ID
func (s *TableService) GetByID(ctx context.Context, id string) (*model.Table, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil data meja")
	}
	if table == nil {
		return nil, errors.NewNotFoundError("meja", "id", id)
	}
	return table, nil
}

// Update updates table details
func (s *TableService) Update(ctx context.Context, id string, location model.TableLocation, capacity int, description string) (*model.Table, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil || table == nil {
		return nil, errors.NewNotFoundError("meja", "id", id)
	}

	if err := table.UpdateDetails(location, capacity, description); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	if err := s.tableRepo.Update(ctx, table); err != nil {
		return nil, errors.NewInternalError("gagal mengupdate meja")
	}

	return table, nil
}

// Delete deletes a table
func (s *TableService) Delete(ctx context.Context, id string) error {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil || table == nil {
		return errors.NewNotFoundError("meja", "id", id)
	}

	if table.IsOccupied() {
		return errors.NewValidationError("meja sedang ditempati, tidak bisa dihapus")
	}

	if err := s.tableRepo.Delete(ctx, id); err != nil {
		return errors.NewInternalError("gagal menghapus meja")
	}

	return nil
}

// List retrieves all tables with optional filtering
func (s *TableService) List(ctx context.Context, filter *repository.TableFilter) ([]*model.Table, error) {
	tables, err := s.tableRepo.List(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil data meja")
	}
	return tables, nil
}

// UpdateStatus updates table status
func (s *TableService) UpdateStatus(ctx context.Context, id string, status model.TableStatus) (*model.Table, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil || table == nil {
		return nil, errors.NewNotFoundError("meja", "id", id)
	}

	if err := table.UpdateStatus(status); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	if err := s.tableRepo.Update(ctx, table); err != nil {
		return nil, errors.NewInternalError("gagal mengupdate status meja")
	}

	return table, nil
}

// GenerateQR generates QR code for a table
func (s *TableService) GenerateQR(ctx context.Context, id string) (string, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil || table == nil {
		return "", errors.NewNotFoundError("meja", "id", id)
	}

	if s.qrService == nil {
		return "", errors.NewInternalError("QR service tidak tersedia")
	}

	// Generate QR code as base64 string
	qrString, err := s.qrService.GenerateQRString(table.Number(), table.ID())
	if err != nil {
		return "", errors.NewInternalError("gagal generate QR code: %v", err)
	}

	// Update table
	table.GenerateQRCode(qrString)
	if err := s.tableRepo.Update(ctx, table); err != nil {
		return "", errors.NewInternalError("gagal menyimpan QR code")
	}

	return qrString, nil
}

// GetAvailableTables retrieves all available tables
func (s *TableService) GetAvailableTables(ctx context.Context, location *model.TableLocation) ([]*model.Table, error) {
	tables, err := s.tableRepo.GetAvailableTables(ctx, location)
	if err != nil {
		return nil, errors.NewInternalError("gagal mengambil meja tersedia")
	}
	return tables, nil
}

// Count returns total number of tables
func (s *TableService) Count(ctx context.Context, filter *repository.TableFilter) (int64, error) {
	count, err := s.tableRepo.Count(ctx, filter)
	if err != nil {
		return 0, errors.NewInternalError("gagal menghitung meja")
	}
	return count, nil
}
