package services

import (
	"errors"
	"hrms/websocket"
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
	leaveRepo           repositories.LeaveRepository
	notificationService NotificationService
	hub                 *websocket.Hub
	userRepo            repositories.UserRepository
}

func NewLeaveService(
	leaveRepo repositories.LeaveRepository,
	notificationService NotificationService,
	hub *websocket.Hub,
	userRepo repositories.UserRepository,
) LeaveService {
	return &leaveService{
		leaveRepo:           leaveRepo,
		notificationService: notificationService,
		hub:                 hub,
		userRepo:            userRepo,
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

	// Notify all admins about new leave request
	if s.notificationService != nil && s.userRepo != nil {
		admins, err := s.userRepo.ListAdmins()
		if err == nil {
			employeeName := "An employee"
			if u, err := s.userRepo.GetByID(userID); err == nil && u != nil && u.Name != "" {
				employeeName = u.Name
			}

			title := "New Leave Request"
			msg := employeeName + " applied for leave (" + req.StartDate + " to " + req.EndDate + ")."

			for _, admin := range admins {
				if admin == nil {
					continue
				}
				_, _ = s.notificationService.CreateNotification(admin.ID, title, msg, models.NotificationTypeInfo)
			}
		}
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

	// Notify employee
	if s.notificationService != nil {
		_, _ = s.notificationService.CreateNotification(
			leave.UserID,
			"Leave Approved",
			"Your leave request has been approved.",
			"success",
		)
	}

	// Broadcast leave update in real-time
	if s.hub != nil {
		s.hub.BroadcastToUser(leave.UserID, websocket.Message{
			Type: "leave_update",
			Payload: websocket.LeaveUpdatePayload{
				LeaveID: leave.ID,
				Status:  leave.Status,
				Leave:   leave,
			},
		})
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

	// Notify employee
	if s.notificationService != nil {
		_, _ = s.notificationService.CreateNotification(
			leave.UserID,
			"Leave Rejected",
			"Your leave request has been rejected.",
			"warning",
		)
	}

	// Broadcast leave update in real-time
	if s.hub != nil {
		s.hub.BroadcastToUser(leave.UserID, websocket.Message{
			Type: "leave_update",
			Payload: websocket.LeaveUpdatePayload{
				LeaveID: leave.ID,
				Status:  leave.Status,
				Leave:   leave,
			},
		})
	}

	return leave, nil
}

func (s *leaveService) GetByID(id uuid.UUID) (*models.Leave, error) {
	return s.leaveRepo.GetByID(id)
}
