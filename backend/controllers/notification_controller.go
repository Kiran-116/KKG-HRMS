package controllers

import (
	"net/http"
	"strconv"

	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {
	notificationService services.NotificationService
}

func NewNotificationController(notificationService services.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

// GetNotifications gets user's notifications
func (c *NotificationController) GetNotifications(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	notifications, err := c.notificationService.GetByUserID(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch notifications",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// GetUnreadCount gets unread notification count
func (c *NotificationController) GetUnreadCount(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	count, err := c.notificationService.GetUnreadCount(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch unread count",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

// MarkAsRead marks a notification as read
func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	notificationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid notification ID",
			"message": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	id := userID.(uuid.UUID)

	if err := c.notificationService.MarkAsRead(notificationID, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to mark as read",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
