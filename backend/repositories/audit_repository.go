package repositories

import (
	"context"
	"database/sql"
	"encoding/json"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type AuditRepository interface {
	Create(ctx context.Context, auditLog *models.AuditLog) error
	GetAll(ctx context.Context, limit, offset int, userID *uuid.UUID, action *string) ([]*models.AuditLog, error)
}

type auditRepository struct {
	db *sql.DB
}

func NewAuditRepository() AuditRepository {
	return &auditRepository{
		db: database.DB,
	}
}

func (r *auditRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	metadataJSON, _ := json.Marshal(auditLog.Metadata)

	query := `
		INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, description, metadata, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		auditLog.ID,
		auditLog.UserID,
		auditLog.Action,
		auditLog.EntityType,
		auditLog.EntityID,
		auditLog.Description,
		metadataJSON,
		auditLog.IPAddress,
		auditLog.UserAgent,
	)

	return err
}

func (r *auditRepository) GetAll(ctx context.Context, limit, offset int, userID *uuid.UUID, action *string) ([]*models.AuditLog, error) {
	query := `
		SELECT a.id, a.user_id, u.name as user_name, a.action, a.entity_type, a.entity_id, a.description, a.metadata, a.ip_address, a.user_agent, a.created_at
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	if userID != nil {
		query += ` AND a.user_id = $` + string(rune('0'+argPos))
		args = append(args, *userID)
		argPos++
	}
	if action != nil {
		query += ` AND a.action = $` + string(rune('0'+argPos))
		args = append(args, *action)
		argPos++
	}

	query += ` ORDER BY a.created_at DESC LIMIT $` + string(rune('0'+argPos)) + ` OFFSET $` + string(rune('0'+argPos+1))
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auditLogs []*models.AuditLog
	for rows.Next() {
		auditLog := &models.AuditLog{}
		var userID, entityID sql.NullString
		var userName sql.NullString
		var description sql.NullString
		var metadataJSON []byte

		err := rows.Scan(
			&auditLog.ID,
			&userID,
			&userName,
			&auditLog.Action,
			&auditLog.EntityType,
			&entityID,
			&description,
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
		if userName.Valid {
			auditLog.UserName = &userName.String
		}
		if description.Valid {
			auditLog.Description = &description.String
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
