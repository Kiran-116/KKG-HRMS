package controllers

import (
	"net/http"

	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardController struct {
	dashboardService services.DashboardService
}

func NewDashboardController(dashboardService services.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
	}
}

// GetAdminDashboard gets admin dashboard data
func (c *DashboardController) GetAdminDashboard(ctx *gin.Context) {
	data, err := c.dashboardService.GetAdminDashboard()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch dashboard data",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// GetEmployeeDashboard gets employee dashboard data
func (c *DashboardController) GetEmployeeDashboard(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	data, err := c.dashboardService.GetEmployeeDashboard(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch dashboard data",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
