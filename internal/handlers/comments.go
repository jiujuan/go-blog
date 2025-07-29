package handlers

import (
	"net/http"

	"go-blog/internal/services"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// GetByArticle handles getting comments by article
func (h *CommentHandler) GetByArticle(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetByArticle endpoint not implemented yet"})
}

// Create handles comment creation
func (h *CommentHandler) Create(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Create endpoint not implemented yet"})
}

// Update handles comment updates
func (h *CommentHandler) Update(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Update endpoint not implemented yet"})
}

// Delete handles comment deletion
func (h *CommentHandler) Delete(c *gin.Context) {
	// Implementation placeholder - will be implemented in later tasks
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Delete endpoint not implemented yet"})
}