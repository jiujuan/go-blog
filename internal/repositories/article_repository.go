package repositories

import (
	"time"

	"go-blog/internal/database"
	"go-blog/internal/models"
)

// SearchFilters represents advanced search filters
type SearchFilters struct {
	Status     string    `json:"status,omitempty"`
	CategoryID uint      `json:"category_id,omitempty"`
	AuthorID   uint      `json:"author_id,omitempty"`
	TagID      uint      `json:"tag_id,omitempty"`
	DateFrom   time.Time `json:"date_from,omitempty"`
	DateTo     time.Time `json:"date_to,omitempty"`
}

type articleRepository struct {
	*BaseRepository
}

// NewArticleRepository creates a new article repository
func NewArticleRepository(db *database.DB) ArticleRepository {
	return &articleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *articleRepository) Create(article *models.Article) error {
	return r.BaseRepository.Create(article)
}

func (r *articleRepository) GetByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.BaseRepository.GetByID(&article, id, "Author", "Category", "Tags")
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.GetDB().GetByField(&article, "slug", slug, "Author", "Category", "Tags")
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) List(offset, limit int, filters map[string]interface{}) ([]models.Article, int64, error) {
	var articles []models.Article
	
	// Convert offset/limit to page-based pagination
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	
	options := &database.QueryOptions{
		Page:     page,
		Limit:    limit,
		OrderBy:  "created_at DESC",
		Filters:  filters,
		Preloads: []string{"Author", "Category", "Tags"},
	}
	
	result, err := r.BaseRepository.List(&articles, options)
	if err != nil {
		return nil, 0, err
	}
	
	return articles, result.Total, nil
}

func (r *articleRepository) Update(article *models.Article) error {
	return r.BaseRepository.Update(article)
}

func (r *articleRepository) Delete(id uint) error {
	return r.BaseRepository.Delete(&models.Article{}, id)
}

func (r *articleRepository) Search(query string, offset, limit int) ([]models.Article, int64, error) {
	return r.AdvancedSearch(query, offset, limit, nil)
}

func (r *articleRepository) AdvancedSearch(query string, offset, limit int, filters *SearchFilters) ([]models.Article, int64, error) {
	var articles []models.Article
	
	// Convert offset/limit to page-based pagination
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	
	// Build the search query using MySQL FULLTEXT search
	db := r.GetDB().GetDB()
	
	// Start with base query
	searchQuery := db.Model(&models.Article{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags")
	
	// Apply FULLTEXT search if query is provided
	if query != "" {
		// Use MySQL FULLTEXT search with relevance scoring
		searchQuery = searchQuery.Where("MATCH(title, content, excerpt) AGAINST(? IN NATURAL LANGUAGE MODE)", query).
			Select("*, MATCH(title, content, excerpt) AGAINST(? IN NATURAL LANGUAGE MODE) as relevance_score", query).
			Order("relevance_score DESC, created_at DESC")
	} else {
		searchQuery = searchQuery.Order("created_at DESC")
	}
	
	// Apply filters if provided
	if filters != nil {
		if filters.Status != "" {
			searchQuery = searchQuery.Where("status = ?", filters.Status)
		}
		if filters.CategoryID > 0 {
			searchQuery = searchQuery.Where("category_id = ?", filters.CategoryID)
		}
		if filters.AuthorID > 0 {
			searchQuery = searchQuery.Where("author_id = ?", filters.AuthorID)
		}
		if filters.TagID > 0 {
			searchQuery = searchQuery.Joins("JOIN article_tags ON articles.id = article_tags.article_id").
				Where("article_tags.tag_id = ?", filters.TagID)
		}
		if !filters.DateFrom.IsZero() {
			searchQuery = searchQuery.Where("created_at >= ?", filters.DateFrom)
		}
		if !filters.DateTo.IsZero() {
			searchQuery = searchQuery.Where("created_at <= ?", filters.DateTo)
		}
	}
	
	// Count total results
	var total int64
	countQuery := searchQuery
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if err := searchQuery.Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	
	return articles, total, nil
}

func (r *articleRepository) SearchWithBoolean(query string, offset, limit int, filters *SearchFilters) ([]models.Article, int64, error) {
	var articles []models.Article
	
	// Convert offset/limit to page-based pagination
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	
	// Build the search query using MySQL FULLTEXT Boolean search
	db := r.GetDB().GetDB()
	
	// Start with base query
	searchQuery := db.Model(&models.Article{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags")
	
	// Apply FULLTEXT Boolean search if query is provided
	if query != "" {
		// Use MySQL FULLTEXT Boolean search for advanced operators
		searchQuery = searchQuery.Where("MATCH(title, content, excerpt) AGAINST(? IN BOOLEAN MODE)", query).
			Select("*, MATCH(title, content, excerpt) AGAINST(? IN BOOLEAN MODE) as relevance_score", query).
			Order("relevance_score DESC, created_at DESC")
	} else {
		searchQuery = searchQuery.Order("created_at DESC")
	}
	
	// Apply filters if provided
	if filters != nil {
		if filters.Status != "" {
			searchQuery = searchQuery.Where("status = ?", filters.Status)
		}
		if filters.CategoryID > 0 {
			searchQuery = searchQuery.Where("category_id = ?", filters.CategoryID)
		}
		if filters.AuthorID > 0 {
			searchQuery = searchQuery.Where("author_id = ?", filters.AuthorID)
		}
		if filters.TagID > 0 {
			searchQuery = searchQuery.Joins("JOIN article_tags ON articles.id = article_tags.article_id").
				Where("article_tags.tag_id = ?", filters.TagID)
		}
		if !filters.DateFrom.IsZero() {
			searchQuery = searchQuery.Where("created_at >= ?", filters.DateFrom)
		}
		if !filters.DateTo.IsZero() {
			searchQuery = searchQuery.Where("created_at <= ?", filters.DateTo)
		}
	}
	
	// Count total results
	var total int64
	countQuery := searchQuery
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if err := searchQuery.Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	
	return articles, total, nil
}

func (r *articleRepository) GetArchive() (map[string]interface{}, error) {
	db := r.GetDB().GetDB()
	
	// Query to get article counts grouped by year and month
	// Only include published articles
	var results []struct {
		Year  int   `json:"year"`
		Month int   `json:"month"`
		Count int64 `json:"count"`
	}
	
	query := `
		SELECT 
			YEAR(published_at) as year,
			MONTH(published_at) as month,
			COUNT(*) as count
		FROM articles 
		WHERE status = 'published' 
			AND published_at IS NOT NULL
			AND deleted_at IS NULL
		GROUP BY YEAR(published_at), MONTH(published_at)
		ORDER BY year DESC, month DESC
	`
	
	if err := db.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	
	// Organize results by year
	yearMap := make(map[int][]map[string]interface{})
	
	for _, result := range results {
		if _, exists := yearMap[result.Year]; !exists {
			yearMap[result.Year] = []map[string]interface{}{}
		}
		
		monthData := map[string]interface{}{
			"month": int64(result.Month),
			"count": result.Count,
		}
		
		yearMap[result.Year] = append(yearMap[result.Year], monthData)
	}
	
	// Convert to the expected format
	var years []map[string]interface{}
	
	for year, months := range yearMap {
		yearData := map[string]interface{}{
			"year":   int64(year),
			"months": months,
		}
		years = append(years, yearData)
	}
	
	response := map[string]interface{}{
		"years": years,
	}
	
	return response, nil
}

func (r *articleRepository) GetByMonth(year, month int, offset, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	
	db := r.GetDB().GetDB()
	
	// Create date range for the specified month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second) // Last second of the month
	
	// Build query for articles in the specified month
	query := db.Model(&models.Article{}).
		Where("status = ?", "published").
		Where("published_at >= ? AND published_at <= ?", startDate, endDate).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC")
	
	// Count total results
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and get results
	if err := query.Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	
	return articles, total, nil
}

func (r *articleRepository) GetByAuthorID(authorID uint, limit, offset int) ([]*models.Article, error) {
	var articles []*models.Article
	
	// Convert offset/limit to page-based pagination
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	
	filters := map[string]interface{}{
		"author_id": authorID,
		"status":    "published", // Only return published articles
	}
	
	options := &database.QueryOptions{
		Page:     page,
		Limit:    limit,
		OrderBy:  "created_at DESC",
		Filters:  filters,
		Preloads: []string{"Author", "Category", "Tags"},
	}
	
	result, err := r.BaseRepository.List(&articles, options)
	if err != nil {
		return nil, err
	}
	
	return articles, nil
}

func (r *articleRepository) CountByAuthorID(authorID uint) (int64, error) {
	filters := map[string]interface{}{
		"author_id": authorID,
		"status":    "published", // Only count published articles
	}
	
	return r.BaseRepository.Count(&models.Article{}, filters)
}

func (r *articleRepository) IncrementViewCount(id uint) error {
	return r.GetDB().Increment(&models.Article{}, id, "view_count", 1)
}

func (r *articleRepository) UpdateStatistics(id uint, viewCount, likeCount, commentCount uint) error {
	updates := map[string]interface{}{
		"view_count":    viewCount,
		"like_count":    likeCount,
		"comment_count": commentCount,
	}
	return r.GetDB().UpdateFields(&models.Article{}, id, updates)
}

func (r *articleRepository) GetStatistics(id uint) (viewCount, likeCount, commentCount uint, err error) {
	var article models.Article
	err = r.GetDB().GetByID(&article, id)
	if err != nil {
		return 0, 0, 0, err
	}
	return article.ViewCount, article.LikeCount, article.CommentCount, nil
}