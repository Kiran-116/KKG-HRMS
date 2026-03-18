package repositories

import (
	"database/sql"
	"errors"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type LeaveRepository interface {
	Create(leave *models.Leave) error
	GetByID(id uuid.UUID) (*models.Leave, error)
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Leave, error)
	GetAll(limit, offset int, status *string) ([]*models.Leave, error)
	Update(leave *models.Leave) error
}

type leaveRepository struct {
	db *sql.DB
}

func NewLeaveRepository() LeaveRepository {
	return &leaveRepository{
		db: database.DB,
	}
}

func (r *leaveRepository) Create(leave *models.Leave) error {
	query := `
		INSERT INTO leaves (id, user_id, start_date, end_date, reason, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	
	err := r.db.QueryRow(
		query,
		leave.ID,
		leave.UserID,
		leave.StartDate,
		leave.EndDate,
		leave.Reason,
		leave.Status,
	).Scan(&leave.CreatedAt, &leave.UpdatedAt)
	
	return err
}

func (r *leaveRepository) GetByID(id uuid.UUID) (*models.Leave, error) {
	leave := &models.Leave{}
	query := `
		SELECT id, user_id, start_date, end_date, reason, status, approved_by, created_at, updated_at
		FROM leaves
		WHERE id = $1
	`
	
	var approvedBy sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&leave.ID,
		&leave.UserID,
		&leave.StartDate,
		&leave.EndDate,
		&leave.Reason,
		&leave.Status,
		&approvedBy,
		&leave.CreatedAt,
		&leave.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("leave not found")
	}
	if err != nil {
		return nil, err
	}
	
	if approvedBy.Valid {
		if id, err := uuid.Parse(approvedBy.String); err == nil {
			leave.ApprovedBy = &id
		}
	}
	
	return leave, nil
}

func (r *leaveRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Leave, error) {
	query := `
		SELECT id, user_id, start_date, end_date, reason, status, approved_by, created_at, updated_at
		FROM leaves
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var leaves []*models.Leave
	for rows.Next() {
		leave := &models.Leave{}
		var approvedBy sql.NullString
		
		err := rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.StartDate,
			&leave.EndDate,
			&leave.Reason,
			&leave.Status,
			&approvedBy,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if approvedBy.Valid {
			if id, err := uuid.Parse(approvedBy.String); err == nil {
				leave.ApprovedBy = &id
			}
		}
		
		leaves = append(leaves, leave)
	}
	
	return leaves, rows.Err()
}

func (r *leaveRepository) GetAll(limit, offset int, status *string) ([]*models.Leave, error) {
	query := `
		SELECT l.id, l.user_id, u.name as user_name, l.start_date, l.end_date, l.reason, l.status, l.approved_by, l.created_at, l.updated_at
		FROM leaves l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1
	
	if status != nil {
		query += ` AND l.status = $` + string(rune('0'+argPos))
		args = append(args, *status)
		argPos++
	}
	
	query += ` ORDER BY l.created_at DESC LIMIT $` + string(rune('0'+argPos)) + ` OFFSET $` + string(rune('0'+argPos+1))
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var leaves []*models.Leave
	for rows.Next() {
		leave := &models.Leave{}
		var approvedBy sql.NullString
		var userName sql.NullString
		
		err := rows.Scan(
			&leave.ID,
			&leave.UserID,
			&userName,
			&leave.StartDate,
			&leave.EndDate,
			&leave.Reason,
			&leave.Status,
			&approvedBy,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if userName.Valid {
			leave.UserName = &userName.String
		}
		if approvedBy.Valid {
			if id, err := uuid.Parse(approvedBy.String); err == nil {
				leave.ApprovedBy = &id
			}
		}
		
		leaves = append(leaves, leave)
	}
	
	return leaves, rows.Err()
}

func (r *leaveRepository) Update(leave *models.Leave) error {
	query := `
		UPDATE leaves
		SET status = $2, approved_by = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	return r.db.QueryRow(
		query,
		leave.ID,
		leave.Status,
		leave.ApprovedBy,
	).Scan(&leave.UpdatedAt)
}
