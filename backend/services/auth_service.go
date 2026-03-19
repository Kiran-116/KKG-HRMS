package services

import (
	"context"
	"errors"
	"time"

	"hrms/models"
	"hrms/repositories"
	"hrms/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	MagicLogin(ctx context.Context, req *models.MagicLoginRequest) (*models.LoginResponse, error)
	SetPassword(ctx context.Context, userID uuid.UUID, req *models.SetPasswordRequest) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = models.RoleEmployee
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     true,
	}

	if req.Department != "" {
		user.Department = &req.Department
	}
	if req.Designation != "" {
		user.Designation = &req.Designation
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Minute * 15 / time.Second),
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.LoginResponse{
		User:              user,
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		TokenType:         "Bearer",
		ExpiresIn:         int64(time.Minute * 15 / time.Second),
		MustChangePassword: user.MustChangePassword,
	}, nil
}

func (s *authService) MagicLogin(ctx context.Context, req *models.MagicLoginRequest) (*models.LoginResponse, error) {
	// Get user by magic token
	user, err := s.userRepo.GetByMagicToken(ctx, req.Token)
	if err != nil {
		return nil, errors.New("invalid or expired magic link")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Check if must_change_password is true
	if !user.MustChangePassword {
		return nil, errors.New("magic link has already been used")
	}

	// Clear magic token (one-time use)
	emptyToken := ""
	emptyTime := time.Time{}
	user.MagicToken = &emptyToken
	user.MagicExpiresAt = &emptyTime
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.New("failed to update user")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.LoginResponse{
		User:              user,
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		TokenType:         "Bearer",
		ExpiresIn:         int64(time.Minute * 15 / time.Second),
		MustChangePassword: true,
	}, nil
}

func (s *authService) SetPassword(ctx context.Context, userID uuid.UUID, req *models.SetPasswordRequest) error {
	// Validate password match
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	// Validate stronger password policy (min 10 chars, number + symbol)
	if len(req.NewPassword) < 10 {
		return errors.New("password must be at least 10 characters long")
	}

	hasNumber := false
	hasSymbol := false
	for _, char := range req.NewPassword {
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if (char >= '!' && char <= '/') || (char >= ':' && char <= '@') || (char >= '[' && char <= '`') || (char >= '{' && char <= '~') {
			hasSymbol = true
		}
	}

	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSymbol {
		return errors.New("password must contain at least one symbol")
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password and clear must_change_password flag
	if err := s.userRepo.UpdatePassword(ctx, userID, passwordHash); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
