package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
)

// Logger is the process-wide logger instance.
var Logger zerolog.Logger
var baseWriter zerologWriter.ZerologWriter

// Init configures the global zerolog Logger.
// If app is non-nil, logs can be enriched with New Relic linking metadata
// by using the WithContext helper at call sites that have a context.
func Init(levelString string, prettyConsole bool, app *newrelic.Application) {
	level := parseLevel(levelString)

	var base zerolog.Logger
	if prettyConsole {
		consoleWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stdout
			w.TimeFormat = time.RFC3339
		})
		base = zerolog.New(consoleWriter).With().Timestamp().Logger()
	} else {
		base = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(level)

	// Use New Relic zerologWriter to enable logs-in-context.
	// Fallback to stdout if app is nil (writer requires an initialized app).
	if app != nil {
		bw := zerologWriter.New(os.Stdout, app)
		baseWriter = bw
		base = zerolog.New(bw).With().Timestamp().Logger()
	}

	// Attach static service metadata.
	l := base.With().Str("service", "HRMS API").Logger()
	Logger = l
}

// WithContext returns a child logger enriched with New Relic linking metadata
// when the provided context carries a transaction/span. Safe to call with nil.
func WithContext(ctx context.Context) zerolog.Logger {
	if ctx == nil || (baseWriter == zerologWriter.ZerologWriter{}) {
		return Logger
	}
	// Create a context-aware writer and return a logger using it.
	txnWriter := baseWriter.WithContext(ctx)
	return Logger.Output(txnWriter)
}

func parseLevel(s string) zerolog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
