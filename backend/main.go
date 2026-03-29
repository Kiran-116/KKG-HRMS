package main

import (
	"os"
	"os/signal"
	"syscall"

	"hrms/config"
	"hrms/database"
	"hrms/logger"
	"hrms/middleware"
	"hrms/routes"

	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		// fall back to std output since logger not initialized yet
		println("Failed to load configuration:", err.Error())
		return
	}

	// Initialize database
	if err := database.Init(); err != nil {
		println("Failed to initialize database:", err.Error())
		return
	}
	defer database.Close()

	// Set Gin mode
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize rate limiter cleanup
	middleware.CleanupVisitors()

	// Initialize application logger first (used by New Relic agent for its own logs)
	logger.Init(config.AppConfig.Logging.Level, config.AppConfig.Logging.Pretty, nil)

	// Initialize New Relic app (optional)
	var nrApp *newrelic.Application
	if config.AppConfig.NewRelic.Enabled {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName(config.AppConfig.NewRelic.AppName),
			newrelic.ConfigLicense(config.AppConfig.NewRelic.LicenseKey),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigAIMonitoringEnabled(config.AppConfig.NewRelic.AIMonitoringEnabled),
			newrelic.ConfigCustomInsightsEventsMaxSamplesStored(config.AppConfig.NewRelic.CustomInsightsEventsMaxSamplesStored),
			newrelic.ConfigDistributedTracerEnabled(true),
			newrelic.ConfigDebugLogger(os.Stdout),
		)
		if err != nil {
			println("New Relic initialization failed; continuing without New Relic:", err.Error())
		} else {
			nrApp = app
		}
	}

	// Reconfigure logger (no-op for now; reserved if we add context-aware linking)
	logger.Init(config.AppConfig.Logging.Level, config.AppConfig.Logging.Pretty, nrApp)

	// Setup router
	router := setupRouter(nrApp)

	// Start server
	serverAddr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	logger.Logger.Info().Str("addr", serverAddr).Msg("server_starting")

	// Graceful shutdown
	go func() {
		if err := router.Run(serverAddr); err != nil {
			logger.Logger.Fatal().Err(err).Msg("failed_to_start_server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info().Msg("server_shutting_down")
}

func setupRouter(nrApp *newrelic.Application) *gin.Engine {
	router := gin.New()

	// Global middleware
	// New Relic middleware should be first to capture the full request lifecycle.
	if nrApp != nil {
		router.Use(nrgin.Middleware(nrApp))
	}
	router.Use(middleware.Recovery())
	// Use the structured request logger with NR context
	router.Use(middleware.RequestLoggingNR())
	router.Use(middleware.SetupCORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "HRMS API",
		})
	})

	// Setup routes
	routes.SetupRoutes(router)

	return router
}
