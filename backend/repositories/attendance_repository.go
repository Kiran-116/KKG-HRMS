package repositories

import (
	"database/sql"
	"errors"
	"time"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type AttendanceRepository interface {
	Create(attendance *models.Attendance) error
	GetByID(id uuid.UUID) (*models.Attendance, error)
	GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Attendance, error)
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Attendance, error)
	GetAll(limit, offset int, userID *uuid.UUID, date *time.Time) ([]*models.Attendance, error)
	Update(attendance *models.Attendance) error
}

type attendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository() AttendanceRepository {
	return &attendanceRepository{
		db: database.DB,
	}
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	query := `
		INSERT INTO attendance (id, user_id, date, check_in, check_out, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	
	err := r.db.QueryRow(
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

func (r *attendanceRepository) GetByID(id uuid.UUID) (*models.Attendance, error) {
	attendance := &models.Attendance{}
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE id = $1
	`
	
	var checkIn, checkOut sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
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

func (r *attendanceRepository) GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Attendance, error) {
	attendance := &models.Attendance{}
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE user_id = $1 AND date = $2
	`
	
	var checkIn, checkOut sql.NullTime
	err := r.db.QueryRow(query, userID, date).Scan(
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

func (r *attendanceRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Attendance, error) {
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
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

func (r *attendanceRepository) GetAll(limit, offset int, userID *uuid.UUID, date *time.Time) ([]*models.Attendance, error) {
	query := `
		SELECT id, user_id, date, check_in, check_out, status, created_at, updated_at
		FROM attendance
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1
	
	if userID != nil {
		query += ` AND user_id = $` + string(rune('0'+argPos))
		args = append(args, *userID)
		argPos++
	}
	if date != nil {
		query += ` AND date = $` + string(rune('0'+argPos))
		args = append(args, *date)
		argPos++
	}
	
	query += ` ORDER BY date DESC LIMIT $` + string(rune('0'+argPos)) + ` OFFSET $` + string(rune('0'+argPos+1))
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
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

func (r *attendanceRepository) Update(attendance *models.Attendance) error {
	query := `
		UPDATE attendance
		SET check_in = $2, check_out = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	return r.db.QueryRow(
		query,
		attendance.ID,
		attendance.CheckIn,
		attendance.CheckOut,
		attendance.Status,
	).Scan(&attendance.UpdatedAt)
}
