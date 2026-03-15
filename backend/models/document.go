package models

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	FileURL      string    `json:"file_url"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	DocumentType string    `json:"document_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UploadDocumentRequest struct {
	DocumentType string `form:"document_type" binding:"required"`
}
