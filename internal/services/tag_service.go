package services

import (
	"errors"
	"fmt"
	"strings"

	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"go-blog/internal/utils"

	"gorm.io/gorm"
)

type TagService struct {
	tagRepo repositories.TagRepository
}

// CreateTagRequest represents tag creation data
type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

// TagWithStats represents a tag with usage statistics
type TagWithStats struct {
	models.Tag
	ArticleCount int64 `json:"article_count"`
}

// NewTagService creates a new tag service
func NewTagService(tagRepo repositories.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// Create creates a new tag
func (s *TagService) Create(req *CreateTagRequest) (*models.Tag, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	tagName := strings.TrimSpace(req.Name)
	
	// Check if tag already exists
	existingTag, err := s.tagRepo.GetByName(tagName)
	if err == nil {
		return nil, fmt.Errorf("tag with name '%s' already exists", tagName)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing tag: %w", err)
	}

	// Generate unique slug
	slug, err := s.generateUniqueSlug(tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate slug: %w", err)
	}

	tag := &models.Tag{
		Name: tagName,
		Slug: slug,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

// GetByID retrieves a tag by ID
func (s *TagService) GetByID(id uint) (*models.Tag, error) {
	if id == 0 {
		return nil, errors.New("tag ID is required")
	}
	return s.tagRepo.GetByID(id)
}

// GetBySlug retrieves a tag by slug
func (s *TagService) GetBySlug(slug string) (*models.Tag, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("slug cannot be empty")
	}
	return s.tagRepo.GetBySlug(slug)
}

// GetByName retrieves a tag by name
func (s *TagService) GetByName(name string) (*models.Tag, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}
	return s.tagRepo.GetByName(name)
}

// List retrieves all tags
func (s *TagService) List() ([]models.Tag, error) {
	return s.tagRepo.List()
}

// FindOrCreateByName finds an existing tag by name or creates a new one
func (s *TagService) FindOrCreateByName(name string) (*models.Tag, error) {
	tagName := strings.TrimSpace(name)
	if tagName == "" {
		return nil, errors.New("tag name cannot be empty")
	}

	// Try to find existing tag
	tag, err := s.tagRepo.GetByName(tagName)
	if err == nil {
		return tag, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing tag: %w", err)
	}

	// Create new tag
	req := &CreateTagRequest{Name: tagName}
	return s.Create(req)
}

// GetArticles retrieves articles associated with a tag
func (s *TagService) GetArticles(tagID uint, page, limit int) ([]models.Article, int64, error) {
	if tagID == 0 {
		return nil, 0, errors.New("tag ID is required")
	}

	// Verify tag exists
	_, err := s.tagRepo.GetByID(tagID)
	if err != nil {
		return nil, 0, errors.New("tag not found")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.tagRepo.GetArticles(tagID, offset, limit)
}

// GetArticlesBySlug retrieves articles associated with a tag by slug
func (s *TagService) GetArticlesBySlug(slug string, page, limit int) ([]models.Article, int64, error) {
	tag, err := s.GetBySlug(slug)
	if err != nil {
		return nil, 0, err
	}

	return s.GetArticles(tag.ID, page, limit)
}

// ProcessTagNames processes a list of tag names and returns tag models
// This method finds existing tags or creates new ones as needed
func (s *TagService) ProcessTagNames(tagNames []string) ([]models.Tag, error) {
	var tags []models.Tag
	processedNames := make(map[string]bool) // To avoid duplicates

	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// Skip if already processed (avoid duplicates)
		if processedNames[tagName] {
			continue
		}
		processedNames[tagName] = true

		tag, err := s.FindOrCreateByName(tagName)
		if err != nil {
			return nil, fmt.Errorf("failed to process tag '%s': %w", tagName, err)
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

// GetPopularTags retrieves tags ordered by usage (article count)
func (s *TagService) GetPopularTags(limit int) ([]TagWithStats, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}

	tags, err := s.tagRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	// Get article counts for each tag
	var tagsWithStats []TagWithStats
	for _, tag := range tags {
		_, count, err := s.tagRepo.GetArticles(tag.ID, 0, 1) // Just get count
		if err != nil {
			return nil, fmt.Errorf("failed to get article count for tag %d: %w", tag.ID, err)
		}

		tagsWithStats = append(tagsWithStats, TagWithStats{
			Tag:          tag,
			ArticleCount: count,
		})
	}

	// Sort by article count (simple bubble sort for small datasets)
	for i := 0; i < len(tagsWithStats)-1; i++ {
		for j := 0; j < len(tagsWithStats)-i-1; j++ {
			if tagsWithStats[j].ArticleCount < tagsWithStats[j+1].ArticleCount {
				tagsWithStats[j], tagsWithStats[j+1] = tagsWithStats[j+1], tagsWithStats[j]
			}
		}
	}

	// Limit results
	if len(tagsWithStats) > limit {
		tagsWithStats = tagsWithStats[:limit]
	}

	return tagsWithStats, nil
}

// generateUniqueSlug generates a unique slug for a tag
func (s *TagService) generateUniqueSlug(name string) (string, error) {
	baseSlug := utils.GenerateSlug(name)
	if baseSlug == "" {
		return "", errors.New("cannot generate slug from name")
	}

	slug := baseSlug
	counter := 1

	// Check if slug exists and generate unique one
	for {
		_, err := s.tagRepo.GetBySlug(slug)
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

// validateCreateRequest validates tag creation request
func (s *TagService) validateCreateRequest(req *CreateTagRequest) error {
	if req == nil {
		return errors.New("create request is required")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New("tag name is required")
	}

	if len(req.Name) > 50 {
		return errors.New("tag name must be less than 50 characters")
	}

	return nil
}