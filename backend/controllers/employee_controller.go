package controllers

import (
	"net/http"
	"strconv"

	"hrms/models"
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeController struct {
	employeeService services.EmployeeService
}

func NewEmployeeController(employeeService services.EmployeeService) *EmployeeController {
	return &EmployeeController{
		employeeService: employeeService,
	}
}

// ListEmployees lists all employees (admin only)
func (c *EmployeeController) ListEmployees(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	response, err := c.employeeService.ListEmployees(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch employees",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateEmployee creates a new employee (admin only)
func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	employee, err := c.employeeService.CreateEmployee(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create employee",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, employee)
}

// GetEmployee gets employee by ID
func (c *EmployeeController) GetEmployee(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid employee ID",
			"message": err.Error(),
		})
		return
	}

	employee, err := c.employeeService.GetEmployee(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Employee not found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// UpdateEmployee updates an employee (admin only)
func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid employee ID",
			"message": err.Error(),
		})
		return
	}

	var req models.UpdateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	employee, err := c.employeeService.UpdateEmployee(ctx.Request.Context(), id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update employee",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// GetMe gets the current user's profile
func (c *EmployeeController) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal error",
			"message": "Invalid user ID type",
		})
		return
	}

	employee, err := c.employeeService.GetEmployee(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Employee not found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}
