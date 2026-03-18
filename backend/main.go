package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"hrms/config"
	"hrms/database"
	"hrms/middleware"
	"hrms/routes"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Set Gin mode
	if config.AppConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize rate limiter cleanup
	middleware.CleanupVisitors()

	// Initialize New Relic app (optional)
	var nrApp *newrelic.Application
	if config.AppConfig.NewRelic.Enabled {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName(config.AppConfig.NewRelic.AppName),
			newrelic.ConfigLicense(config.AppConfig.NewRelic.LicenseKey),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigAIMonitoringEnabled(config.AppConfig.NewRelic.AIMonitoringEnabled),
			newrelic.ConfigCustomInsightsEventsMaxSamplesStored(config.AppConfig.NewRelic.CustomInsightsEventsMaxSamplesStored),
		)
		if err != nil {
			log.Printf("New Relic initialization failed; continuing without New Relic: %v", err)
		} else {
			nrApp = app
		}
	}

	// Setup router
	router := setupRouter(nrApp)

	// Start server
	serverAddr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on %s", serverAddr)

	// Graceful shutdown
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

func setupRouter(nrApp *newrelic.Application) *gin.Engine {
	router := gin.New()

	// Global middleware
	// New Relic middleware should be first to capture the full request lifecycle.
	router.Use(middleware.NewRelicMiddleware(nrApp))
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
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
