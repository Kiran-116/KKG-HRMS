package services

import (
	"errors"
	"time"

	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type LeaveService interface {
	ApplyLeave(userID uuid.UUID, req *models.ApplyLeaveRequest) (*models.Leave, error)
	GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Leave, error)
	GetAll(page, limit int, status *string) ([]*models.Leave, error)
	ApproveLeave(leaveID uuid.UUID, approvedBy uuid.UUID) (*models.Leave, error)
	RejectLeave(leaveID uuid.UUID, approvedBy uuid.UUID) (*models.Leave, error)
	GetByID(id uuid.UUID) (*models.Leave, error)
}

type leaveService struct {
	leaveRepo repositories.LeaveRepository
}

func NewLeaveService(leaveRepo repositories.LeaveRepository) LeaveService {
	return &leaveService{
		leaveRepo: leaveRepo,
	}
}

func (s *leaveService) ApplyLeave(userID uuid.UUID, req *models.ApplyLeaveRequest) (*models.Leave, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("start date cannot be in the past")
	}

	leave := &models.Leave{
		ID:        uuid.New(),
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    req.Reason,
		Status:    models.LeaveStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.leaveRepo.Create(leave); err != nil {
		return nil, errors.New("failed to apply leave")
	}

	return leave, nil
}

func (s *leaveService) GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Leave, error) {
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
	return s.leaveRepo.GetByUserID(userID, limit, offset)
}

func (s *leaveService) GetAll(page, limit int, status *string) ([]*models.Leave, error) {
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
	return s.leaveRepo.GetAll(limit, offset, status)
}

func (s *leaveService) ApproveLeave(leaveID uuid.UUID, approvedBy uuid.UUID) (*models.Leave, error) {
	leave, err := s.leaveRepo.GetByID(leaveID)
	if err != nil {
		return nil, errors.New("leave not found")
	}

	if leave.Status != models.LeaveStatusPending {
		return nil, errors.New("leave is not pending")
	}

	leave.Status = models.LeaveStatusApproved
	leave.ApprovedBy = &approvedBy

	if err := s.leaveRepo.Update(leave); err != nil {
		return nil, errors.New("failed to approve leave")
	}

	return leave, nil
}

func (s *leaveService) RejectLeave(leaveID uuid.UUID, approvedBy uuid.UUID) (*models.Leave, error) {
	leave, err := s.leaveRepo.GetByID(leaveID)
	if err != nil {
		return nil, errors.New("leave not found")
	}

	if leave.Status != models.LeaveStatusPending {
		return nil, errors.New("leave is not pending")
	}

	leave.Status = models.LeaveStatusRejected
	leave.ApprovedBy = &approvedBy

	if err := s.leaveRepo.Update(leave); err != nil {
		return nil, errors.New("failed to reject leave")
	}

	return leave, nil
}

func (s *leaveService) GetByID(id uuid.UUID) (*models.Leave, error) {
	return s.leaveRepo.GetByID(id)
}
