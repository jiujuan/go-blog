package services

import (
	"fmt"
	"time"

	"go-blog/internal/models"
	"go-blog/internal/repositories"
)

// ArchiveEntry represents a single archive entry with date and count
type ArchiveEntry struct {
	Year         int   `json:"year"`
	Month        int   `json:"month"`
	MonthName    string `json:"month_name"`
	ArticleCount int64 `json:"article_count"`
}

// ArchiveYear represents a year with its months
type ArchiveYear struct {
	Year   int            `json:"year"`
	Months []ArchiveEntry `json:"months"`
	Total  int64          `json:"total"`
}

// ArchiveResponse represents the complete archive structure
type ArchiveResponse struct {
	Years []ArchiveYear `json:"years"`
	Total int64         `json:"total"`
}

// ArchiveService handles archive-related operations
type ArchiveService struct {
	articleRepo repositories.ArticleRepository
}

// NewArchiveService creates a new archive service
func NewArchiveService(articleRepo repositories.ArticleRepository) *ArchiveService {
	return &ArchiveService{
		articleRepo: articleRepo,
	}
}

// GetArchive returns the complete archive structure with year/month organization
func (s *ArchiveService) GetArchive() (*ArchiveResponse, error) {
	// Get all published articles grouped by year and month
	archiveData, err := s.articleRepo.GetArchive()
	if err != nil {
		return nil, fmt.Errorf("failed to get archive data: %w", err)
	}

	// Convert the raw data to structured response
	response := &ArchiveResponse{
		Years: []ArchiveYear{},
		Total: 0,
	}

	// Process the archive data
	if yearData, ok := archiveData["years"].([]map[string]interface{}); ok {
		for _, year := range yearData {
			archiveYear := ArchiveYear{
				Year:   int(year["year"].(int64)),
				Months: []ArchiveEntry{},
				Total:  0,
			}

			if monthData, ok := year["months"].([]map[string]interface{}); ok {
				for _, month := range monthData {
					entry := ArchiveEntry{
						Year:         archiveYear.Year,
						Month:        int(month["month"].(int64)),
						MonthName:    s.getMonthName(int(month["month"].(int64))),
						ArticleCount: month["count"].(int64),
					}
					archiveYear.Months = append(archiveYear.Months, entry)
					archiveYear.Total += entry.ArticleCount
				}
			}

			response.Years = append(response.Years, archiveYear)
			response.Total += archiveYear.Total
		}
	}

	return response, nil
}

// GetArchiveByMonth returns articles for a specific year and month
func (s *ArchiveService) GetArchiveByMonth(year, month, page, limit int) ([]models.Article, int64, error) {
	if year < 1900 || year > time.Now().Year()+1 {
		return nil, 0, fmt.Errorf("invalid year: %d", year)
	}

	if month < 1 || month > 12 {
		return nil, 0, fmt.Errorf("invalid month: %d", month)
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	articles, total, err := s.articleRepo.GetByMonth(year, month, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get articles for %d/%d: %w", year, month, err)
	}

	return articles, total, nil
}

// GetArchiveStatistics returns statistics about the archive
func (s *ArchiveService) GetArchiveStatistics() (map[string]interface{}, error) {
	archive, err := s.GetArchive()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_articles": archive.Total,
		"total_years":    len(archive.Years),
		"total_months":   0,
	}

	totalMonths := 0
	for _, year := range archive.Years {
		totalMonths += len(year.Months)
	}
	stats["total_months"] = totalMonths

	// Find most active year and month
	var mostActiveYear *ArchiveYear
	var mostActiveMonth *ArchiveEntry

	for i, year := range archive.Years {
		if mostActiveYear == nil || year.Total > mostActiveYear.Total {
			mostActiveYear = &archive.Years[i]
		}

		for j, month := range year.Months {
			if mostActiveMonth == nil || month.ArticleCount > mostActiveMonth.ArticleCount {
				mostActiveMonth = &year.Months[j]
			}
		}
	}

	if mostActiveYear != nil {
		stats["most_active_year"] = map[string]interface{}{
			"year":  mostActiveYear.Year,
			"count": mostActiveYear.Total,
		}
	}

	if mostActiveMonth != nil {
		stats["most_active_month"] = map[string]interface{}{
			"year":       mostActiveMonth.Year,
			"month":      mostActiveMonth.Month,
			"month_name": mostActiveMonth.MonthName,
			"count":      mostActiveMonth.ArticleCount,
		}
	}

	return stats, nil
}

// getMonthName returns the month name for a given month number
func (s *ArchiveService) getMonthName(month int) string {
	months := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	if month < 1 || month > 12 {
		return "Unknown"
	}

	return months[month]
}

// ValidateArchiveRequest validates archive request parameters
func (s *ArchiveService) ValidateArchiveRequest(year, month int) error {
	currentYear := time.Now().Year()

	if year < 1900 || year > currentYear+1 {
		return fmt.Errorf("year must be between 1900 and %d", currentYear+1)
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}

	return nil
}