package services

import (
	"fmt"
	"os"
	"testing"
)

// RunStatisticsTests runs all statistics-related tests
func RunStatisticsTests() {
	fmt.Println("Running Statistics Service Tests...")
	
	// Run unit tests
	fmt.Println("\n=== Running Statistics Service Unit Tests ===")
	if code := testing.Main(func(pat, str string) (bool, error) { return true, nil }, 
		[]testing.InternalTest{
			{"TestStatisticsService_IncrementViewCount", TestStatisticsService_IncrementViewCount},
			{"TestStatisticsService_GetArticleStats", TestStatisticsService_GetArticleStats},
			{"TestStatisticsService_GetAuthorStats", TestStatisticsService_GetAuthorStats},
			{"TestStatisticsService_GetPopularArticles", TestStatisticsService_GetPopularArticles},
			{"TestStatisticsService_GetTrendingArticles", TestStatisticsService_GetTrendingArticles},
			{"TestStatisticsService_GetAuthorSummaryStats", TestStatisticsService_GetAuthorSummaryStats},
			{"TestStatisticsService_GetPeriodStats", TestStatisticsService_GetPeriodStats},
			{"TestStatisticsService_DefaultValues", TestStatisticsService_DefaultValues},
		}, 
		nil, nil); code != 0 {
		fmt.Printf("Unit tests failed with code: %d\n", code)
		os.Exit(code)
	}
	
	// Run integration tests
	fmt.Println("\n=== Running Statistics Service Integration Tests ===")
	if code := testing.Main(func(pat, str string) (bool, error) { return true, nil }, 
		[]testing.InternalTest{
			{"TestStatisticsService_Integration", TestStatisticsService_Integration},
			{"TestStatisticsService_EdgeCases", TestStatisticsService_EdgeCases},
		}, 
		nil, nil); code != 0 {
		fmt.Printf("Integration tests failed with code: %d\n", code)
		os.Exit(code)
	}
	
	fmt.Println("\n✅ All Statistics Service tests passed!")
}

func main() {
	RunStatisticsTests()
}