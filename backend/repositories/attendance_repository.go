package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *models.Attendance) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Attendance, error)
	GetByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.Attendance, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Attendance, error)
	GetAll(ctx context.Context, limit, offset int, userID *uuid.UUID, date *time.Time) ([]*models.Attendance, error)
	Update(ctx context.Context, attendance *models.Attendance) error
}

type attendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository() AttendanceRepository {
	return &attendanceRepository{
		db: database.DB,
	}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *models.Attendance) error {
	query := `
		INSERT INTO attendance (id, user_id, date, check_in, check_out, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		attendance.ID,
		attendance.UserID,
		attendance.Date,
		attendance.CheckIn,
		attendance.CheckOut,
		attendance.Status,
	).Scan(&attendance.CreatedAt, &attendance.UpdatedAt)

	return err
}

func (r *attendanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Attendance, error) {
	attendance := &models.Attendance{}
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE id = $1
	`

	var checkIn, checkOut sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Date,
		&checkIn,
		&checkOut,
		&attendance.Status,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("attendance not found")
	}
	if err != nil {
		return nil, err
	}

	if checkIn.Valid {
		attendance.CheckIn = &checkIn.Time
	}
	if checkOut.Valid {
		attendance.CheckOut = &checkOut.Time
	}

	return attendance, nil
}

func (r *attendanceRepository) GetByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*models.Attendance, error) {
	attendance := &models.Attendance{}
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE user_id = $1 AND date = $2
	`

	var checkIn, checkOut sql.NullTime
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Date,
		&checkIn,
		&checkOut,
		&attendance.Status,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found, but not an error
	}
	if err != nil {
		return nil, err
	}

	if checkIn.Valid {
		attendance.CheckIn = &checkIn.Time
	}
	if checkOut.Valid {
		attendance.CheckOut = &checkOut.Time
	}

	return attendance, nil
}

func (r *attendanceRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Attendance, error) {
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*models.Attendance
	for rows.Next() {
		attendance := &models.Attendance{}
		var checkIn, checkOut sql.NullTime

		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.Date,
			&checkIn,
			&checkOut,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if checkIn.Valid {
			attendance.CheckIn = &checkIn.Time
		}
		if checkOut.Valid {
			attendance.CheckOut = &checkOut.Time
		}

		attendances = append(attendances, attendance)
	}

	return attendances, rows.Err()
}

func (r *attendanceRepository) GetAll(ctx context.Context, limit, offset int, userID *uuid.UUID, date *time.Time) ([]*models.Attendance, error) {
	query := `
		SELECT a.id, a.user_id, u.name as user_name, a.date, a.check_in, a.check_out, a.status, a.created_at, a.updated_at
		FROM attendance a
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
	if date != nil {
		query += ` AND a.date = $` + string(rune('0'+argPos))
		args = append(args, *date)
		argPos++
	}

	query += ` ORDER BY a.date DESC LIMIT $` + string(rune('0'+argPos)) + ` OFFSET $` + string(rune('0'+argPos+1))
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*models.Attendance
	for rows.Next() {
		attendance := &models.Attendance{}
		var checkIn, checkOut sql.NullTime
		var userName sql.NullString

		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&userName,
			&attendance.Date,
			&checkIn,
			&checkOut,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if userName.Valid {
			attendance.UserName = &userName.String
		}
		if checkIn.Valid {
			attendance.CheckIn = &checkIn.Time
		}
		if checkOut.Valid {
			attendance.CheckOut = &checkOut.Time
		}

		attendances = append(attendances, attendance)
	}

	return attendances, rows.Err()
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *models.Attendance) error {
	query := `
		UPDATE attendance
		SET check_in = $2, check_out = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		attendance.ID,
		attendance.CheckIn,
		attendance.CheckOut,
		attendance.Status,
	).Scan(&attendance.UpdatedAt)
}
