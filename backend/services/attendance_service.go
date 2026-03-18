package services

import (
	"errors"
	"time"

	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type AttendanceService interface {
	CheckIn(userID uuid.UUID, date time.Time) (*models.Attendance, error)
	CheckOut(userID uuid.UUID, date time.Time) (*models.Attendance, error)
	GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Attendance, error)
	GetAll(page, limit int, userID *uuid.UUID, date *time.Time) ([]*models.Attendance, error)
}

type attendanceService struct {
	attendanceRepo repositories.AttendanceRepository
}

func NewAttendanceService(attendanceRepo repositories.AttendanceRepository) AttendanceService {
	return &attendanceService{
		attendanceRepo: attendanceRepo,
	}
}

func (s *attendanceService) CheckIn(userID uuid.UUID, date time.Time) (*models.Attendance, error) {
	// Get or create attendance record for today
	attendance, err := s.attendanceRepo.GetByUserAndDate(userID, date)
	if err != nil {
		return nil, err
	}

	nowUTC := time.Now().UTC()
	checkInTime := time.Date(date.Year(), date.Month(), date.Day(), nowUTC.Hour(), nowUTC.Minute(), nowUTC.Second(), 0, time.UTC)

	if attendance == nil {
		// Create new attendance record
		attendance = &models.Attendance{
			ID:        uuid.New(),
			UserID:    userID,
			Date:      date,
			CheckIn:   &checkInTime,
			Status:    models.StatusPresent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expectedCheckIn := time.Date(date.Year(), date.Month(), date.Day(), 3, 30, 0, 0, time.UTC)
		if checkInTime.After(expectedCheckIn) {
			attendance.Status = models.StatusLate
		}

		if err := s.attendanceRepo.Create(attendance); err != nil {
			return nil, errors.New("failed to check in")
		}
	} else {
		// Update existing record
		if attendance.CheckIn != nil {
			return nil, errors.New("already checked in for this date")
		}
		attendance.CheckIn = &checkInTime
		attendance.Status = models.StatusPresent

		expectedCheckIn := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
		if checkInTime.After(expectedCheckIn) {
			attendance.Status = models.StatusLate
		}

		if err := s.attendanceRepo.Update(attendance); err != nil {
			return nil, errors.New("failed to update check-in")
		}
	}

	return attendance, nil
}

func (s *attendanceService) CheckOut(userID uuid.UUID, date time.Time) (*models.Attendance, error) {
	attendance, err := s.attendanceRepo.GetByUserAndDate(userID, date)
	if err != nil {
		return nil, err
	}

	if attendance == nil {
		return nil, errors.New("no check-in found for this date")
	}

	if attendance.CheckIn == nil {
		return nil, errors.New("must check in before checking out")
	}

	if attendance.CheckOut != nil {
		return nil, errors.New("already checked out for this date")
	}

	nowUTC := time.Now().UTC()
	checkOutTime := time.Date(date.Year(), date.Month(), date.Day(), nowUTC.Hour(), nowUTC.Minute(), nowUTC.Second(), 0, time.UTC)
	attendance.CheckOut = &checkOutTime

	// Check if half day (less than 4 hours)
	duration := checkOutTime.Sub(*attendance.CheckIn)
	if duration < 4*time.Hour {
		attendance.Status = models.StatusHalfDay
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, errors.New("failed to check out")
	}

	return attendance, nil
}

func (s *attendanceService) GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Attendance, error) {
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
	return s.attendanceRepo.GetByUserID(userID, limit, offset)
}

func (s *attendanceService) GetAll(page, limit int, userID *uuid.UUID, date *time.Time) ([]*models.Attendance, error) {
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
	return s.attendanceRepo.GetAll(limit, offset, userID, date)
}
