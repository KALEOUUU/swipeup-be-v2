package auth

import (
	"net/http"
	"swipeup-admin-v2/internal/app/auth"
	"swipeup-admin-v2/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User      models.User `json:"user"`
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Check if user exists by name or email
	if err := h.db.Where("(name = ? OR email = ?) AND is_active = ?", req.Username, req.Username, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token with user information
	token := auth.GenerateToken(user.ID, user.Name, user.Role)
	expiresAt := time.Now().Add(10 * 24 * time.Hour)

	response := LoginResponse{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	// Extract token from Authorization header
	token := authHeader[7:] // Remove "Bearer " prefix
	auth.RemoveToken(token)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	// Extract token from Authorization header
	token := authHeader[7:] // Remove "Bearer " prefix
	
	// Validate token
	userInfo, isValid := auth.ValidateToken(token)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Generate new token
	newToken := auth.GenerateToken(userInfo.UserID, userInfo.Username, userInfo.Role)
	
	// Get new token expiry
	newUserInfo, _ := auth.ValidateToken(newToken)

	// Remove old token
	auth.RemoveToken(token)

	c.JSON(http.StatusOK, gin.H{
		"token":     newToken,
		"expiresAt": newUserInfo.Expiry,
	})
}
