package handlers

import (
	"net/http"

	"go-blog/internal/services"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// List handles category listing
func (h *CategoryHandler) List(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "List endpoint not implemented yet"})
}

// Create handles category creation
func (h *CategoryHandler) Create(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Create endpoint not implemented yet"})
}

// GetBySlug handles getting category by slug
func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetBySlug endpoint not implemented yet"})
}

// Update handles category updates
func (h *CategoryHandler) Update(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Update endpoint not implemented yet"})
}

// Delete handles category deletion
func (h *CategoryHandler) Delete(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Delete endpoint not implemented yet"})
}

// GetCategoryArticles handles getting articles by category
func (h *CategoryHandler) GetCategoryArticles(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetCategoryArticles endpoint not implemented yet"})
}