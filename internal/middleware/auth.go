package middleware

import (
	"net/http"
	"strings"

	"go-blog/internal/services"
	"go-blog/internal/utils"

	"github.com/gin-gonic/gin"
)

// Auth middleware for JWT authentication
func Auth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Authorization header required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid authorization header format"))
			c.Abort()
			return
		}

		token := tokenParts[1]
		user, err := authService.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user information in context for use in handlers
		c.Set("userID", user.ID)
		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuth middleware for optional JWT authentication
func OptionalAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		user, err := authService.GetUserFromToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context if valid
		c.Set("userID", user.ID)
		c.Set("user", user)
		c.Next()
	}
}