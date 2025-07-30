package services

import (
	"errors"
	"fmt"
	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"go-blog/internal/utils"
	"strings"

	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
	articleRepo  repositories.ArticleRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repositories.CategoryRepository, articleRepo repositories.ArticleRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		articleRepo:  articleRepo,
	}
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

// CategoryWithStats represents a category with article statistics
type CategoryWithStats struct {
	*models.Category
	ArticleCount int64 `json:"article_count"`
}

// Create creates a new category
func (s *CategoryService) Create(req *CreateCategoryRequest) (*models.Category, error) {
	if req == nil {
		return nil, errors.New("create request cannot be nil")
	}

	// Generate slug from name
	slug := utils.GenerateSlug(req.Name)
	if slug == "" {
		return nil, errors.New("failed to generate slug from category name")
	}

	// Check if category with same slug already exists
	existing, err := s.categoryRepo.GetBySlug(slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing category: %w", err)
	}
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	category := &models.Category{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Slug:        slug,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetByID retrieves a category by ID
func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	if id == 0 {
		return nil, errors.New("category ID cannot be zero")
	}

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// GetBySlug retrieves a category by slug
func (s *CategoryService) GetBySlug(slug string) (*models.Category, error) {
	if slug == "" {
		return nil, errors.New("category slug cannot be empty")
	}

	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// List retrieves all categories
func (s *CategoryService) List() ([]models.Category, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// ListWithStats retrieves all categories with article count statistics
func (s *CategoryService) ListWithStats() ([]CategoryWithStats, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	result := make([]CategoryWithStats, len(categories))
	for i, category := range categories {
		// Get article count for this category
		_, count, err := s.categoryRepo.GetArticles(category.ID, 0, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to get article count for category %d: %w", category.ID, err)
		}

		result[i] = CategoryWithStats{
			Category:     &categories[i],
			ArticleCount: count,
		}
	}

	return result, nil
}

// Update updates a category
func (s *CategoryService) Update(id uint, req *UpdateCategoryRequest) (*models.Category, error) {
	if id == 0 {
		return nil, errors.New("category ID cannot be zero")
	}
	if req == nil {
		return nil, errors.New("update request cannot be nil")
	}

	// Get existing category
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Generate new slug if name changed
	newSlug := utils.GenerateSlug(req.Name)
	if newSlug == "" {
		return nil, errors.New("failed to generate slug from category name")
	}

	// Check if another category with same slug exists (excluding current category)
	if newSlug != category.Slug {
		existing, err := s.categoryRepo.GetBySlug(newSlug)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check existing category: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("category with this name already exists")
		}
	}

	// Update category fields
	category.Name = strings.TrimSpace(req.Name)
	category.Description = strings.TrimSpace(req.Description)
	category.Slug = newSlug

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// Delete deletes a category
func (s *CategoryService) Delete(id uint) error {
	if id == 0 {
		return errors.New("category ID cannot be zero")
	}

	// Check if category exists
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return fmt.Errorf("failed to get category: %w", err)
	}

	// Check if category has articles
	_, count, err := s.categoryRepo.GetArticles(id, 0, 1)
	if err != nil {
		return fmt.Errorf("failed to check category articles: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category '%s' because it has %d articles", category.Name, count)
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// GetCategoryArticles retrieves articles for a specific category with pagination
func (s *CategoryService) GetCategoryArticles(categoryID uint, page, limit int) ([]models.Article, int64, error) {
	if categoryID == 0 {
		return nil, 0, errors.New("category ID cannot be zero")
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Check if category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("category not found")
		}
		return nil, 0, fmt.Errorf("failed to get category: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	articles, total, err := s.categoryRepo.GetArticles(categoryID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get category articles: %w", err)
	}

	return articles, total, nil
}

// GetCategoryArticlesBySlug retrieves articles for a specific category by slug with pagination
func (s *CategoryService) GetCategoryArticlesBySlug(slug string, page, limit int) ([]models.Article, int64, error) {
	if slug == "" {
		return nil, 0, errors.New("category slug cannot be empty")
	}

	// Get category by slug
	category, err := s.GetBySlug(slug)
	if err != nil {
		return nil, 0, err
	}

	return s.GetCategoryArticles(category.ID, page, limit)
}