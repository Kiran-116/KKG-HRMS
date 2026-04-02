package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hrms/services"
	"hrms/logger"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// responseWriter is a custom writer to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// normalizeEntityType converts entity type to user-friendly format
func normalizeEntityType(entityType string) string {
	// Handle special cases
	if entityType == "auth" {
		return "User"
	}

	// Convert to singular and capitalize
	entityType = strings.ToLower(entityType)
	if strings.HasSuffix(entityType, "ies") {
		return capitalizeFirst(strings.TrimSuffix(entityType, "ies") + "y")
	} else if strings.HasSuffix(entityType, "es") {
		return capitalizeFirst(strings.TrimSuffix(entityType, "es"))
	} else if strings.HasSuffix(entityType, "s") {
		return capitalizeFirst(strings.TrimSuffix(entityType, "s"))
	}
	return capitalizeFirst(entityType)
}

// generateDescription creates human-readable description from route and data
func generateDescription(action, entityType string, routePath string, requestBody map[string]interface{}, responseBody map[string]interface{}) string {
	normalizedEntity := normalizeEntityType(entityType)

	switch action {
	case "CREATE_EMPLOYEE":
		if name, ok := requestBody["name"].(string); ok {
			if email, ok := requestBody["email"].(string); ok {
				return "New Employee created: " + name + " (" + email + ")"
			}
			return "New Employee created: " + name
		}
		return "New Employee created"

	case "UPDATE_EMPLOYEE":
		changes := []string{}
		if name, ok := requestBody["name"].(string); ok {
			changes = append(changes, "Name: "+name)
		}
		if email, ok := requestBody["email"].(string); ok {
			changes = append(changes, "Email: "+email)
		}
		if len(changes) > 0 {
			return "Updated Employee: " + strings.Join(changes, ", ")
		}
		return "Updated Employee"

	case "CREATE_SALARY", "UPDATE_SALARY":
		if amount, ok := requestBody["amount"].(float64); ok {
			return "Salary record created: ₹" + formatCurrency(amount)
		}
		return "Salary record created"

	case "APPLY_LEAVE":
		if startDate, ok := requestBody["start_date"].(string); ok {
			if endDate, ok := requestBody["end_date"].(string); ok {
				return "Leave applied: " + startDate + " to " + endDate
			}
		}
		return "Leave applied"

	case "APPROVE_LEAVE":
		if leave, ok := responseBody["leave"].(map[string]interface{}); ok {
			if startDate, ok := leave["start_date"].(string); ok {
				if endDate, ok := leave["end_date"].(string); ok {
					return "Leave approved: " + startDate + " to " + endDate
				}
			}
		}
		return "Leave approved"

	case "REJECT_LEAVE":
		return "Leave rejected"

	case "CHECKIN":
		return "Checked in"

	case "CHECKOUT":
		return "Checked out"

	case "LOGIN":
		if email, ok := requestBody["email"].(string); ok {
			return "User logged in: " + email
		}
		return "User logged in"

	case "REGISTER":
		if email, ok := requestBody["email"].(string); ok {
			return "User registered: " + email
		}
		return "User registered"

	case "DELETE_DOCUMENT":
		if doc, ok := responseBody["document"].(map[string]interface{}); ok {
			if fileName, ok := doc["file_name"].(string); ok {
				return "Document deleted: " + fileName
			}
		}
		return "Document deleted"

	case "UPLOAD_DOCUMENT":
		if doc, ok := responseBody["document"].(map[string]interface{}); ok {
			if fileName, ok := doc["file_name"].(string); ok {
				return "Document uploaded: " + fileName
			}
		}
		return "Document uploaded"
	}

	// Generic descriptions
	if strings.HasPrefix(action, "CREATE_") {
		return "New " + normalizedEntity + " created"
	} else if strings.HasPrefix(action, "UPDATE_") {
		return "Updated " + normalizedEntity
	} else if strings.HasPrefix(action, "DELETE_") {
		return normalizedEntity + " deleted"
	}

	return action + " on " + normalizedEntity
}

func formatCurrency(amount float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", amount), "0"), ".")
}

// mapRouteToAction maps route patterns to descriptive action names
func mapRouteToAction(method, routePath string) (string, string) {
	routePath = strings.TrimPrefix(routePath, "/api/")
	parts := strings.Split(routePath, "/")

	if len(parts) == 0 {
		return method, "unknown"
	}

	entityType := parts[0]

	// Handle special routes
	if routePath == "attendance/checkin" {
		return "CHECKIN", "attendance"
	}
	if routePath == "attendance/checkout" {
		return "CHECKOUT", "attendance"
	}
	if routePath == "leaves/apply" {
		return "APPLY_LEAVE", "leaves"
	}
	if routePath == "auth/login" {
		return "LOGIN", "auth"
	}
	if routePath == "auth/register" {
		return "REGISTER", "auth"
	}

	// Handle routes with actions (e.g., /leaves/:id/approve)
	if len(parts) >= 3 {
		lastPart := parts[len(parts)-1]
		if lastPart == "approve" {
			return "APPROVE_LEAVE", "leaves"
		}
		if lastPart == "reject" {
			return "REJECT_LEAVE", "leaves"
		}
	}

	// Map HTTP methods to actions
	switch method {
	case "POST":
		switch entityType {
		case "employees":
			return "CREATE_EMPLOYEE", "employees"
		case "salary":
			return "CREATE_SALARY", "salary"
		case "documents":
			return "UPLOAD_DOCUMENT", "documents"
		}
	case "PUT":
		switch entityType {
		case "employees":
			return "UPDATE_EMPLOYEE", "employees"
		case "salary":
			return "UPDATE_SALARY", "salary"
		}
	case "DELETE":
		switch entityType {
		case "documents":
			return "DELETE_DOCUMENT", "documents"
		}
	}

	return method, entityType
}

// AuditMiddleware creates audit logs for important actions
func AuditMiddleware(auditService services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log POST, PUT, DELETE requests
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			c.Next()
			return
		}

		// Capture request body
		var requestBody map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			json.Unmarshal(bodyBytes, &requestBody)
		}

		// Capture response body
		responseBody := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = writer

		c.Next()

		// Get user ID if authenticated
		var userID interface{}
		if val, exists := c.Get("user_id"); exists {
			userID = val
		}

		// Extract action and entity type from path
		path := c.FullPath()
		requestPath := c.Request.URL.Path
		action, entityType := mapRouteToAction(method, path)
		var entityID *uuid.UUID

		// Extract entity ID from request path
		if path != "" && strings.HasPrefix(path, "/api/") {
			requestPathClean := requestPath
			if strings.HasPrefix(requestPath, "/api/") {
				requestPathClean = requestPath[5:]
			}
			requestParts := strings.Split(requestPathClean, "/")

			// Try to extract entity ID
			for _, part := range requestParts {
				if id, err := uuid.Parse(part); err == nil {
					entityID = &id
					break
				}
			}
		}

		// Parse response body
		var responseData map[string]interface{}
		if responseBody.Len() > 0 {
			json.Unmarshal(responseBody.Bytes(), &responseData)
		}

		// Generate description
		description := generateDescription(action, entityType, path, requestBody, responseData)

		// Normalize entity type
		normalizedEntityType := normalizeEntityType(entityType)

		// Capture values needed for persistence. (Do this in the main goroutine to avoid losing
		// inserts due to request-context cancellation or unsafe access to `c`.)
		var uid *uuid.UUID
		if userID != nil {
			if id, ok := userID.(uuid.UUID); ok {
				uid = &id
			}
		}

		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()
		metadata := map[string]interface{}{
			"path":   path,
			"method": method,
			"url":    requestPath,
		}

		if err := auditService.Log(
			c.Request.Context(),
			action,
			normalizedEntityType,
			entityID,
			uid,
			description,
			metadata,
			ipAddress,
			userAgent,
		); err != nil {
			logger.Logger.Error().
				Err(err).
				Str("action", action).
				Str("entity_type", normalizedEntityType).
				Str("url", requestPath).
				Str("method", method).
				Str("ip_address", ipAddress).
				Msg("audit_log_insert_failed")
		}
	}
}
