package services

import (
	"context"
	"errors"
	"time"

	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type SalaryService interface {
	CreateSalary(ctx context.Context, req *models.CreateSalaryRequest) (*models.Salary, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Salary, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Salary, error)
}

type salaryService struct {
	salaryRepo repositories.SalaryRepository
}

func NewSalaryService(salaryRepo repositories.SalaryRepository) SalaryService {
	return &salaryService{
		salaryRepo: salaryRepo,
	}
}

func (s *salaryService) CreateSalary(ctx context.Context, req *models.CreateSalaryRequest) (*models.Salary, error) {
	// Check if salary already exists for this month/year
	existing, err := s.salaryRepo.GetByUserIDAndMonthYear(ctx, req.UserID, req.Month, req.Year)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("salary already exists for this month and year")
	}

	salary := &models.Salary{
		ID:         uuid.New(),
		UserID:     req.UserID,
		BaseSalary: req.BaseSalary,
		Bonus:      req.Bonus,
		Deductions: req.Deductions,
		Month:      req.Month,
		Year:       req.Year,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.salaryRepo.Create(ctx, salary); err != nil {
		return nil, errors.New("failed to create salary record")
	}

	return salary, nil
}

func (s *salaryService) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Salary, error) {
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
	return s.salaryRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *salaryService) GetByID(ctx context.Context, id uuid.UUID) (*models.Salary, error) {
	return s.salaryRepo.GetByID(ctx, id)
}
