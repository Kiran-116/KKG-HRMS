package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByMagicToken(ctx context.Context, token string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	ListAdmins(ctx context.Context) ([]*models.User, error)
	Count(ctx context.Context) (int, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role, department, designation, joining_date, salary, is_active, magic_token, magic_expires_at, must_change_password)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
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
		user.MagicToken,
		user.MagicExpiresAt,
		user.MustChangePassword,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password_hash, role, department, designation, 
		       joining_date, salary, is_active, magic_token, magic_expires_at, must_change_password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var magicToken sql.NullString
	var magicExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		&magicToken,
		&magicExpiresAt,
		&user.MustChangePassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if magicToken.Valid {
		user.MagicToken = &magicToken.String
	}
	if magicExpiresAt.Valid {
		user.MagicExpiresAt = &magicExpiresAt.Time
	}

	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password_hash, role, department, designation, 
		       joining_date, salary, is_active, magic_token, magic_expires_at, must_change_password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var magicToken sql.NullString
	var magicExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
		&magicToken,
		&magicExpiresAt,
		&user.MustChangePassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if magicToken.Valid {
		user.MagicToken = &magicToken.String
	}
	if magicExpiresAt.Valid {
		user.MagicExpiresAt = &magicExpiresAt.Time
	}

	return user, err
}

func (r *userRepository) GetByMagicToken(ctx context.Context, token string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password_hash, role, department, designation, 
		       joining_date, salary, is_active, magic_token, magic_expires_at, must_change_password, created_at, updated_at
		FROM users
		WHERE magic_token = $1 AND magic_expires_at > CURRENT_TIMESTAMP
	`

	var magicToken sql.NullString
	var magicExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, token).Scan(
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
		&magicToken,
		&magicExpiresAt,
		&user.MustChangePassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid or expired magic token")
	}

	if magicToken.Valid {
		user.MagicToken = &magicToken.String
	}
	if magicExpiresAt.Valid {
		user.MagicExpiresAt = &magicExpiresAt.Time
	}

	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, role = $4, department = $5, designation = $6,
		    joining_date = $7, salary = $8, is_active = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx,
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

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, must_change_password = false, magic_token = NULL, magic_expires_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, userID, passwordHash).Scan(&updatedAt)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, name, email, role, department, designation, 
		       joining_date, salary, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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

func (r *userRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *userRepository) ListAdmins(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, name, email, role, department, designation,
		       joining_date, salary, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true AND role = 'admin'
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
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
