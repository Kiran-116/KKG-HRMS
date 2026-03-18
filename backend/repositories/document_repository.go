package repositories

import (
	"context"
	"database/sql"
	"errors"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(ctx context.Context, document *models.Document) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Document, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type documentRepository struct {
	db *sql.DB
}

func NewDocumentRepository() DocumentRepository {
	return &documentRepository{
		db: database.DB,
	}
}

func (r *documentRepository) Create(ctx context.Context, document *models.Document) error {
	query := `
		INSERT INTO documents (id, user_id, file_url, file_name, file_size, document_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uploaded_at, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		document.ID,
		document.UserID,
		document.FileURL,
		document.FileName,
		document.FileSize,
		document.DocumentType,
	).Scan(&document.UploadedAt, &document.CreatedAt, &document.UpdatedAt)

	return err
}

func (r *documentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	document := &models.Document{}
	query := `
		SELECT id, user_id, file_url, file_name, file_size, document_type, uploaded_at, created_at, updated_at
		FROM documents
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&document.ID,
		&document.UserID,
		&document.FileURL,
		&document.FileName,
		&document.FileSize,
		&document.DocumentType,
		&document.UploadedAt,
		&document.CreatedAt,
		&document.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("document not found")
	}

	return document, err
}

func (r *documentRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Document, error) {
	query := `
		SELECT id, user_id, file_url, file_name, file_size, document_type, uploaded_at, created_at, updated_at
		FROM documents
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		document := &models.Document{}
		err := rows.Scan(
			&document.ID,
			&document.UserID,
			&document.FileURL,
			&document.FileName,
			&document.FileSize,
			&document.DocumentType,
			&document.UploadedAt,
			&document.CreatedAt,
			&document.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, rows.Err()
}

func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("document not found")
	}

	return nil
}
