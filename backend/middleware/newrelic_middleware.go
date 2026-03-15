package middleware

import (
	"github.com/gin-gonic/gin"
)

// NewRelicMiddleware provides New Relic instrumentation
// This is a placeholder - actual implementation would use New Relic Go agent
func NewRelicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In production, this would start a New Relic transaction
		// transaction := newrelic.StartTransaction(c.Request.URL.Path, c.Writer, c.Request)
		// defer transaction.End()

		c.Next()

		// Record metrics
		// transaction.AddAttribute("status_code", c.Writer.Status())
		// transaction.AddAttribute("method", c.Request.Method)
	}
}
