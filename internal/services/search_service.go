package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go-blog/internal/models"
	"go-blog/internal/repositories"
)

// SearchService handles search operations
type SearchService struct {
	articleRepo  repositories.ArticleRepository
	categoryRepo repositories.CategoryRepository
	tagRepo      repositories.TagRepository
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query      string    `json:"query" validate:"required,min=1,max=255"`
	CategoryID uint      `json:"category_id,omitempty"`
	AuthorID   uint      `json:"author_id,omitempty"`
	TagID      uint      `json:"tag_id,omitempty"`
	Status     string    `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	DateFrom   time.Time `json:"date_from,omitempty"`
	DateTo     time.Time `json:"date_to,omitempty"`
	SearchMode string    `json:"search_mode,omitempty" validate:"omitempty,oneof=natural boolean"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Articles    []models.Article `json:"articles"`
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	Limit       int              `json:"limit"`
	TotalPages  int              `json:"total_pages"`
	Query       string           `json:"query"`
	SearchTime  time.Duration    `json:"search_time_ms"`
	Suggestions []string         `json:"suggestions,omitempty"`
}

// SearchSuggestion represents a search suggestion
type SearchSuggestion struct {
	Term        string `json:"term"`
	Type        string `json:"type"` // "category", "tag", "author"
	Count       int64  `json:"count"`
	Description string `json:"description,omitempty"`
}

// NewSearchService creates a new search service
func NewSearchService(
	articleRepo repositories.ArticleRepository,
	categoryRepo repositories.CategoryRepository,
	tagRepo repositories.TagRepository,
) *SearchService {
	return &SearchService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// Search performs a search with the given parameters
func (s *SearchService) Search(req *SearchRequest, page, limit int) (*SearchResponse, error) {
	startTime := time.Now()

	// Validate input
	if err := s.validateSearchRequest(req); err != nil {
		return nil, err
	}

	// Sanitize and prepare query
	query := s.sanitizeQuery(req.Query)
	if query == "" {
		return nil, errors.New("search query cannot be empty after sanitization")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build search filters
	filters := &repositories.SearchFilters{
		Status:     req.Status,
		CategoryID: req.CategoryID,
		AuthorID:   req.AuthorID,
		TagID:      req.TagID,
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
	}

	// Perform search based on mode
	var articles []models.Article
	var total int64
	var err error

	switch req.SearchMode {
	case "boolean":
		// Prepare query for boolean search
		booleanQuery := s.prepareBooleanQuery(query)
		articles, total, err = s.articleRepo.SearchWithBoolean(booleanQuery, offset, limit, filters)
	default:
		// Default to natural language search
		articles, total, err = s.articleRepo.AdvancedSearch(query, offset, limit, filters)
	}

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Generate suggestions if no results found
	var suggestions []string
	if total == 0 {
		suggestions = s.generateSuggestions(query)
	}

	searchTime := time.Since(startTime)

	return &SearchResponse{
		Articles:    articles,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		Query:       req.Query,
		SearchTime:  searchTime,
		Suggestions: suggestions,
	}, nil
}

// QuickSearch performs a simple search without advanced filters
func (s *SearchService) QuickSearch(query string, page, limit int) (*SearchResponse, error) {
	req := &SearchRequest{
		Query:      query,
		SearchMode: "natural",
	}
	return s.Search(req, page, limit)
}

// GetSearchSuggestions returns search suggestions based on partial query
func (s *SearchService) GetSearchSuggestions(partialQuery string, limit int) ([]SearchSuggestion, error) {
	if limit < 1 || limit > 20 {
		limit = 10
	}

	suggestions := make([]SearchSuggestion, 0)

	// Get category suggestions
	categories, err := s.categoryRepo.List()
	if err == nil {
		for _, category := range categories {
			if strings.Contains(strings.ToLower(category.Name), strings.ToLower(partialQuery)) {
				suggestions = append(suggestions, SearchSuggestion{
					Term:        category.Name,
					Type:        "category",
					Description: fmt.Sprintf("Search in %s category", category.Name),
				})
			}
		}
	}

	// Get tag suggestions
	tags, err := s.tagRepo.List()
	if err == nil {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag.Name), strings.ToLower(partialQuery)) {
				suggestions = append(suggestions, SearchSuggestion{
					Term:        tag.Name,
					Type:        "tag",
					Description: fmt.Sprintf("Search articles tagged with %s", tag.Name),
				})
			}
		}
	}

	// Limit results
	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return suggestions, nil
}

// validateSearchRequest validates the search request
func (s *SearchService) validateSearchRequest(req *SearchRequest) error {
	if req == nil {
		return errors.New("search request is required")
	}

	if strings.TrimSpace(req.Query) == "" {
		return errors.New("search query is required")
	}

	if len(req.Query) > 255 {
		return errors.New("search query must be less than 255 characters")
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
			return errors.New("invalid status filter")
		}
	}

	if req.SearchMode != "" && req.SearchMode != "natural" && req.SearchMode != "boolean" {
		return errors.New("search mode must be 'natural' or 'boolean'")
	}

	if !req.DateFrom.IsZero() && !req.DateTo.IsZero() && req.DateFrom.After(req.DateTo) {
		return errors.New("date_from must be before date_to")
	}

	return nil
}

// sanitizeQuery sanitizes the search query
func (s *SearchService) sanitizeQuery(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// Remove excessive whitespace
	re := regexp.MustCompile(`\s+`)
	query = re.ReplaceAllString(query, " ")

	// Remove potentially dangerous characters for FULLTEXT search
	// Keep alphanumeric, spaces, and common punctuation
	re = regexp.MustCompile(`[^\w\s\-\.\,\!\?\:\;\(\)\[\]\"\']+`)
	query = re.ReplaceAllString(query, " ")

	// Trim again after sanitization
	query = strings.TrimSpace(query)

	return query
}

// prepareBooleanQuery prepares query for MySQL FULLTEXT Boolean search
func (s *SearchService) prepareBooleanQuery(query string) string {
	// Split query into words
	words := strings.Fields(query)
	if len(words) == 0 {
		return query
	}

	// For boolean search, we can add operators
	// For now, we'll just ensure each word is properly formatted
	var booleanTerms []string
	for _, word := range words {
		// Remove any existing boolean operators to prevent injection
		word = strings.Trim(word, "+-<>()~*\"")
		if word != "" {
			// Add word as-is for natural boolean search
			booleanTerms = append(booleanTerms, word)
		}
	}

	return strings.Join(booleanTerms, " ")
}

// generateSuggestions generates search suggestions when no results are found
func (s *SearchService) generateSuggestions(query string) []string {
	suggestions := make([]string, 0)

	// Simple suggestions based on common search patterns
	words := strings.Fields(strings.ToLower(query))
	if len(words) > 0 {
		// Suggest removing words
		if len(words) > 1 {
			for i := range words {
				newWords := make([]string, 0, len(words)-1)
				for j, word := range words {
					if i != j {
						newWords = append(newWords, word)
					}
				}
				suggestions = append(suggestions, strings.Join(newWords, " "))
			}
		}

		// Suggest alternative spellings (simple character removal/addition)
		firstWord := words[0]
		if len(firstWord) > 3 {
			// Remove last character
			suggestions = append(suggestions, firstWord[:len(firstWord)-1])
			// Remove first character
			suggestions = append(suggestions, firstWord[1:])
		}
	}

	// Remove duplicates and limit
	seen := make(map[string]bool)
	uniqueSuggestions := make([]string, 0)
	for _, suggestion := range suggestions {
		if !seen[suggestion] && suggestion != strings.ToLower(query) {
			seen[suggestion] = true
			uniqueSuggestions = append(uniqueSuggestions, suggestion)
		}
	}

	// Limit to 5 suggestions
	if len(uniqueSuggestions) > 5 {
		uniqueSuggestions = uniqueSuggestions[:5]
	}

	return uniqueSuggestions
}

// GetPopularSearchTerms returns popular search terms (placeholder for future implementation)
func (s *SearchService) GetPopularSearchTerms(limit int) ([]string, error) {
	// This would typically be implemented with search analytics
	// For now, return empty slice
	return []string{}, nil
}

// LogSearch logs a search query for analytics (placeholder for future implementation)
func (s *SearchService) LogSearch(query string, userID *uint, resultCount int64) error {
	// This would typically log to a search analytics table
	// For now, just return nil
	return nil
}