package service

import (
	"context"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
	"github.com/example/jwt-ddd-clean/internal/domain/repository"
)

// AuditService handles audit logging
type AuditService struct {
	auditRepo repository.AuditLogRepository
}

// NewAuditService creates a new AuditService
func NewAuditService(auditRepo repository.AuditLogRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

// LogCreate logs a create action
func (s *AuditService) LogCreate(
	ctx context.Context,
	userID, userName, entityType, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionCreate,
		entityType, entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogUpdate logs an update action
func (s *AuditService) LogUpdate(
	ctx context.Context,
	userID, userName, entityType, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionUpdate,
		entityType, entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogDelete logs a delete action
func (s *AuditService) LogDelete(
	ctx context.Context,
	userID, userName, entityType, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionDelete,
		entityType, entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogCheckout logs a checkout action
func (s *AuditService) LogCheckout(
	ctx context.Context,
	userID, userName, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionCheckout,
		"TRANSACTION", entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogCancel logs a cancel action
func (s *AuditService) LogCancel(
	ctx context.Context,
	userID, userName, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionCancel,
		"TRANSACTION", entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogRefund logs a refund action
func (s *AuditService) LogRefund(
	ctx context.Context,
	userID, userName, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionRefund,
		"TRANSACTION", entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogAdjustStock logs a stock adjustment action
func (s *AuditService) LogAdjustStock(
	ctx context.Context,
	userID, userName, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		model.ActionAdjustStock,
		"INVENTORY", entityID,
		details,
		ipAddress, userAgent,
		true,
	)
	return s.auditRepo.Create(ctx, auditLog)
}

// LogWithSuccess logs an action with success flag
func (s *AuditService) LogWithSuccess(
	ctx context.Context,
	userID, userName string,
	action model.AuditAction,
	entityType, entityID string,
	details map[string]interface{},
	ipAddress, userAgent string,
	success bool,
	errorMsg string,
) error {
	auditLog := model.NewAuditLog(
		userID, userName,
		action,
		entityType, entityID,
		details,
		ipAddress, userAgent,
		success,
	)

	if !success && errorMsg != "" {
		auditLog.SetError(errorMsg)
	}

	return s.auditRepo.Create(ctx, auditLog)
}
