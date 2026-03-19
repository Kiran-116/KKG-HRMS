package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Department  *string    `json:"department,omitempty"`
	Designation *string    `json:"designation,omitempty"`
	JoiningDate *time.Time `json:"joining_date,omitempty"`
	Salary      *float64   `json:"salary,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateEmployeeRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Email       string  `json:"email" binding:"required,email"`
	Role        string  `json:"role" binding:"omitempty,oneof=admin employee"`
	Department  string  `json:"department" binding:"omitempty,max=255"`
	Designation string  `json:"designation" binding:"omitempty,max=255"`
	JoiningDate string  `json:"joining_date" binding:"omitempty"`
	Salary      float64 `json:"salary" binding:"omitempty,min=0"`
}

type UpdateEmployeeRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=2,max=255"`
	Email       string  `json:"email" binding:"omitempty,email"`
	Role        string  `json:"role" binding:"omitempty,oneof=admin employee"`
	Department  string  `json:"department" binding:"omitempty,max=255"`
	Designation string  `json:"designation" binding:"omitempty,max=255"`
	JoiningDate string  `json:"joining_date" binding:"omitempty"`
	Salary      float64 `json:"salary" binding:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
}

type EmployeeListResponse struct {
	Employees []*Employee `json:"employees"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
}
