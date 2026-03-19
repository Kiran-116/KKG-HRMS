package controllers

import (
	"net/http"
	"strconv"

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
