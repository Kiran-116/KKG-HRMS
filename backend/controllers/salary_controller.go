package controllers

import (
	"net/http"
	"strconv"

	"hrms/models"
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SalaryController struct {
	salaryService services.SalaryService
}

func NewSalaryController(salaryService services.SalaryService) *SalaryController {
	return &SalaryController{
		salaryService: salaryService,
	}
}

// CreateSalary creates a salary record (admin only)
func (c *SalaryController) CreateSalary(ctx *gin.Context) {
	var req models.CreateSalaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	salary, err := c.salaryService.CreateSalary(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create salary",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, salary)
}

// GetMySalary gets current user's salary history
func (c *SalaryController) GetMySalary(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	salaries, err := c.salaryService.GetByUserID(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch salary",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"salaries": salaries})
}

// GetSalaryByUserID gets salary by user ID (admin only)
func (c *SalaryController) GetSalaryByUserID(ctx *gin.Context) {
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

	salaries, err := c.salaryService.GetByUserID(userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch salary",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"salaries": salaries})
}
