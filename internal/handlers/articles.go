package handlers

import (
	"net/http"

	"go-blog/internal/services"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *services.ArticleService
}

// NewArticleHandler creates a new article handler
func NewArticleHandler(articleService *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// List handles article listing
func (h *ArticleHandler) List(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "List endpoint not implemented yet"})
}

// Create handles article creation
func (h *ArticleHandler) Create(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Create endpoint not implemented yet"})
}

// GetBySlug handles getting article by slug
func (h *ArticleHandler) GetBySlug(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetBySlug endpoint not implemented yet"})
}

// Update handles article updates
func (h *ArticleHandler) Update(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Update endpoint not implemented yet"})
}

// Delete handles article deletion
func (h *ArticleHandler) Delete(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Delete endpoint not implemented yet"})
}

// Search handles article search
func (h *ArticleHandler) Search(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Search endpoint not implemented yet"})
}

// ToggleLike handles article like/unlike
func (h *ArticleHandler) ToggleLike(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "ToggleLike endpoint not implemented yet"})
}

// GetArchive handles archive listing
func (h *ArticleHandler) GetArchive(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetArchive endpoint not implemented yet"})
}

// GetArchiveByMonth handles archive by month
func (h *ArticleHandler) GetArchiveByMonth(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetArchiveByMonth endpoint not implemented yet"})
}