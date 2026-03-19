package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"hrms/models"
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentController struct {
	documentService services.DocumentService
}

func NewDocumentController(documentService services.DocumentService) *DocumentController {
	return &DocumentController{
		documentService: documentService,
	}
}

// UploadDocument handles document upload
func (c *DocumentController) UploadDocument(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	var req models.UploadDocumentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "File required",
			"message": "Please upload a file",
		})
		return
	}

	document, err := c.documentService.UploadDocument(ctx.Request.Context(), id, file, req.DocumentType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, document)
}

// GetMyDocuments gets current user's documents
func (c *DocumentController) GetMyDocuments(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	documents, err := c.documentService.GetByUserID(ctx.Request.Context(), id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch documents",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"documents": documents})
}

// GetAllDocuments gets all documents (admin only)
func (c *DocumentController) GetAllDocuments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	documents, total, err := c.documentService.GetAll(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch documents",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"documents": documents,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// GetDocumentsByUserID gets documents by user ID (admin only)
func (c *DocumentController) GetDocumentsByUserID(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": err.Error(),
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	documents, err := c.documentService.GetByUserID(ctx.Request.Context(), userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch documents",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"documents": documents})
}

// DownloadDocument downloads a document
func (c *DocumentController) DownloadDocument(ctx *gin.Context) {
	documentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid document ID",
			"message": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)
	userRole, _ := ctx.Get("user_role")
	role := userRole.(string)

	// Get document
	document, err := c.documentService.GetByID(ctx.Request.Context(), documentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Document not found",
			"message": err.Error(),
		})
		return
	}

	// Check ownership (user owns document OR user is admin)
	if document.UserID != id && role != models.RoleAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "You do not have permission to access this document",
		})
		return
	}

	// Get file from storage
	file, err := c.documentService.GetFile(ctx.Request.Context(), document.FileURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve file",
			"message": err.Error(),
		})
		return
	}
	defer file.Close()

	// Determine content type
	ext := strings.ToLower(filepath.Ext(document.FileName))
	contentType := getContentType(ext)

	// Set headers
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", `attachment; filename="`+document.FileName+`"`)

	// Stream file to response
	ctx.DataFromReader(http.StatusOK, document.FileSize, contentType, file, nil)
}

// getContentType returns the MIME type based on file extension
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

// DeleteDocument deletes a document
func (c *DocumentController) DeleteDocument(ctx *gin.Context) {
	documentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid document ID",
			"message": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	if err := c.documentService.DeleteDocument(ctx.Request.Context(), documentID, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete document",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
