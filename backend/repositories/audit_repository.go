package repositories

import (
	"database/sql"
	"encoding/json"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type AuditRepository interface {
	Create(auditLog *models.AuditLog) error
	GetAll(limit, offset int, userID *uuid.UUID, action *string) ([]*models.AuditLog, error)
}

type auditRepository struct {
	db *sql.DB
}

func NewAuditRepository() AuditRepository {
	return &auditRepository{
		db: database.DB,
	}
}

func (r *auditRepository) Create(auditLog *models.AuditLog) error {
	metadataJSON, _ := json.Marshal(auditLog.Metadata)
	
	query := `
		INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, metadata, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.Exec(
		query,
		auditLog.ID,
		auditLog.UserID,
		auditLog.Action,
		auditLog.EntityType,
		auditLog.EntityID,
		metadataJSON,
		auditLog.IPAddress,
		auditLog.UserAgent,
	)
	
	return err
}

func (r *auditRepository) GetAll(limit, offset int, userID *uuid.UUID, action *string) ([]*models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, entity_type, entity_id, metadata, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1
	
	if userID != nil {
		query += ` AND user_id = $` + string(rune('0'+argPos))
		args = append(args, *userID)
		argPos++
	}
	if action != nil {
		query += ` AND action = $` + string(rune('0'+argPos))
		args = append(args, *action)
		argPos++
	}
	
	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argPos)) + ` OFFSET $` + string(rune('0'+argPos+1))
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var auditLogs []*models.AuditLog
	for rows.Next() {
		auditLog := &models.AuditLog{}
		var userID, entityID sql.NullString
		var metadataJSON []byte
		
		err := rows.Scan(
			&auditLog.ID,
			&userID,
			&auditLog.Action,
			&auditLog.EntityType,
			&entityID,
			&metadataJSON,
			&auditLog.IPAddress,
			&auditLog.UserAgent,
			&auditLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if userID.Valid {
			if id, err := uuid.Parse(userID.String); err == nil {
				auditLog.UserID = &id
			}
		}
		if entityID.Valid {
			if id, err := uuid.Parse(entityID.String); err == nil {
				auditLog.EntityID = &id
			}
		}
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &auditLog.Metadata)
		}
		
		auditLogs = append(auditLogs, auditLog)
	}
	
	return auditLogs, rows.Err()
}
