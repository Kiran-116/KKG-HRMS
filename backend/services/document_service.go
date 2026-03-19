package services

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"hrms/config"
	"hrms/models"
	"hrms/repositories"

	"github.com/google/uuid"
)

type DocumentService interface {
	UploadDocument(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, documentType string) (*models.Document, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error)
	GetFile(ctx context.Context, filePath string) (io.ReadCloser, error)
	GetAll(ctx context.Context, page, limit int) ([]*models.DocumentWithUser, int, error)
	DeleteDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type documentService struct {
	documentRepo   repositories.DocumentRepository
	storageService StorageService
}

func NewDocumentService(documentRepo repositories.DocumentRepository, storageService StorageService) DocumentService {
	return &documentService{
		documentRepo:   documentRepo,
		storageService: storageService,
	}
}

func (s *documentService) UploadDocument(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, documentType string) (*models.Document, error) {
	// Validate file size
	if file.Size > config.AppConfig.Storage.MaxFileSize {
		return nil, errors.New("file size exceeds maximum allowed size")
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	allowed := false
	for _, allowedExt := range config.AppConfig.Storage.AllowedTypes {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, errors.New("file type not allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	defer src.Close()

	// Generate unique filename
	filename := uuid.New().String() + ext

	// Save file
	filePath, err := s.storageService.SaveFile(filename, src)
	if err != nil {
		return nil, errors.New("failed to save file")
	}

	document := &models.Document{
		ID:           uuid.New(),
		UserID:       userID,
		FileURL:      filePath,
		FileName:     file.Filename,
		FileSize:     file.Size,
		DocumentType: documentType,
		UploadedAt:   time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.documentRepo.Create(ctx, document); err != nil {
		// Clean up file if database insert fails
		s.storageService.DeleteFile(filePath)
		return nil, errors.New("failed to create document record")
	}

	return document, nil
}

func (s *documentService) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Document, error) {
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
	return s.documentRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *documentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	return s.documentRepo.GetByID(ctx, id)
}

func (s *documentService) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return s.storageService.GetFile(filePath)
}

func (s *documentService) GetAll(ctx context.Context, page, limit int) ([]*models.DocumentWithUser, int, error) {
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
	documents, err := s.documentRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.documentRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

func (s *documentService) DeleteDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	document, err := s.documentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("document not found")
	}

	// Check ownership (unless admin)
	if document.UserID != userID {
		return errors.New("unauthorized to delete this document")
	}

	// Delete file
	if err := s.storageService.DeleteFile(document.FileURL); err != nil {
		// Log error but continue with database deletion
	}

	// Delete database record
	return s.documentRepo.Delete(ctx, id)
}
