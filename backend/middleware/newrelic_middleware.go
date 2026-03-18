package middleware

import (
	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewRelicMiddleware provides New Relic instrumentation
func NewRelicMiddleware(app *newrelic.Application) gin.HandlerFunc {
	// When New Relic is disabled (or failed to initialize), keep behavior unchanged.
	if app == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return nrgin.Middleware(app)
}
