package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"`
	Role        string    `json:"role" db:"role"`
	Department  *string   `json:"department,omitempty" db:"department"`
	Designation *string   `json:"designation,omitempty" db:"designation"`
	JoiningDate *time.Time `json:"joining_date,omitempty" db:"joining_date"`
	Salary      *float64  `json:"salary,omitempty" db:"salary"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	Role        string `json:"role" binding:"omitempty,oneof=admin employee"`
	Department  string `json:"department" binding:"omitempty,max=255"`
	Designation string `json:"designation" binding:"omitempty,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

const (
	RoleAdmin    = "admin"
	RoleEmployee = "employee"
)
