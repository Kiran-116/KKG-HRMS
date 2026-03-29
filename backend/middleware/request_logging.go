package middleware

import (
	"time"

	"hrms/logger"

	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
)

// RequestLoggingNR logs one structured line per HTTP request and, when available,
// includes New Relic correlation (trace/span) via the nrgin transaction.
func RequestLoggingNR() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Derive a per-request logger bound to NR context if available.
		reqLogger := deriveLogger(c)

		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()

		// Attempt to fetch NR transaction for correlation
		var traceID, spanID, entityGUID string
		if txn := nrgin.Transaction(c); txn != nil {
			if md := txn.GetLinkingMetadata(); md != (newrelic.LinkingMetadata{}) {
				traceID = md.TraceID
				spanID = md.SpanID
				entityGUID = md.EntityGUID
			}
		}

		reqLogger.
			Info().
			Int("status", status).
			Str("method", method).
			Str("path", path).
			Str("client_ip", clientIP).
			Str("user_agent", ua).
			Int64("duration_ms", latency.Milliseconds()).
			Str("trace.id", traceID).
			Str("span.id", spanID).
			Str("entity.guid", entityGUID).
			Msg("http_request")
	}
}

func deriveLogger(c *gin.Context) zerolog.Logger {
	if txn := nrgin.Transaction(c); txn != nil {
		// Use WithTransaction for logs-in-context correlation
		if md := txn.GetLinkingMetadata(); md != (newrelic.LinkingMetadata{}) {
			// enrich fields even if writer fails
			return logger.WithContext(newrelic.NewContext(c.Request.Context(), txn)).
				With().
				Str("trace.id", md.TraceID).
				Str("span.id", md.SpanID).
				Str("entity.guid", md.EntityGUID).
				Logger()
		}
		return logger.WithContext(newrelic.NewContext(c.Request.Context(), txn))
	}
	return logger.WithContext(c.Request.Context())
}
