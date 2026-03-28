package routes

import (
	"hrms/controllers"
	"hrms/middleware"
	"hrms/repositories"
	"hrms/services"
	"hrms/websocket"

	"hrms/database"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// WebSocket hub + handler
	wsHub := websocket.NewHub()
	wsHandler := websocket.NewHandler(wsHub)
	router.GET("/ws", wsHandler.HandleWebSocket)

	api := router.Group("/api")
	{
		// Initialize repositories
		userRepo := repositories.NewUserRepository()

		// Initialize audit service early (needed for middleware)
		auditRepo := repositories.NewAuditRepository()
		auditService := services.NewAuditService(auditRepo)

		// Initialize services
		authService := services.NewAuthService(userRepo)

		// Initialize controllers
		authController := controllers.NewAuthController(authService)

		// Apply audit middleware to all API routes (after auth middleware will be applied per route)
		// This will capture all POST/PUT/DELETE requests
		api.Use(middleware.AuditMiddleware(auditService))

		// Auth routes
		auth := api.Group("/auth")
		{
			// Register is now admin-only (or can be removed entirely)
			auth.POST("/register", middleware.RateLimitAuth(), middleware.AuthMiddleware(), middleware.RequireAdmin(), authController.Register)
			auth.POST("/login", middleware.RateLimitAuth(), authController.Login)
			auth.POST("/magic-login", middleware.RateLimitAuth(), authController.MagicLogin)
			auth.POST("/set-password", middleware.AuthMiddleware(), authController.SetPassword)
			auth.GET("/me", middleware.AuthMiddleware(), authController.GetMe)
		}

		// Employee routes
		emailService := services.NewEmailService()
		employeeService := services.NewEmployeeService(userRepo, repositories.NewEmployeeRepository(), emailService)
		employeeController := controllers.NewEmployeeController(employeeService)
		employees := api.Group("/employees")
		{
			employees.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), employeeController.ListEmployees)
			employees.POST("", middleware.AuthMiddleware(), middleware.RequireAdmin(), employeeController.CreateEmployee)
			employees.GET("/me", middleware.AuthMiddleware(), employeeController.GetMe)
			employees.GET("/:id", middleware.AuthMiddleware(), employeeController.GetEmployee)
			employees.PUT("/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), employeeController.UpdateEmployee)
		}
		// Attendance routes
		attendanceRepo := repositories.NewAttendanceRepository()
		attendanceService := services.NewAttendanceService(attendanceRepo)
		attendanceController := controllers.NewAttendanceController(attendanceService)
		attendance := api.Group("/attendance")
		{
			attendance.POST("/checkin", middleware.AuthMiddleware(), attendanceController.CheckIn)
			attendance.POST("/checkout", middleware.AuthMiddleware(), attendanceController.CheckOut)
			attendance.GET("/me", middleware.AuthMiddleware(), attendanceController.GetMyAttendance)
			attendance.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), attendanceController.GetAllAttendance)
		}
		// Leave routes
		leaveRepo := repositories.NewLeaveRepository()
		// Notification service is needed by leave service for real-time updates
		notificationRepo := repositories.NewNotificationRepository()
		notificationService := services.NewNotificationService(notificationRepo, emailService, wsHub)
		leaveService := services.NewLeaveService(leaveRepo, notificationService, wsHub, userRepo)
		leaveController := controllers.NewLeaveController(leaveService)
		leaves := api.Group("/leaves")
		{
			leaves.POST("/apply", middleware.AuthMiddleware(), leaveController.ApplyLeave)
			leaves.GET("/me", middleware.AuthMiddleware(), leaveController.GetMyLeaves)
			leaves.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), leaveController.GetAllLeaves)
			leaves.PUT("/:id/approve", middleware.AuthMiddleware(), middleware.RequireAdmin(), leaveController.ApproveLeave)
			leaves.PUT("/:id/reject", middleware.AuthMiddleware(), middleware.RequireAdmin(), leaveController.RejectLeave)
		}
		// Payroll routes
		salaryRepo := repositories.NewSalaryRepository()
		salaryService := services.NewSalaryService(salaryRepo)
		salaryController := controllers.NewSalaryController(salaryService)
		salary := api.Group("/salary")
		{
			salary.POST("", middleware.AuthMiddleware(), middleware.RequireAdmin(), salaryController.CreateSalary)
			salary.GET("/me", middleware.AuthMiddleware(), salaryController.GetMySalary)
			salary.GET("/:userId", middleware.AuthMiddleware(), middleware.RequireAdmin(), salaryController.GetSalaryByUserID)
		}
		// Document routes
		documentRepo := repositories.NewDocumentRepository()
		documentService := services.NewDocumentService(documentRepo, services.NewLocalStorageService())
		documentController := controllers.NewDocumentController(documentService)
		documents := api.Group("/documents")
		{
			documents.POST("", middleware.AuthMiddleware(), documentController.UploadDocument)
			documents.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), documentController.GetAllDocuments)
			documents.GET("/me", middleware.AuthMiddleware(), documentController.GetMyDocuments)
			documents.GET("/user/:userId", middleware.AuthMiddleware(), middleware.RequireAdmin(), documentController.GetDocumentsByUserID)
			documents.GET("/:id/download", middleware.AuthMiddleware(), documentController.DownloadDocument)
			documents.DELETE("/:id", middleware.AuthMiddleware(), documentController.DeleteDocument)
		}
		// Notification routes (notificationService already initialized above)
		notificationController := controllers.NewNotificationController(notificationService)
		notifications := api.Group("/notifications")
		{
			notifications.GET("", middleware.AuthMiddleware(), notificationController.GetNotifications)
			notifications.GET("/unread-count", middleware.AuthMiddleware(), notificationController.GetUnreadCount)
			notifications.PUT("/:id/read", middleware.AuthMiddleware(), notificationController.MarkAsRead)
		}

		// Dashboard routes
		dashboardService := services.NewDashboardService()
		dashboardController := controllers.NewDashboardController(dashboardService)
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/admin", middleware.AuthMiddleware(), middleware.RequireAdmin(), dashboardController.GetAdminDashboard)
			dashboard.GET("/employee", middleware.AuthMiddleware(), dashboardController.GetEmployeeDashboard)
		}

		// Audit routes (auditService already initialized above)
		auditController := controllers.NewAuditController(auditService)
		audit := api.Group("/audit-logs")
		{
			audit.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), auditController.GetAuditLogs)
		}

		// AI routes
		_ = database.DB // ensure database initialized elsewhere; injected into service
		aiService := services.NewAIService(
			userRepo,
			repositories.NewLeaveRepository(),
			repositories.NewAttendanceRepository(),
			repositories.NewSalaryRepository(),
		)
		aiController := controllers.NewAIController(aiService)
		ai := api.Group("/ai")
		{
			ai.POST("/hr-assistant", middleware.AuthMiddleware(), aiController.ProcessHRQuery)
		}
	}
}
