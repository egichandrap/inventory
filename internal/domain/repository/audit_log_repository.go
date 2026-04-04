package repository

import (
	"context"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// AuditLogRepository defines the interface for audit log data operations
type AuditLogRepository interface {
	// Create creates a new audit log
	Create(ctx context.Context, auditLog *model.AuditLog) error

	// GetByID retrieves an audit log by ID
	GetByID(ctx context.Context, id string) (*model.AuditLog, error)

	// List retrieves audit logs with filtering
	List(ctx context.Context, filter AuditLogFilter) ([]*model.AuditLog, error)

	// Count returns total number of audit logs
	Count(ctx context.Context, filter AuditLogFilter) (int64, error)
}

// AuditLogFilter defines filter options for listing audit logs
type AuditLogFilter struct {
	UserID      string
	Action      model.AuditAction
	EntityType  string
	EntityID    string
	StartDate   time.Time
	EndDate     time.Time
	SuccessOnly *bool
	Limit       int
	Offset      int
}
