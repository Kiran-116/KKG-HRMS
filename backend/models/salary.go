package models

import (
	"time"

	"github.com/google/uuid"
)

type Salary struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	BaseSalary float64   `json:"base_salary"`
	Bonus      float64   `json:"bonus"`
	Deductions float64   `json:"deductions"`
	NetSalary  float64   `json:"net_salary"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateSalaryRequest struct {
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	BaseSalary float64   `json:"base_salary" binding:"required,min=0"`
	Bonus      float64   `json:"bonus" binding:"omitempty,min=0"`
	Deductions float64   `json:"deductions" binding:"omitempty,min=0"`
	Month      int       `json:"month" binding:"required,min=1,max=12"`
	Year       int       `json:"year" binding:"required,min=2000"`
}
