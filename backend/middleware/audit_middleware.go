package middleware

import (
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditMiddleware creates audit logs for important actions
func AuditMiddleware(auditService services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only log POST, PUT, DELETE requests
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			return
		}

		// Get user ID if authenticated
		var userID interface{}
		if val, exists := c.Get("user_id"); exists {
			userID = val
		}

		// Extract action and entity type from path
		path := c.FullPath()
		action := method
		entityType := "unknown"

		// Determine entity type from path
		if path != "" {
			// Extract entity type from routes like /api/employees/:id
			if len(path) > 5 {
				parts := path[5:] // Remove /api/
				if len(parts) > 0 {
					entityType = parts
				}
			}
		}

		// Log audit event (async in production)
		go func() {
			var uid *uuid.UUID
			if userID != nil {
				if id, ok := userID.(uuid.UUID); ok {
					uid = &id
				}
			}
			auditService.Log(
				action,
				entityType,
				nil,
				uid,
				map[string]interface{}{
					"path":   path,
					"method": method,
				},
				c.ClientIP(),
				c.Request.UserAgent(),
			)
		}()
	}
}
