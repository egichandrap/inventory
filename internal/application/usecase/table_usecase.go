package usecase

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/application/dto"
	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
	"github.com/example/jwt-ddd-clean/internal/domain/service"
)

// TableUsecase defines the table usecase interface
type TableUsecase interface {
	CreateTable(ctx context.Context, req dto.CreateTableRequest) (*dto.TableResponse, error)
	GetTable(ctx context.Context, id string) (*dto.TableResponse, error)
	UpdateTable(ctx context.Context, id string, req dto.UpdateTableRequest) (*dto.TableResponse, error)
	DeleteTable(ctx context.Context, id string) error
	ListTables(ctx context.Context, filter repository.TableFilter) (*dto.TableListResponse, error)
	UpdateTableStatus(ctx context.Context, id string, status model.TableStatus) (*dto.TableResponse, error)
	GenerateQRCode(ctx context.Context, id string) (string, error)
	GetAvailableTables(ctx context.Context, location *model.TableLocation) ([]dto.TableResponse, error)
	GetTableCount(ctx context.Context, filter repository.TableFilter) (int64, error)
}

type tableUsecase struct {
	tableService *service.TableService
}

// NewTableUsecase creates a new TableUsecase
func NewTableUsecase(tableService *service.TableService) TableUsecase {
	return &tableUsecase{
		tableService: tableService,
	}
}

func (u *tableUsecase) CreateTable(ctx context.Context, req dto.CreateTableRequest) (*dto.TableResponse, error) {
	table, err := u.tableService.Create(ctx, req.Number, req.Location, req.Capacity, req.Description)
	if err != nil {
		return nil, err
	}

	resp := dto.ToTableResponse(table)
	return &resp, nil
}

func (u *tableUsecase) GetTable(ctx context.Context, id string) (*dto.TableResponse, error) {
	table, err := u.tableService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToTableResponse(table)
	return &resp, nil
}

func (u *tableUsecase) UpdateTable(ctx context.Context, id string, req dto.UpdateTableRequest) (*dto.TableResponse, error) {
	table, err := u.tableService.Update(ctx, id, req.Location, req.Capacity, req.Description)
	if err != nil {
		return nil, err
	}

	resp := dto.ToTableResponse(table)
	return &resp, nil
}

func (u *tableUsecase) DeleteTable(ctx context.Context, id string) error {
	return u.tableService.Delete(ctx, id)
}

func (u *tableUsecase) ListTables(ctx context.Context, filter repository.TableFilter) (*dto.TableListResponse, error) {
	tables, err := u.tableService.List(ctx, &filter)
	if err != nil {
		return nil, err
	}

	total, err := u.tableService.Count(ctx, &filter)
	if err != nil {
		return nil, err
	}

	resp := dto.ToTableListResponse(tables, total, filter.Limit, filter.Offset)
	return &resp, nil
}

func (u *tableUsecase) UpdateTableStatus(ctx context.Context, id string, status model.TableStatus) (*dto.TableResponse, error) {
	table, err := u.tableService.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	resp := dto.ToTableResponse(table)
	return &resp, nil
}

func (u *tableUsecase) GenerateQRCode(ctx context.Context, id string) (string, error) {
	return u.tableService.GenerateQR(ctx, id)
}

func (u *tableUsecase) GetAvailableTables(ctx context.Context, location *model.TableLocation) ([]dto.TableResponse, error) {
	tables, err := u.tableService.GetAvailableTables(ctx, location)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TableResponse, len(tables))
	for i, t := range tables {
		responses[i] = dto.ToTableResponse(t)
	}

	return responses, nil
}

func (u *tableUsecase) GetTableCount(ctx context.Context, filter repository.TableFilter) (int64, error) {
	return u.tableService.Count(ctx, &filter)
}
