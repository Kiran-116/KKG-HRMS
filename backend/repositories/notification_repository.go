package repositories

import (
	"database/sql"
	"errors"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	GetUnreadCount(userID uuid.UUID) (int, error)
	MarkAsRead(id uuid.UUID, userID uuid.UUID) error
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{
		db: database.DB,
	}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, title, message, type, is_read)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		notification.ID,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.Type,
		notification.IsRead,
	).Scan(&notification.CreatedAt, &notification.UpdatedAt)

	return err
}

func (r *notificationRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	query := `
		SELECT id, user_id, title, message, type, is_read, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]*models.Notification, 0)
	for rows.Next() {
		notification := &models.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Message,
			&notification.Type,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (r *notificationRepository) GetUnreadCount(userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = true, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("notification not found")
	}

	return nil
}
