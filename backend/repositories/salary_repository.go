package repositories

import (
	"database/sql"
	"errors"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type SalaryRepository interface {
	Create(salary *models.Salary) error
	GetByID(id uuid.UUID) (*models.Salary, error)
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Salary, error)
	GetByUserIDAndMonthYear(userID uuid.UUID, month, year int) (*models.Salary, error)
}

type salaryRepository struct {
	db *sql.DB
}

func NewSalaryRepository() SalaryRepository {
	return &salaryRepository{
		db: database.DB,
	}
}

func (r *salaryRepository) Create(salary *models.Salary) error {
	query := `
		INSERT INTO salaries (id, user_id, base_salary, bonus, deductions, month, year)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING net_salary, created_at, updated_at
	`
	
	err := r.db.QueryRow(
		query,
		salary.ID,
		salary.UserID,
		salary.BaseSalary,
		salary.Bonus,
		salary.Deductions,
		salary.Month,
		salary.Year,
	).Scan(&salary.NetSalary, &salary.CreatedAt, &salary.UpdatedAt)
	
	return err
}

func (r *salaryRepository) GetByID(id uuid.UUID) (*models.Salary, error) {
	salary := &models.Salary{}
	query := `
		SELECT id, user_id, base_salary, bonus, deductions, net_salary, month, year, created_at, updated_at
		FROM salaries
		WHERE id = $1
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&salary.ID,
		&salary.UserID,
		&salary.BaseSalary,
		&salary.Bonus,
		&salary.Deductions,
		&salary.NetSalary,
		&salary.Month,
		&salary.Year,
		&salary.CreatedAt,
		&salary.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("salary not found")
	}
	
	return salary, err
}

func (r *salaryRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Salary, error) {
	query := `
		SELECT id, user_id, base_salary, bonus, deductions, net_salary, month, year, created_at, updated_at
		FROM salaries
		WHERE user_id = $1
		ORDER BY year DESC, month DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var salaries []*models.Salary
	for rows.Next() {
		salary := &models.Salary{}
		err := rows.Scan(
			&salary.ID,
			&salary.UserID,
			&salary.BaseSalary,
			&salary.Bonus,
			&salary.Deductions,
			&salary.NetSalary,
			&salary.Month,
			&salary.Year,
			&salary.CreatedAt,
			&salary.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		salaries = append(salaries, salary)
	}
	
	return salaries, rows.Err()
}

func (r *salaryRepository) GetByUserIDAndMonthYear(userID uuid.UUID, month, year int) (*models.Salary, error) {
	salary := &models.Salary{}
	query := `
		SELECT id, user_id, base_salary, bonus, deductions, net_salary, month, year, created_at, updated_at
		FROM salaries
		WHERE user_id = $1 AND month = $2 AND year = $3
	`
	
	err := r.db.QueryRow(query, userID, month, year).Scan(
		&salary.ID,
		&salary.UserID,
		&salary.BaseSalary,
		&salary.Bonus,
		&salary.Deductions,
		&salary.NetSalary,
		&salary.Month,
		&salary.Year,
		&salary.CreatedAt,
		&salary.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return salary, err
}
