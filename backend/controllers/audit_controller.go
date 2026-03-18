package controllers

import (
	"net/http"
	"strconv"

	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditController struct {
	auditService services.AuditService
}

func NewAuditController(auditService services.AuditService) *AuditController {
	return &AuditController{
		auditService: auditService,
	}
}

// GetAuditLogs gets audit logs (admin only)
func (c *AuditController) GetAuditLogs(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var userID *uuid.UUID
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	var action *string
	if actionStr := ctx.Query("action"); actionStr != "" {
		action = &actionStr
	}

	logs, err := c.auditService.GetAll(ctx.Request.Context(), page, limit, userID, action)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch audit logs",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"audit_logs": logs})
}
