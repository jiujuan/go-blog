package repositories

import (
	"fmt"
	"testing"
	"time"

	"go-blog/internal/database"
	"go-blog/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticleRepository_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	db, err := database.SetupTestDB()
	require.NoError(t, err)
	defer database.CleanupTestDB(db)

	// Create repositories
	articleRepo := NewArticleRepository(db)
	userRepo := NewUserRepository(db)
	categoryRepo := NewCategoryRepository(db)
	tagRepo := NewTagRepository(db)

	// Create test user
	user := &models.User{
		Username: "testauthor",
		Email:    "author@example.com",
		Password: "hashedpassword",
	}
	err = userRepo.Create(user)
	require.NoError(t, err)

	// Create test category
	category := &models.Category{
		Name: "Technology",
		Slug: "technology",
	}
	err = categoryRepo.Create(category)
	require.NoError(t, err)

	// Create test tags
	tag1 := &models.Tag{Name: "golang", Slug: "golang"}
	tag2 := &models.Tag{Name: "web", Slug: "web"}
	err = tagRepo.Create(tag1)
	require.NoError(t, err)
	err = tagRepo.Create(tag2)
	require.NoError(t, err)

	t.Run("Create and retrieve article", func(t *testing.T) {
		now := time.Now()
		article := &models.Article{
			Title:       "Test Article",
			Slug:        "test-article",
			Content:     "This is test content for the article",
			Excerpt:     "Test excerpt",
			AuthorID:    user.ID,
			CategoryID:  &category.ID,
			Status:      models.StatusPublished,
			PublishedAt: &now,
			Tags:        []models.Tag{*tag1, *tag2},
		}

		// Create article
		err := articleRepo.Create(article)
		assert.NoError(t, err)
		assert.NotZero(t, article.ID)
		assert.NotZero(t, article.CreatedAt)

		// Retrieve by ID
		retrieved, err := articleRepo.GetByID(article.ID)
		assert.NoError(t, err)
		assert.Equal(t, article.Title, retrieved.Title)
		assert.Equal(t, article.Slug, retrieved.Slug)
		assert.Equal(t, article.Content, retrieved.Content)
		assert.Equal(t, article.Status, retrieved.Status)

		// Verify preloaded relationships
		assert.NotZero(t, retrieved.Author.ID)
		assert.Equal(t, user.Username, retrieved.Author.Username)
		assert.NotNil(t, retrieved.Category)
		assert.Equal(t, category.Name, retrieved.Category.Name)
		assert.Len(t, retrieved.Tags, 2)

		// Retrieve by slug
		retrievedBySlug, err := articleRepo.GetBySlug(article.Slug)
		assert.NoError(t, err)
		assert.Equal(t, article.ID, retrievedBySlug.ID)
		assert.Equal(t, user.Username, retrievedBySlug.Author.Username)
	})

	t.Run("Update article", func(t *testing.T) {
		article := &models.Article{
			Title:    "Update Test Article",
			Slug:     "update-test-article",
			Content:  "Original content",
			AuthorID: user.ID,
			Status:   models.StatusDraft,
		}

		// Create article
		err := articleRepo.Create(article)
		require.NoError(t, err)

		// Update article
		article.Title = "Updated Article Title"
		article.Content = "Updated content"
		article.Status = models.StatusPublished
		now := time.Now()
		article.PublishedAt = &now

		err = articleRepo.Update(article)
		assert.NoError(t, err)

		// Retrieve and verify update
		updated, err := articleRepo.GetByID(article.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Article Title", updated.Title)
		assert.Equal(t, "Updated content", updated.Content)
		assert.Equal(t, models.StatusPublished, updated.Status)
		assert.NotNil(t, updated.PublishedAt)
	})

	t.Run("Delete article", func(t *testing.T) {
		article := &models.Article{
			Title:    "Delete Test Article",
			Slug:     "delete-test-article",
			Content:  "Content to be deleted",
			AuthorID: user.ID,
			Status:   models.StatusDraft,
		}

		// Create article
		err := articleRepo.Create(article)
		require.NoError(t, err)

		// Delete article
		err = articleRepo.Delete(article.ID)
		assert.NoError(t, err)

		// Verify article is deleted
		_, err = articleRepo.GetByID(article.ID)
		assert.Error(t, err)
		assert.True(t, database.IsRecordNotFound(err))
	})

	t.Run("List articles with pagination and filters", func(t *testing.T) {
		// Create multiple articles
		articles := []*models.Article{
			{
				Title:    "Published Article 1",
				Slug:     "published-article-1",
				Content:  "Content 1",
				AuthorID: user.ID,
				Status:   models.StatusPublished,
			},
			{
				Title:    "Published Article 2",
				Slug:     "published-article-2",
				Content:  "Content 2",
				AuthorID: user.ID,
				Status:   models.StatusPublished,
			},
			{
				Title:    "Draft Article",
				Slug:     "draft-article",
				Content:  "Draft content",
				AuthorID: user.ID,
				Status:   models.StatusDraft,
			},
		}

		for _, article := range articles {
			err := articleRepo.Create(article)
			require.NoError(t, err)
		}

		// Test list all articles
		allArticles, total, err := articleRepo.List(0, 10, map[string]interface{}{})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(allArticles), 3)
		assert.GreaterOrEqual(t, total, int64(3))

		// Test filter by status
		publishedArticles, publishedTotal, err := articleRepo.List(0, 10, map[string]interface{}{
			"status": models.StatusPublished,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(publishedArticles), 2)
		assert.GreaterOrEqual(t, publishedTotal, int64(2))

		// Test pagination
		page1, total1, err := articleRepo.List(0, 1, map[string]interface{}{})
		assert.NoError(t, err)
		assert.Len(t, page1, 1)
		assert.GreaterOrEqual(t, total1, int64(3))

		page2, total2, err := articleRepo.List(1, 1, map[string]interface{}{})
		assert.NoError(t, err)
		assert.Len(t, page2, 1)
		assert.Equal(t, total1, total2) // Total should be same

		// Verify different articles on different pages
		assert.NotEqual(t, page1[0].ID, page2[0].ID)
	})

	t.Run("Search articles", func(t *testing.T) {
		// Create articles with searchable content
		searchArticles := []*models.Article{
			{
				Title:    "Golang Tutorial",
				Slug:     "golang-tutorial",
				Content:  "Learn Go programming language basics",
				Excerpt:  "Go tutorial for beginners",
				AuthorID: user.ID,
				Status:   models.StatusPublished,
			},
			{
				Title:    "Web Development",
				Slug:     "web-development",
				Content:  "Building web applications with modern frameworks",
				Excerpt:  "Web dev guide",
				AuthorID: user.ID,
				Status:   models.StatusPublished,
			},
			{
				Title:    "Database Design",
				Slug:     "database-design",
				Content:  "Designing efficient database schemas",
				AuthorID: user.ID,
				Status:   models.StatusPublished,
			},
		}

		for _, article := range searchArticles {
			err := articleRepo.Create(article)
			require.NoError(t, err)
		}

		// Search by title
		results1, total1, err := articleRepo.Search("Golang", 0, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results1), 1)
		assert.GreaterOrEqual(t, total1, int64(1))

		// Search by content
		results2, total2, err := articleRepo.Search("programming", 0, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results2), 1)
		assert.GreaterOrEqual(t, total2, int64(1))

		// Search with no results
		results3, total3, err := articleRepo.Search("nonexistent", 0, 10)
		assert.NoError(t, err)
		assert.Len(t, results3, 0)
		assert.Equal(t, int64(0), total3)
	})

	t.Run("Get articles by author", func(t *testing.T) {
		// Create another user
		author2 := &models.User{
			Username: "author2",
			Email:    "author2@example.com",
			Password: "hashedpassword",
		}
		err := userRepo.Create(author2)
		require.NoError(t, err)

		// Create articles for both authors
		article1 := &models.Article{
			Title:    "Author 1 Article",
			Slug:     "author-1-article",
			Content:  "Content by author 1",
			AuthorID: user.ID,
			Status:   models.StatusPublished,
		}
		article2 := &models.Article{
			Title:    "Author 2 Article",
			Slug:     "author-2-article",
			Content:  "Content by author 2",
			AuthorID: author2.ID,
			Status:   models.StatusPublished,
		}

		err = articleRepo.Create(article1)
		require.NoError(t, err)
		err = articleRepo.Create(article2)
		require.NoError(t, err)

		// Get articles by author 1
		author1Articles, err := articleRepo.GetByAuthorID(user.ID, 10, 0)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(author1Articles), 1)

		// Verify all articles belong to author 1
		for _, article := range author1Articles {
			assert.Equal(t, user.ID, article.AuthorID)
		}

		// Count articles by author 1
		count1, err := articleRepo.CountByAuthorID(user.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count1, int64(1))

		// Count articles by author 2
		count2, err := articleRepo.CountByAuthorID(author2.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count2, int64(1))
	})

	t.Run("Increment view count", func(t *testing.T) {
		article := &models.Article{
			Title:     "View Count Test",
			Slug:      "view-count-test",
			Content:   "Testing view count increment",
			AuthorID:  user.ID,
			Status:    models.StatusPublished,
			ViewCount: 0,
		}

		// Create article
		err := articleRepo.Create(article)
		require.NoError(t, err)

		// Increment view count
		err = articleRepo.IncrementViewCount(article.ID)
		assert.NoError(t, err)

		// Verify view count increased
		updated, err := articleRepo.GetByID(article.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), updated.ViewCount)

		// Increment again
		err = articleRepo.IncrementViewCount(article.ID)
		assert.NoError(t, err)

		updated2, err := articleRepo.GetByID(article.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(2), updated2.ViewCount)
	})

	t.Run("Article with tags many-to-many", func(t *testing.T) {
		// Create additional tags
		tag3 := &models.Tag{Name: "testing", Slug: "testing"}
		err := tagRepo.Create(tag3)
		require.NoError(t, err)

		article := &models.Article{
			Title:    "Tagged Article",
			Slug:     "tagged-article",
			Content:  "Article with multiple tags",
			AuthorID: user.ID,
			Status:   models.StatusPublished,
			Tags:     []models.Tag{*tag1, *tag2, *tag3},
		}

		// Create article with tags
		err = articleRepo.Create(article)
		assert.NoError(t, err)

		// Retrieve and verify tags
		retrieved, err := articleRepo.GetByID(article.ID)
		assert.NoError(t, err)
		assert.Len(t, retrieved.Tags, 3)

		// Verify tag names
		tagNames := make([]string, len(retrieved.Tags))
		for i, tag := range retrieved.Tags {
			tagNames[i] = tag.Name
		}
		assert.Contains(t, tagNames, "golang")
		assert.Contains(t, tagNames, "web")
		assert.Contains(t, tagNames, "testing")
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test GetByID with non-existent ID
		_, err := articleRepo.GetByID(99999)
		assert.Error(t, err)
		assert.True(t, database.IsRecordNotFound(err))

		// Test GetBySlug with non-existent slug
		_, err = articleRepo.GetBySlug("non-existent-slug")
		assert.Error(t, err)
		assert.True(t, database.IsRecordNotFound(err))

		// Test IncrementViewCount with non-existent ID
		err = articleRepo.IncrementViewCount(99999)
		assert.Error(t, err)

		// Test GetByAuthorID with non-existent author
		articles, err := articleRepo.GetByAuthorID(99999, 10, 0)
		assert.NoError(t, err) // Should not error, just return empty
		assert.Len(t, articles, 0)

		// Test CountByAuthorID with non-existent author
		count, err := articleRepo.CountByAuthorID(99999)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestArticleRepository_StatusFiltering(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	db, err := database.SetupTestDB()
	require.NoError(t, err)
	defer database.CleanupTestDB(db)

	articleRepo := NewArticleRepository(db)
	userRepo := NewUserRepository(db)

	// Create test user
	user := &models.User{
		Username: "statustest",
		Email:    "statustest@example.com",
		Password: "hashedpassword",
	}
	err = userRepo.Create(user)
	require.NoError(t, err)

	// Create articles with different statuses
	statuses := []models.ArticleStatus{
		models.StatusDraft,
		models.StatusPublished,
		models.StatusArchived,
	}

	for i, status := range statuses {
		article := &models.Article{
			Title:    fmt.Sprintf("Article %d", i+1),
			Slug:     fmt.Sprintf("article-%d", i+1),
			Content:  fmt.Sprintf("Content %d", i+1),
			AuthorID: user.ID,
			Status:   status,
		}
		err := articleRepo.Create(article)
		require.NoError(t, err)
	}

	// Test filtering by each status
	for _, status := range statuses {
		articles, total, err := articleRepo.List(0, 10, map[string]interface{}{
			"status": status,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(articles), 1)
		assert.GreaterOrEqual(t, total, int64(1))

		// Verify all returned articles have the correct status
		for _, article := range articles {
			assert.Equal(t, status, article.Status)
		}
	}
}