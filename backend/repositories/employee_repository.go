package repositories

import (
	"database/sql"
	"errors"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type EmployeeRepository interface {
	List(limit, offset int) ([]*models.Employee, error)
	GetByID(id uuid.UUID) (*models.Employee, error)
	Count() (int, error)
}

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository() EmployeeRepository {
	return &employeeRepository{
		db: database.DB,
	}
}

func (r *employeeRepository) List(limit, offset int) ([]*models.Employee, error) {
	query := `
		SELECT id, name, email, role, department, designation, 
		       joining_date, salary, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		emp := &models.Employee{}
		var joiningDate sql.NullTime

		err := rows.Scan(
			&emp.ID,
			&emp.Name,
			&emp.Email,
			&emp.Role,
			&emp.Department,
			&emp.Designation,
			&joiningDate,
			&emp.Salary,
			&emp.IsActive,
			&emp.CreatedAt,
			&emp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if joiningDate.Valid {
			emp.JoiningDate = &joiningDate.Time
		}

		employees = append(employees, emp)
	}

	return employees, rows.Err()
}

func (r *employeeRepository) GetByID(id uuid.UUID) (*models.Employee, error) {
	emp := &models.Employee{}
	query := `
		SELECT id, name, email, role, department, designation, 
		       joining_date, salary, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var joiningDate sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&emp.ID,
		&emp.Name,
		&emp.Email,
		&emp.Role,
		&emp.Department,
		&emp.Designation,
		&joiningDate,
		&emp.Salary,
		&emp.IsActive,
		&emp.CreatedAt,
		&emp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("employee not found")
	}
	if err != nil {
		return nil, err
	}

	if joiningDate.Valid {
		emp.JoiningDate = &joiningDate.Time
	}

	return emp, nil
}

func (r *employeeRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}
