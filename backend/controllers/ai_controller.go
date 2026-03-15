package controllers

import (
	"net/http"

	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AIController struct {
	aiService services.AIService
}

func NewAIController(aiService services.AIService) *AIController {
	return &AIController{
		aiService: aiService,
	}
}

type HRQueryRequest struct {
	Query string `json:"query" binding:"required"`
}

// ProcessHRQuery processes HR-related queries
func (c *AIController) ProcessHRQuery(ctx *gin.Context) {
	var req HRQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	answer, err := c.aiService.ProcessHRQuery(id, req.Query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process query",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"answer": answer})
}
