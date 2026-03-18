package models

import (
	"time"

	"github.com/google/uuid"
)

type Leave struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	UserName   *string    `json:"user_name,omitempty"` // Optional, populated when joining with users table
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	Reason     string     `json:"reason"`
	Status     string     `json:"status"`
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ApplyLeaveRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason" binding:"required,min=10"`
}

const (
	LeaveStatusPending  = "pending"
	LeaveStatusApproved = "approved"
	LeaveStatusRejected = "rejected"
)
