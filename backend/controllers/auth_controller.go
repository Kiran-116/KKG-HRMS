package controllers

import (
	"net/http"

	"hrms/models"
	"hrms/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Registration failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := c.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMe returns the current authenticated user
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/me [get]
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal error",
			"message": "Invalid user ID type",
		})
		return
	}

	user, err := c.authService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Not found",
			"message": "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// MagicLogin handles magic link authentication
func (c *AuthController) MagicLogin(ctx *gin.Context) {
	var req models.MagicLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	response, err := c.authService.MagicLogin(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// SetPassword handles password setup for new users
func (c *AuthController) SetPassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal error",
			"message": "Invalid user ID type",
		})
		return
	}

	var req models.SetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": err.Error(),
		})
		return
	}

	if err := c.authService.SetPassword(ctx.Request.Context(), id, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to set password",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password set successfully",
	})
}