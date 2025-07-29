package handlers

import (
	"net/http"
	"strconv"

	"go-blog/internal/models"
	"go-blog/internal/services"
	"go-blog/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetByID handles getting user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID"))
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, utils.SuccessResponse("User retrieved successfully", user))
}

// Update handles user profile updates
func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID"))
		return
	}

	// Check if user is updating their own profile
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	currentUserModel, ok := currentUser.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid user data"))
		return
	}

	if currentUserModel.ID != uint(id) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("You can only update your own profile"))
		return
	}

	var updateReq services.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format"))
		return
	}

	updatedUser, err := h.userService.UpdateProfile(uint(id), &updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	// Remove password from response
	updatedUser.Password = ""

	c.JSON(http.StatusOK, utils.SuccessResponse("Profile updated successfully", updatedUser))
}

// GetUserArticles handles getting articles by user
func (h *UserHandler) GetUserArticles(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID"))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := h.userService.GetUserArticles(uint(id), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user articles"))
		return
	}

	utils.PaginatedSuccessResponse(c, articles, page, limit, total)
}