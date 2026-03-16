package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	UserName  *string    `json:"user_name,omitempty"` // Optional, populated when joining with users table
	Date      time.Time  `json:"date"`
	CheckIn   *time.Time `json:"check_in,omitempty"`
	CheckOut  *time.Time `json:"check_out,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CheckInRequest struct {
	Date string `json:"date" binding:"omitempty"`
}

type CheckOutRequest struct {
	Date string `json:"date" binding:"omitempty"`
}

const (
	StatusPresent = "present"
	StatusAbsent  = "absent"
	StatusLate    = "late"
	StatusHalfDay = "half_day"
)
