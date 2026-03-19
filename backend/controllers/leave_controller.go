package controllers

import (
	"net/http"
	"strconv"

	"hrms/models"
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LeaveController struct {
	leaveService services.LeaveService
}

func NewLeaveController(leaveService services.LeaveService) *LeaveController {
	return &LeaveController{
		leaveService: leaveService,
	}
}

// ApplyLeave handles leave application
func (c *LeaveController) ApplyLeave(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	var req models.ApplyLeaveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	leave, err := c.leaveService.ApplyLeave(ctx.Request.Context(), id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to apply leave",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, leave)
}

// GetMyLeaves gets current user's leave history
func (c *LeaveController) GetMyLeaves(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	leaves, err := c.leaveService.GetByUserID(ctx.Request.Context(), id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch leaves",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

// GetAllLeaves gets all leaves (admin only)
func (c *LeaveController) GetAllLeaves(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var status *string
	if statusStr := ctx.Query("status"); statusStr != "" {
		status = &statusStr
	}

	leaves, err := c.leaveService.GetAll(ctx.Request.Context(), page, limit, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch leaves",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

// ApproveLeave approves a leave request (admin only)
func (c *LeaveController) ApproveLeave(ctx *gin.Context) {
	leaveID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid leave ID",
			"message": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	approvedBy := userID.(uuid.UUID)

	leave, err := c.leaveService.ApproveLeave(ctx.Request.Context(), leaveID, approvedBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to approve leave",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, leave)
}

// RejectLeave rejects a leave request (admin only)
func (c *LeaveController) RejectLeave(ctx *gin.Context) {
	leaveID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid leave ID",
			"message": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	approvedBy := userID.(uuid.UUID)

	leave, err := c.leaveService.RejectLeave(ctx.Request.Context(), leaveID, approvedBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to reject leave",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, leave)
}
