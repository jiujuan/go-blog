package handlers

import (
	"net/http"
	"strconv"

	"go-blog/internal/services"
	"go-blog/internal/utils"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *services.TagService
}

// NewTagHandler creates a new tag handler
func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// List handles tag listing
// GET /api/tags
func (h *TagHandler) List(c *gin.Context) {
	// Check for popular tags request
	if c.Query("popular") == "true" {
		h.getPopularTags(c)
		return
	}

	tags, err := h.tagService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve tags"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Tags retrieved successfully", tags))
}

// Create handles tag creation
// POST /api/tags
func (h *TagHandler) Create(c *gin.Context) {
	var req services.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request data"))
		return
	}

	tag, err := h.tagService.Create(&req)
	if err != nil {
		if err.Error() == "tag name is required" || 
		   err.Error() == "tag name must be less than 50 characters" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid tag data"))
			return
		}
		if err.Error() == "tag with name '"+req.Name+"' already exists" {
			c.JSON(http.StatusConflict, utils.ErrorResponse("Tag already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create tag"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Tag created successfully", tag))
}

// GetBySlug handles getting tag by slug
// GET /api/tags/:slug
func (h *TagHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Tag slug is required"))
		return
	}

	tag, err := h.tagService.GetBySlug(slug)
	if err != nil {
		if err.Error() == "slug cannot be empty" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid slug"))
			return
		}
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Tag not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Tag retrieved successfully", tag))
}

// GetTagArticles handles getting articles by tag
// GET /api/tags/:slug/articles
func (h *TagHandler) GetTagArticles(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Tag slug is required"))
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := h.tagService.GetArticlesBySlug(slug, page, limit)
	if err != nil {
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Tag not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve articles"))
		return
	}

	// Calculate pagination info
	totalPages := (int(total) + limit - 1) / limit
	
	response := gin.H{
		"articles": articles,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Articles retrieved successfully", response))
}

// getPopularTags handles getting popular tags with usage statistics
// GET /api/tags?popular=true&limit=20
func (h *TagHandler) getPopularTags(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	tags, err := h.tagService.GetPopularTags(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve popular tags"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Popular tags retrieved successfully", tags))
}