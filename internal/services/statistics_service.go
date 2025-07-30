package services

import (
	"fmt"
	"go-blog/internal/models"
	"go-blog/internal/repositories"
	"time"
)

// StatisticsService handles article statistics and analytics
type StatisticsService struct {
	articleRepo repositories.ArticleRepository
	likeRepo    repositories.LikeRepository
	commentRepo repositories.CommentRepository
}

// NewStatisticsService creates a new statistics service
func NewStatisticsService(
	articleRepo repositories.ArticleRepository,
	likeRepo repositories.LikeRepository,
	commentRepo repositories.CommentRepository,
) *StatisticsService {
	return &StatisticsService{
		articleRepo: articleRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
	}
}

// ArticleStats represents comprehensive article statistics
type ArticleStats struct {
	ArticleID    uint      `json:"article_id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	ViewCount    uint      `json:"view_count"`
	LikeCount    uint      `json:"like_count"`
	CommentCount uint      `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	PublishedAt  *time.Time `json:"published_at"`
}

// PopularArticle represents popular article data
type PopularArticle struct {
	ArticleID    uint      `json:"article_id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	ViewCount    uint      `json:"view_count"`
	LikeCount    uint      `json:"like_count"`
	CommentCount uint      `json:"comment_count"`
	Score        float64   `json:"score"` // Calculated popularity score
	CreatedAt    time.Time `json:"created_at"`
}

// TrendingArticle represents trending article data
type TrendingArticle struct {
	ArticleID     uint      `json:"article_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	ViewCount     uint      `json:"view_count"`
	LikeCount     uint      `json:"like_count"`
	CommentCount  uint      `json:"comment_count"`
	TrendingScore float64   `json:"trending_score"` // Recent activity score
	CreatedAt     time.Time `json:"created_at"`
}

// AuthorStats represents author statistics
type AuthorStats struct {
	AuthorID     uint `json:"author_id"`
	ArticleCount uint `json:"article_count"`
	TotalViews   uint `json:"total_views"`
	TotalLikes   uint `json:"total_likes"`
	TotalComments uint `json:"total_comments"`
}

// PeriodStats represents statistics for a time period
type PeriodStats struct {
	Period       string `json:"period"`
	ArticleCount uint   `json:"article_count"`
	TotalViews   uint   `json:"total_views"`
	TotalLikes   uint   `json:"total_likes"`
	TotalComments uint  `json:"total_comments"`
}

// IncrementViewCount increments the view count for an article
func (s *StatisticsService) IncrementViewCount(articleID uint) error {
	return s.articleRepo.IncrementViewCount(articleID)
}

// GetArticleStats retrieves comprehensive statistics for a specific article
func (s *StatisticsService) GetArticleStats(articleID uint) (*ArticleStats, error) {
	article, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return &ArticleStats{
		ArticleID:    article.ID,
		Title:        article.Title,
		Slug:         article.Slug,
		ViewCount:    article.ViewCount,
		LikeCount:    article.LikeCount,
		CommentCount: article.CommentCount,
		CreatedAt:    article.CreatedAt,
		PublishedAt:  article.PublishedAt,
	}, nil
}

// GetAuthorStats retrieves statistics for all articles by a specific author
func (s *StatisticsService) GetAuthorStats(authorID uint) ([]*ArticleStats, error) {
	articles, err := s.articleRepo.GetByAuthorID(authorID, 100, 0) // Get up to 100 articles
	if err != nil {
		return nil, fmt.Errorf("failed to get author articles: %w", err)
	}

	stats := make([]*ArticleStats, len(articles))
	for i, article := range articles {
		stats[i] = &ArticleStats{
			ArticleID:    article.ID,
			Title:        article.Title,
			Slug:         article.Slug,
			ViewCount:    article.ViewCount,
			LikeCount:    article.LikeCount,
			CommentCount: article.CommentCount,
			CreatedAt:    article.CreatedAt,
			PublishedAt:  article.PublishedAt,
		}
	}

	return stats, nil
}

// GetPopularArticles retrieves the most popular articles based on engagement metrics
func (s *StatisticsService) GetPopularArticles(limit int) ([]*PopularArticle, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get published articles with statistics
	filters := map[string]interface{}{
		"status": models.StatusPublished,
	}
	
	articles, _, err := s.articleRepo.List(0, limit*2, filters) // Get more to calculate scores
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	popularArticles := make([]*PopularArticle, 0, len(articles))
	for _, article := range articles {
		// Calculate popularity score based on views, likes, and comments
		// Weight: views (1x), likes (3x), comments (5x)
		score := float64(article.ViewCount) + 
				float64(article.LikeCount)*3 + 
				float64(article.CommentCount)*5

		popularArticles = append(popularArticles, &PopularArticle{
			ArticleID:    article.ID,
			Title:        article.Title,
			Slug:         article.Slug,
			ViewCount:    article.ViewCount,
			LikeCount:    article.LikeCount,
			CommentCount: article.CommentCount,
			Score:        score,
			CreatedAt:    article.CreatedAt,
		})
	}

	// Sort by score (descending) and limit results
	for i := 0; i < len(popularArticles)-1; i++ {
		for j := i + 1; j < len(popularArticles); j++ {
			if popularArticles[i].Score < popularArticles[j].Score {
				popularArticles[i], popularArticles[j] = popularArticles[j], popularArticles[i]
			}
		}
	}

	if len(popularArticles) > limit {
		popularArticles = popularArticles[:limit]
	}

	return popularArticles, nil
}

// GetTrendingArticles retrieves trending articles based on recent activity
func (s *StatisticsService) GetTrendingArticles(limit int, days int) ([]*TrendingArticle, error) {
	if limit <= 0 {
		limit = 10
	}
	if days <= 0 {
		days = 7 // Default to last 7 days
	}

	// Get recent published articles
	filters := map[string]interface{}{
		"status": models.StatusPublished,
	}
	
	articles, _, err := s.articleRepo.List(0, limit*2, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	trendingArticles := make([]*TrendingArticle, 0, len(articles))

	for _, article := range articles {
		// Only consider articles published within the trending period or with recent activity
		if article.PublishedAt != nil && article.PublishedAt.After(cutoffDate) {
			// Calculate trending score with time decay
			daysSincePublished := time.Since(*article.PublishedAt).Hours() / 24
			timeDecay := 1.0 / (1.0 + daysSincePublished/float64(days))
			
			// Base score with time decay
			baseScore := (float64(article.ViewCount) + 
						 float64(article.LikeCount)*3 + 
						 float64(article.CommentCount)*5) * timeDecay

			trendingArticles = append(trendingArticles, &TrendingArticle{
				ArticleID:     article.ID,
				Title:         article.Title,
				Slug:          article.Slug,
				ViewCount:     article.ViewCount,
				LikeCount:     article.LikeCount,
				CommentCount:  article.CommentCount,
				TrendingScore: baseScore,
				CreatedAt:     article.CreatedAt,
			})
		}
	}

	// Sort by trending score (descending) and limit results
	for i := 0; i < len(trendingArticles)-1; i++ {
		for j := i + 1; j < len(trendingArticles); j++ {
			if trendingArticles[i].TrendingScore < trendingArticles[j].TrendingScore {
				trendingArticles[i], trendingArticles[j] = trendingArticles[j], trendingArticles[i]
			}
		}
	}

	if len(trendingArticles) > limit {
		trendingArticles = trendingArticles[:limit]
	}

	return trendingArticles, nil
}

// GetAuthorSummaryStats retrieves aggregated statistics for an author
func (s *StatisticsService) GetAuthorSummaryStats(authorID uint) (*AuthorStats, error) {
	articles, err := s.articleRepo.GetByAuthorID(authorID, 1000, 0) // Get all articles
	if err != nil {
		return nil, fmt.Errorf("failed to get author articles: %w", err)
	}

	stats := &AuthorStats{
		AuthorID: authorID,
	}

	for _, article := range articles {
		stats.ArticleCount++
		stats.TotalViews += article.ViewCount
		stats.TotalLikes += article.LikeCount
		stats.TotalComments += article.CommentCount
	}

	return stats, nil
}

// GetPeriodStats retrieves statistics aggregated by time periods
func (s *StatisticsService) GetPeriodStats(period string, limit int) ([]*PeriodStats, error) {
	// This is a simplified implementation
	// In a production system, you would use more sophisticated SQL queries
	// to aggregate data by time periods
	
	filters := map[string]interface{}{
		"status": models.StatusPublished,
	}
	
	articles, _, err := s.articleRepo.List(0, 1000, filters) // Get many articles for aggregation
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	// Group articles by period (simplified monthly grouping)
	periodMap := make(map[string]*PeriodStats)
	
	for _, article := range articles {
		var periodKey string
		if article.PublishedAt != nil {
			switch period {
			case "monthly":
				periodKey = article.PublishedAt.Format("2006-01")
			case "yearly":
				periodKey = article.PublishedAt.Format("2006")
			default:
				periodKey = article.PublishedAt.Format("2006-01-02")
			}
		} else {
			periodKey = "unpublished"
		}

		if _, exists := periodMap[periodKey]; !exists {
			periodMap[periodKey] = &PeriodStats{
				Period: periodKey,
			}
		}

		stats := periodMap[periodKey]
		stats.ArticleCount++
		stats.TotalViews += article.ViewCount
		stats.TotalLikes += article.LikeCount
		stats.TotalComments += article.CommentCount
	}

	// Convert map to slice and sort by period
	result := make([]*PeriodStats, 0, len(periodMap))
	for _, stats := range periodMap {
		result = append(result, stats)
	}

	// Simple sorting by period (descending)
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Period < result[j].Period {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}