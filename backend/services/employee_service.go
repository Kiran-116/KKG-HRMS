package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"hrms/config"
	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req *models.CreateEmployeeRequest) (*models.Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, req *models.UpdateEmployeeRequest) (*models.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	ListEmployees(ctx context.Context, page, limit int) (*models.EmployeeListResponse, error)
	DeactivateEmployee(ctx context.Context, id uuid.UUID) error
}

type employeeService struct {
	userRepo     repositories.UserRepository
	employeeRepo repositories.EmployeeRepository
	emailService EmailService
}

func NewEmployeeService(userRepo repositories.UserRepository, employeeRepo repositories.EmployeeRepository, emailService EmailService) EmployeeService {
	return &employeeService{
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
		emailService: emailService,
	}
}

func (s *employeeService) CreateEmployee(ctx context.Context, req *models.CreateEmployeeRequest) (*models.Employee, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("employee with this email already exists")
	}

	// Generate magic token
	magicToken, err := generateMagicToken()
	if err != nil {
		return nil, errors.New("failed to generate magic token")
	}

	// Set expiry to 24 hours from now
	expiresAt := time.Now().Add(24 * time.Hour)

	// Set default role
	role := req.Role
	if role == "" {
		role = models.RoleEmployee
	}

	// Parse joining date if provided
	var joiningDate *time.Time
	if req.JoiningDate != "" {
		parsed, err := time.Parse("2006-01-02", req.JoiningDate)
		if err == nil {
			joiningDate = &parsed
		}
	}

	// Create user with magic token (no password hash initially)
	user := &models.User{
		ID:                uuid.New(),
		Name:              req.Name,
		Email:             req.Email,
		PasswordHash:      "", // Empty password hash - will be set when user sets password
		Role:              role,
		IsActive:          true,
		MagicToken:        &magicToken,
		MagicExpiresAt:    &expiresAt,
		MustChangePassword: true,
	}

	if req.Department != "" {
		user.Department = &req.Department
	}
	if req.Designation != "" {
		user.Designation = &req.Designation
	}
	if joiningDate != nil {
		user.JoiningDate = joiningDate
	}
	if req.Salary > 0 {
		user.Salary = &req.Salary
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create employee")
	}

	// Generate magic link
	frontendURL := config.AppConfig.Server.FrontendURL
	magicLink := fmt.Sprintf("%s/magic-login?token=%s", frontendURL, magicToken)

	// Send magic link email (async)
	go func() {
		if err := s.emailService.SendMagicLinkEmail(user.Email, user.Name, magicLink); err != nil {
			// Log error but don't fail employee creation
			fmt.Printf("Failed to send magic link email to %s: %v\n", user.Email, err)
		}
	}()

	return s.employeeRepo.GetByID(ctx, user.ID)
}

// generateMagicToken generates a secure random token for magic link
func generateMagicToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, id uuid.UUID, req *models.UpdateEmployeeRequest) (*models.Employee, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		// Check if new email already exists
		existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Department != "" {
		user.Department = &req.Department
	} else if req.Department == "" && user.Department != nil {
		// Allow clearing department
		user.Department = nil
	}
	if req.Designation != "" {
		user.Designation = &req.Designation
	} else if req.Designation == "" && user.Designation != nil {
		user.Designation = nil
	}
	if req.JoiningDate != "" {
		parsed, err := time.Parse("2006-01-02", req.JoiningDate)
		if err == nil {
			user.JoiningDate = &parsed
		}
	}
	if req.Salary > 0 {
		user.Salary = &req.Salary
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.New("failed to update employee")
	}

	return s.employeeRepo.GetByID(ctx, id)
}

func (s *employeeService) GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	return s.employeeRepo.GetByID(ctx, id)
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	return s.employeeRepo.GetByID(ctx, id)
}

func (s *employeeService) ListEmployees(ctx context.Context, page, limit int) (*models.EmployeeListResponse, error) {
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

	employees, err := s.employeeRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.employeeRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &models.EmployeeListResponse{
		Employees: employees,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

func (s *employeeService) DeactivateEmployee(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
