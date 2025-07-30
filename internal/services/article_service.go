package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"go-blog/internal/utils"

	"gorm.io/gorm"
)

type ArticleService struct {
	articleRepo  repositories.ArticleRepository
	userRepo     repositories.UserRepository
	categoryRepo repositories.CategoryRepository
	tagRepo      repositories.TagRepository
}

// CreateArticleRequest represents article creation data
type CreateArticleRequest struct {
	Title      string   `json:"title" validate:"required,min=1,max=255"`
	Content    string   `json:"content" validate:"required,min=1"`
	Excerpt    string   `json:"excerpt,omitempty" validate:"omitempty,max=500"`
	CategoryID *uint    `json:"category_id,omitempty" validate:"omitempty,min=1"`
	TagNames   []string `json:"tag_names,omitempty"`
	Status     string   `json:"status,omitempty" validate:"omitempty,oneof=draft published"`
}

// UpdateArticleRequest represents article update data
type UpdateArticleRequest struct {
	Title      string   `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content    string   `json:"content,omitempty" validate:"omitempty,min=1"`
	Excerpt    string   `json:"excerpt,omitempty" validate:"omitempty,max=500"`
	CategoryID *uint    `json:"category_id,omitempty" validate:"omitempty,min=1"`
	TagNames   []string `json:"tag_names,omitempty"`
	Status     string   `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

// ArticleListFilters represents filters for article listing
type ArticleListFilters struct {
	Status     string `json:"status,omitempty"`
	CategoryID uint   `json:"category_id,omitempty"`
	AuthorID   uint   `json:"author_id,omitempty"`
	TagID      uint   `json:"tag_id,omitempty"`
}

// NewArticleService creates a new article service
func NewArticleService(
	articleRepo repositories.ArticleRepository,
	userRepo repositories.UserRepository,
	categoryRepo repositories.CategoryRepository,
	tagRepo repositories.TagRepository,
) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// Create creates a new article
func (s *ArticleService) Create(authorID uint, req *CreateArticleRequest) (*models.Article, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Verify author exists
	author, err := s.userRepo.GetByID(authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	// Generate unique slug
	slug, err := s.generateUniqueSlug(req.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to generate slug: %w", err)
	}

	// Create article model
	article := &models.Article{
		Title:    strings.TrimSpace(req.Title),
		Slug:     slug,
		Content:  strings.TrimSpace(req.Content),
		Excerpt:  strings.TrimSpace(req.Excerpt),
		AuthorID: authorID,
		Author:   *author,
		Status:   models.StatusDraft, // Default to draft
	}

	// Set category if provided
	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		article.CategoryID = req.CategoryID
		article.Category = category
	}

	// Set status if provided
	if req.Status != "" {
		article.Status = models.ArticleStatus(req.Status)
	}

	// Handle publishing
	if article.Status == models.StatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	// Process tags
	if len(req.TagNames) > 0 {
		tags, err := s.processTagNames(req.TagNames)
		if err != nil {
			return nil, fmt.Errorf("failed to process tags: %w", err)
		}
		article.Tags = tags
	}

	// Create article
	if err := s.articleRepo.Create(article); err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return article, nil
}

// GetByID retrieves an article by ID
func (s *ArticleService) GetByID(id uint) (*models.Article, error) {
	return s.articleRepo.GetByID(id)
}

// GetBySlug retrieves an article by slug
func (s *ArticleService) GetBySlug(slug string) (*models.Article, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("slug cannot be empty")
	}
	return s.articleRepo.GetBySlug(slug)
}

// List retrieves articles with pagination and filters
func (s *ArticleService) List(page, limit int, filters *ArticleListFilters) ([]models.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	
	// Convert filters to map
	filterMap := make(map[string]interface{})
	if filters != nil {
		if filters.Status != "" {
			filterMap["status"] = filters.Status
		}
		if filters.CategoryID > 0 {
			filterMap["category_id"] = filters.CategoryID
		}
		if filters.AuthorID > 0 {
			filterMap["author_id"] = filters.AuthorID
		}
		if filters.TagID > 0 {
			filterMap["tag_id"] = filters.TagID
		}
	}

	return s.articleRepo.List(offset, limit, filterMap)
}

// Update updates an article
func (s *ArticleService) Update(id uint, authorID uint, req *UpdateArticleRequest) (*models.Article, error) {
	// Validate input
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing article
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	// Check authorization
	if article.AuthorID != authorID {
		return nil, errors.New("unauthorized: you can only edit your own articles")
	}

	// Update fields if provided
	updated := false

	if req.Title != "" && req.Title != article.Title {
		article.Title = strings.TrimSpace(req.Title)
		// Generate new slug if title changed
		slug, err := s.generateUniqueSlug(article.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to generate slug: %w", err)
		}
		article.Slug = slug
		updated = true
	}

	if req.Content != "" && req.Content != article.Content {
		article.Content = strings.TrimSpace(req.Content)
		updated = true
	}

	if req.Excerpt != article.Excerpt {
		article.Excerpt = strings.TrimSpace(req.Excerpt)
		updated = true
	}

	// Handle category change
	if req.CategoryID != nil {
		if (article.CategoryID == nil && *req.CategoryID != 0) ||
			(article.CategoryID != nil && *article.CategoryID != *req.CategoryID) {
			if *req.CategoryID == 0 {
				article.CategoryID = nil
				article.Category = nil
			} else {
				category, err := s.categoryRepo.GetByID(*req.CategoryID)
				if err != nil {
					return nil, errors.New("category not found")
				}
				article.CategoryID = req.CategoryID
				article.Category = category
			}
			updated = true
		}
	}

	// Handle status change
	if req.Status != "" && string(article.Status) != req.Status {
		oldStatus := article.Status
		article.Status = models.ArticleStatus(req.Status)

		// Handle publishing
		if article.Status == models.StatusPublished && oldStatus != models.StatusPublished {
			if article.PublishedAt == nil {
				now := time.Now()
				article.PublishedAt = &now
			}
		}
		updated = true
	}

	// Handle tags
	if req.TagNames != nil {
		tags, err := s.processTagNames(req.TagNames)
		if err != nil {
			return nil, fmt.Errorf("failed to process tags: %w", err)
		}
		article.Tags = tags
		updated = true
	}

	if !updated {
		return article, nil
	}

	// Update article
	if err := s.articleRepo.Update(article); err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return article, nil
}

// Delete deletes an article
func (s *ArticleService) Delete(id uint, authorID uint) error {
	// Get existing article
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	// Check authorization
	if article.AuthorID != authorID {
		return errors.New("unauthorized: you can only delete your own articles")
	}

	return s.articleRepo.Delete(id)
}

// Publish publishes an article
func (s *ArticleService) Publish(id uint, authorID uint) (*models.Article, error) {
	return s.changeStatus(id, authorID, models.StatusPublished)
}

// Unpublish unpublishes an article (sets to draft)
func (s *ArticleService) Unpublish(id uint, authorID uint) (*models.Article, error) {
	return s.changeStatus(id, authorID, models.StatusDraft)
}

// Archive archives an article
func (s *ArticleService) Archive(id uint, authorID uint) (*models.Article, error) {
	return s.changeStatus(id, authorID, models.StatusArchived)
}

// IncrementViewCount increments the view count for an article
func (s *ArticleService) IncrementViewCount(id uint) error {
	return s.articleRepo.IncrementViewCount(id)
}

// Search searches articles using basic search
func (s *ArticleService) Search(query string, page, limit int) ([]models.Article, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, errors.New("search query cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.articleRepo.Search(query, offset, limit)
}

// AdvancedSearch performs advanced search with filters
func (s *ArticleService) AdvancedSearch(query string, page, limit int, filters *repositories.SearchFilters) ([]models.Article, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, errors.New("search query cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.articleRepo.AdvancedSearch(query, offset, limit, filters)
}

// changeStatus changes the status of an article
func (s *ArticleService) changeStatus(id uint, authorID uint, status models.ArticleStatus) (*models.Article, error) {
	// Get existing article
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	// Check authorization
	if article.AuthorID != authorID {
		return nil, errors.New("unauthorized: you can only modify your own articles")
	}

	// Don't update if status is the same
	if article.Status == status {
		return article, nil
	}

	// Update status
	article.Status = status

	// Handle publishing timestamp
	if status == models.StatusPublished && article.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}

	// Update article
	if err := s.articleRepo.Update(article); err != nil {
		return nil, fmt.Errorf("failed to update article status: %w", err)
	}

	return article, nil
}

// generateUniqueSlug generates a unique slug for an article
func (s *ArticleService) generateUniqueSlug(title string) (string, error) {
	baseSlug := utils.GenerateSlug(title)
	if baseSlug == "" {
		return "", errors.New("cannot generate slug from title")
	}

	slug := baseSlug
	counter := 1

	// Check if slug exists and generate unique one
	for {
		_, err := s.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Slug is available
				break
			}
			return "", fmt.Errorf("error checking slug availability: %w", err)
		}

		// Slug exists, try with counter
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		// Prevent infinite loop
		if counter > 1000 {
			return "", errors.New("unable to generate unique slug")
		}
	}

	return slug, nil
}

// processTagNames processes tag names and returns tag models
func (s *ArticleService) processTagNames(tagNames []string) ([]models.Tag, error) {
	// Create a tag service instance to handle tag processing
	tagService := NewTagService(s.tagRepo)
	return tagService.ProcessTagNames(tagNames)
}

// validateCreateRequest validates article creation request
func (s *ArticleService) validateCreateRequest(req *CreateArticleRequest) error {
	if req == nil {
		return errors.New("create request is required")
	}

	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}

	if len(req.Title) > 255 {
		return errors.New("title must be less than 255 characters")
	}

	if strings.TrimSpace(req.Content) == "" {
		return errors.New("content is required")
	}

	if req.Excerpt != "" && len(req.Excerpt) > 500 {
		return errors.New("excerpt must be less than 500 characters")
	}

	if req.Status != "" && req.Status != "draft" && req.Status != "published" {
		return errors.New("status must be either 'draft' or 'published'")
	}

	return nil
}

// validateUpdateRequest validates article update request
func (s *ArticleService) validateUpdateRequest(req *UpdateArticleRequest) error {
	if req == nil {
		return errors.New("update request is required")
	}

	if req.Title != "" {
		if strings.TrimSpace(req.Title) == "" {
			return errors.New("title cannot be empty")
		}
		if len(req.Title) > 255 {
			return errors.New("title must be less than 255 characters")
		}
	}

	if req.Content != "" && strings.TrimSpace(req.Content) == "" {
		return errors.New("content cannot be empty")
	}

	if req.Excerpt != "" && len(req.Excerpt) > 500 {
		return errors.New("excerpt must be less than 500 characters")
	}

	if req.Status != "" {
		validStatuses := []string{"draft", "published", "archived"}
		valid := false
		for _, status := range validStatuses {
			if req.Status == status {
				valid = true
				break
			}
		}
		if !valid {
			return errors.New("status must be one of: draft, published, archived")
		}
	}

	return nil
}