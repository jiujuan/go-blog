package services

import (
	"errors"
	"strings"

	"go-blog/internal/models"
	"go-blog/internal/repositories"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo    repositories.UserRepository
	articleRepo repositories.ArticleRepository
}

// UpdateUserRequest represents user profile update data
type UpdateUserRequest struct {
	Username  string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email     string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"omitempty,url,max=255"`
	Bio       string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// SetArticleRepository sets the article repository (for dependency injection)
func (s *UserService) SetArticleRepository(articleRepo repositories.ArticleRepository) {
	s.articleRepo = articleRepo
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// Update updates a user
func (s *UserService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(userID uint, req *UpdateUserRequest) (*models.User, error) {
	// Validate input
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get current user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if username is being changed and if it's available
	if req.Username != "" && req.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(req.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("username is already taken")
		}
		user.Username = req.Username
	}

	// Check if email is being changed and if it's available
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("email is already taken")
		}
		user.Email = req.Email
	}

	// Update other fields
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user profile")
	}

	return user, nil
}

// GetUserArticles retrieves articles by user with pagination
func (s *UserService) GetUserArticles(userID uint, page, limit int) ([]*models.Article, int64, error) {
	if s.articleRepo == nil {
		return nil, 0, errors.New("article repository not available")
	}

	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	// Get articles by user
	offset := (page - 1) * limit
	articles, err := s.articleRepo.GetByAuthorID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := s.articleRepo.CountByAuthorID(userID)
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// validateUpdateRequest validates user update request
func (s *UserService) validateUpdateRequest(req *UpdateUserRequest) error {
	if req == nil {
		return errors.New("update request is required")
	}

	if req.Username != "" {
		if len(strings.TrimSpace(req.Username)) < 3 || len(req.Username) > 50 {
			return errors.New("username must be between 3 and 50 characters")
		}
	}

	if req.Email != "" {
		if len(strings.TrimSpace(req.Email)) == 0 {
			return errors.New("email cannot be empty")
		}
		if len(req.Email) > 100 {
			return errors.New("email must be less than 100 characters")
		}
	}

	if req.AvatarURL != "" && len(req.AvatarURL) > 255 {
		return errors.New("avatar URL must be less than 255 characters")
	}

	if req.Bio != "" && len(req.Bio) > 500 {
		return errors.New("bio must be less than 500 characters")
	}

	return nil
}