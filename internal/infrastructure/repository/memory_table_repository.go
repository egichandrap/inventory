package repository

import (
	"context"
	"strings"
	"sync"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// MemoryTableRepository is an in-memory implementation for testing
type MemoryTableRepository struct {
	mu    sync.RWMutex
	tables map[string]*model.Table
}

// NewMemoryTableRepository creates a new in-memory table repository
func NewMemoryTableRepository() repository.TableRepository {
	return &MemoryTableRepository{
		tables: make(map[string]*model.Table),
	}
}

func (r *MemoryTableRepository) Create(ctx context.Context, table *model.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tables[table.ID()] = table
	return nil
}

func (r *MemoryTableRepository) GetByID(ctx context.Context, id string) (*model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if table, ok := r.tables[id]; ok {
		return table, nil
	}
	return nil, nil
}

func (r *MemoryTableRepository) GetByNumber(ctx context.Context, number int) (*model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, table := range r.tables {
		if table.Number() == number {
			return table, nil
		}
	}
	return nil, nil
}

func (r *MemoryTableRepository) Update(ctx context.Context, table *model.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tables[table.ID()] = table
	return nil
}

func (r *MemoryTableRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tables, id)
	return nil
}

func (r *MemoryTableRepository) List(ctx context.Context, filter *repository.TableFilter) ([]*model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.Table
	for _, table := range r.tables {
		if filter != nil {
			if filter.Location != nil && table.Location() != *filter.Location {
				continue
			}
			if filter.Status != nil && table.Status() != *filter.Status {
				continue
			}
			if filter.MinCapacity > 0 && table.Capacity() < filter.MinCapacity {
				continue
			}
			if filter.MaxCapacity > 0 && table.Capacity() > filter.MaxCapacity {
				continue
			}
			if filter.Search != "" && !strings.Contains(table.Description(), filter.Search) {
				continue
			}
		}
		result = append(result, table)
	}

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		start := filter.Offset
		end := start + filter.Limit

		if start >= len(result) {
			return []*model.Table{}, nil
		}
		if end > len(result) {
			end = len(result)
		}

		result = result[start:end]
	}

	return result, nil
}

func (r *MemoryTableRepository) Count(ctx context.Context, filter *repository.TableFilter) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, table := range r.tables {
		if filter != nil {
			if filter.Location != nil && table.Location() != *filter.Location {
				continue
			}
			if filter.Status != nil && table.Status() != *filter.Status {
				continue
			}
			if filter.MinCapacity > 0 && table.Capacity() < filter.MinCapacity {
				continue
			}
			if filter.MaxCapacity > 0 && table.Capacity() > filter.MaxCapacity {
				continue
			}
			if filter.Search != "" && !strings.Contains(table.Description(), filter.Search) {
				continue
			}
		}
		count++
	}
	return int64(count), nil
}

func (r *MemoryTableRepository) GetAvailableTables(ctx context.Context, location *model.TableLocation) ([]*model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.Table
	for _, table := range r.tables {
		if table.IsAvailable() {
			if location != nil && table.Location() != *location {
				continue
			}
			result = append(result, table)
		}
	}
	return result, nil
}

func (r *MemoryTableRepository) ExistsByNumber(ctx context.Context, number int, excludeID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for id, table := range r.tables {
		if table.Number() == number && id != excludeID {
			return true, nil
		}
	}
	return false, nil
}

// Ensure it implements the interface
var _ repository.TableRepository = (*MemoryTableRepository)(nil)
