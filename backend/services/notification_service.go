package services

import (
	"hrms/websocket"
	"time"

	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(userID uuid.UUID, title, message, notificationType string) (*models.Notification, error)
	GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Notification, error)
	GetUnreadCount(userID uuid.UUID) (int, error)
	MarkAsRead(id uuid.UUID, userID uuid.UUID) error
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
	emailService     EmailService
	hub              *websocket.Hub
}

func NewNotificationService(notificationRepo repositories.NotificationRepository, emailService EmailService, hub *websocket.Hub) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		emailService:     emailService,
		hub:              hub,
	}
}

func (s *notificationService) CreateNotification(userID uuid.UUID, title, message, notificationType string) (*models.Notification, error) {
	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      notificationType,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return nil, err
	}

	// Send email notification (async in production)
	go func() {
		// Get user email from repository if needed
		s.emailService.SendEmail("user@example.com", title, message)
	}()

	// Broadcast real-time notification
	if s.hub != nil {
		s.hub.BroadcastToUser(userID, websocket.Message{
			Type: "notification",
			Payload: websocket.NotificationPayload{
				Notification: notification,
			},
		})
	}

	return notification, nil
}

func (s *notificationService) GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Notification, error) {
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
	return s.notificationRepo.GetByUserID(userID, limit, offset)
}

func (s *notificationService) GetUnreadCount(userID uuid.UUID) (int, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

func (s *notificationService) MarkAsRead(id uuid.UUID, userID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(id, userID)
}
