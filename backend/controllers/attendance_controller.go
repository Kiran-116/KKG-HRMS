package controllers

import (
	"net/http"
	"strconv"
	"time"

	"hrms/models"
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceController struct {
	attendanceService services.AttendanceService
}

func NewAttendanceController(attendanceService services.AttendanceService) *AttendanceController {
	return &AttendanceController{
		attendanceService: attendanceService,
	}
}

// CheckIn handles employee check-in
func (c *AttendanceController) CheckIn(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	var req models.CheckInRequest
	ctx.ShouldBindJSON(&req)

	date := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			date = parsed
		}
	}

	attendance, err := c.attendanceService.CheckIn(ctx.Request.Context(), id, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Check-in failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, attendance)
}

// CheckOut handles employee check-out
func (c *AttendanceController) CheckOut(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	var req models.CheckOutRequest
	ctx.ShouldBindJSON(&req)

	date := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			date = parsed
		}
	}

	attendance, err := c.attendanceService.CheckOut(ctx.Request.Context(), id, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Check-out failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, attendance)
}

// GetMyAttendance gets current user's attendance history
func (c *AttendanceController) GetMyAttendance(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	attendances, err := c.attendanceService.GetByUserID(ctx.Request.Context(), id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch attendance",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attendances": attendances})
}

// GetAllAttendance gets all attendance (admin only)
func (c *AttendanceController) GetAllAttendance(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var userID *uuid.UUID
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	var date *time.Time
	if dateStr := ctx.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = &parsed
		}
	}

	attendances, err := c.attendanceService.GetAll(ctx.Request.Context(), page, limit, userID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch attendance",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"attendances": attendances})
}
