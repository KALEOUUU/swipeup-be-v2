package admin

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler handles user-related requests for admin
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name      string  `json:"name" binding:"required"`
	Email     string  `json:"email" binding:"required"`
	Phone     string  `json:"phone"`
	Role      string  `json:"role" binding:"required,oneof=student admin stand_admin"`
	Class     string  `json:"class"`
	Balance   float64 `json:"balance"`
	IsActive  bool    `json:"is_active"`
	RFIDCard  string  `json:"rfid_card"`
	StudentId string  `json:"student_id"`
	Password  string  `json:"password" binding:"required"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Role      string  `json:"role" binding:"omitempty,oneof=student admin stand_admin"`
	Class     string  `json:"class"`
	Balance   float64 `json:"balance"`
	IsActive  bool    `json:"is_active"`
	RFIDCard  string  `json:"rfid_card"`
	StudentId string  `json:"student_id"`
	Password  string  `json:"password"` // Optional for update
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Validate that all users have valid roles
	for i := range users {
		if users[i].Role == "" {
			// Set default role for users with empty role
			users[i].Role = "student"
			// Update in database
			h.db.Model(&users[i]).Update("role", "student")
		}
	}

	c.JSON(http.StatusOK, users)
}

// GetUser returns a single user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fix empty role if found
	if user.Role == "" {
		user.Role = "student"
		h.db.Model(&user).Update("role", "student")
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user model
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Role:      req.Role, // Role is now validated and required
		Class:     req.Class,
		Balance:   req.Balance,
		IsActive:  true,
		RFIDCard:  req.RFIDCard,
		StudentId: req.StudentId,
		Password:  string(hashedPassword),
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var existingUser models.User
	if err := h.db.First(&existingUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If password is provided, hash it; otherwise keep existing password
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		existingUser.Password = string(hashedPassword)
	}

	// Update other fields
	existingUser.Name = req.Name
	existingUser.Email = req.Email
	existingUser.Phone = req.Phone
	existingUser.Role = req.Role
	existingUser.Class = req.Class
	existingUser.Balance = req.Balance
	existingUser.IsActive = req.IsActive
	existingUser.RFIDCard = req.RFIDCard
	existingUser.StudentId = req.StudentId

	if err := h.db.Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// TopUpBalance adds balance to a user account
func (h *UserHandler) TopUpBalance(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		TransactionNumber: "TOPUP-" + time.Now().Format("20060102150405"),
		UserID:            user.ID,
		Type:              "top_up",
		Amount:            req.Amount,
		BalanceBefore:     user.Balance,
		BalanceAfter:      user.Balance + req.Amount,
		Description:       "Balance top-up",
	}

	// Update user balance
	user.Balance += req.Amount

	// Start transaction
	tx := h.db.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Balance topped up successfully",
		"balance": user.Balance,
	})
}
