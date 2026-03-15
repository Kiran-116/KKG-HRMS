package repositories

import (
	"database/sql"
	"errors"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*models.User, error)
	Count() (int, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role, department, designation, joining_date, salary, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Department,
		user.Designation,
		user.JoiningDate,
		user.Salary,
		user.IsActive,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password_hash, role, department, designation, 
		       joining_date, salary, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Department,
		&user.Designation,
		&user.JoiningDate,
		&user.Salary,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password_hash, role, department, designation, 
		       joining_date, salary, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Department,
		&user.Designation,
		&user.JoiningDate,
		&user.Salary,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, role = $4, department = $5, designation = $6,
		    joining_date = $7, salary = $8, is_active = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	return r.db.QueryRow(
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.Department,
		user.Designation,
		user.JoiningDate,
		user.Salary,
		user.IsActive,
	).Scan(&user.UpdatedAt)
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
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

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var joiningDate sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Department,
			&user.Designation,
			&joiningDate,
			&user.Salary,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if joiningDate.Valid {
			user.JoiningDate = &joiningDate.Time
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *userRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}
