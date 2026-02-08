package stand

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CategoryHandler handles category-related requests for stand admins
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategories returns all active categories for stand admins
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := h.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}