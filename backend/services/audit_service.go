package services

import (
	"context"
	"time"

	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type AuditService interface {
	Log(ctx context.Context, action, entityType string, entityID *uuid.UUID, userID *uuid.UUID, description string, metadata map[string]interface{}, ipAddress, userAgent string) error
	GetAll(ctx context.Context, page, limit int, userID *uuid.UUID, action *string) ([]*models.AuditLog, error)
}

type auditService struct {
	auditRepo repositories.AuditRepository
}

func NewAuditService(auditRepo repositories.AuditRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

func (s *auditService) Log(ctx context.Context, action, entityType string, entityID *uuid.UUID, userID *uuid.UUID, description string, metadata map[string]interface{}, ipAddress, userAgent string) error {
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	auditLog := &models.AuditLog{
		ID:          uuid.New(),
		UserID:      userID,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: descPtr,
		Metadata:    metadata,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		CreatedAt:   time.Now(),
	}

	return s.auditRepo.Create(ctx, auditLog)
}

func (s *auditService) GetAll(ctx context.Context, page, limit int, userID *uuid.UUID, action *string) ([]*models.AuditLog, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	return s.auditRepo.GetAll(ctx, limit, offset, userID, action)
}
