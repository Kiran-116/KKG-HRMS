package websocket

import (
	"hrms/models"

	"github.com/google/uuid"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type NotificationPayload struct {
	Notification *models.Notification `json:"notification"`
}

type LeaveUpdatePayload struct {
	LeaveID uuid.UUID     `json:"leave_id"`
	Status  string        `json:"status"`
	Leave   *models.Leave `json:"leave,omitempty"`
}
