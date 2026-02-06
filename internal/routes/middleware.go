package routes

import (
	"net/http"
	"strings"

	"swipeup-admin-v2/internal/app/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the authentication token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token and get user information
		userInfo, isValid := auth.ValidateToken(token)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", userInfo.UserID)
		c.Set("user_role", userInfo.Role)
		c.Set("username", userInfo.Username)

		c.Next()
	}
}

// AdminMiddleware validates that the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get user role from context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// Check if user has admin role
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StandMiddleware validates that the user has stand role
func StandMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get user role from context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// Check if user has stand role (or admin role for superadmin access)
		if userRole != "stand_admin" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Stand role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SiswaMiddleware validates that the user has student role
func SiswaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get user role from context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// Check if user has student role (or admin role for superadmin access)
		if userRole != "student" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Student role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
