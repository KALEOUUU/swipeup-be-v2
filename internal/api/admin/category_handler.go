package admin

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CategoryHandler handles category-related requests for admin
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategories returns all categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := h.db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategory returns a single category by ID
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.db.Preload("Products").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category by ID
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetGlobalSettings returns all global settings
func (h *CategoryHandler) GetGlobalSettings(c *gin.Context) {
	var settings []models.GlobalSettings
	if err := h.db.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch global settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// UpdateGlobalSetting updates a global setting
func (h *CategoryHandler) UpdateGlobalSetting(c *gin.Context) {
	key := c.Param("key")
	var req struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.GlobalSettings
	if err := h.db.Where("`key` = ?", key).First(&setting).Error; err != nil {
		// Create new setting if not found
		setting = models.GlobalSettings{
			Key:      key,
			Value:     req.Value,
			IsActive:  true,
		}

		if err := h.db.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create setting"})
			return
		}
	} else {
		// Update existing setting
		setting.Value = req.Value
		if err := h.db.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
			return
		}
	}

	c.JSON(http.StatusOK, setting)
}
