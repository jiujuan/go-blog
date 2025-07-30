package services

import (
	"errors"
	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"go-blog/internal/utils"
	"strings"

	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

// LoginRequest represents user login data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User   *models.User       `json:"user"`
	Tokens *utils.TokenPair   `json:"tokens"`
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists by email
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if username is taken
	existingUser, err = s.userRepo.GetByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username is already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	tokens, err := utils.GenerateTokenPair(user.ID, user.Username, user.Email, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Validate input
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	// Find user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	tokens, err := utils.GenerateTokenPair(user.ID, user.Username, user.Email, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	return utils.ExtractUserID(tokenString, s.jwtSecret)
}

// GetUserFromToken validates token and returns user information
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := utils.ValidateJWT(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// RefreshToken generates new tokens using refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*utils.TokenPair, error) {
	claims, err := utils.ValidateJWT(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Verify user still exists
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new token pair
	tokens, err := utils.GenerateTokenPair(user.ID, user.Username, user.Email, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	return tokens, nil
}

// validateRegisterRequest validates registration request
func (s *AuthService) validateRegisterRequest(req *RegisterRequest) error {
	if req == nil {
		return errors.New("registration request is required")
	}

	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username is required")
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}

	if len(req.Email) > 100 {
		return errors.New("email must be less than 100 characters")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 || len(req.Password) > 255 {
		return errors.New("password must be between 8 and 255 characters")
	}

	return nil
}

// validateLoginRequest validates login request
func (s *AuthService) validateLoginRequest(req *LoginRequest) error {
	if req == nil {
		return errors.New("login request is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	return nil
}