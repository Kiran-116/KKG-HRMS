package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID              `json:"id"`
	UserID      *uuid.UUID             `json:"user_id,omitempty"`
	UserName    *string                `json:"user_name,omitempty"` // Optional, populated when joining with users table
	Action      string                 `json:"action"`
	EntityType  string                 `json:"entity_type"`
	EntityID    *uuid.UUID             `json:"entity_id,omitempty"`
	Description *string                `json:"description,omitempty"` // Human-readable change description
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
