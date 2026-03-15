package middleware

import (
	"net/http"

	"hrms/models"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		hasRole := false
		for _, r := range roles {
			if role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only routes
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireEmployee is a convenience middleware for employee routes
func RequireEmployee() gin.HandlerFunc {
	return RequireRole(models.RoleEmployee, models.RoleAdmin)
}
