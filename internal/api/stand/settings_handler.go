package stand

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingsHandler handles stand settings-related requests
type SettingsHandler struct {
	db *gorm.DB
}

// NewSettingsHandler creates a new SettingsHandler instance
func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// GetSettings returns the current stand's settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var settings models.StandSettings
	if err := h.db.Where("stand_id = ?", standID).First(&settings).Error; err != nil {
		// Return default settings if not found
		settings = models.StandSettings{
			StandID:  standID.(uint),
			StoreName: "My Canteen Stand",
			QRIS:       "",
			IsActive:   true,
		}
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates the stand's settings
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		StoreName string `json:"store_name" binding:"required"`
		QRIS      string `json:"qris" binding:"required"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var settings models.StandSettings
	if err := h.db.Where("stand_id = ?", standID).First(&settings).Error; err != nil {
		// Create new settings if not found
		settings = models.StandSettings{
			StandID:   standID.(uint),
			StoreName: req.StoreName,
			QRIS:       req.QRIS,
			IsActive:   true,
		}

		if req.IsActive != nil {
			settings.IsActive = *req.IsActive
		}

		if err := h.db.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	} else {
		// Update existing settings
		settings.StoreName = req.StoreName
		settings.QRIS = req.QRIS
		if req.IsActive != nil {
			settings.IsActive = *req.IsActive
		}

		if err := h.db.Save(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
			return
		}
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateQRIS updates only the QRIS code
func (h *SettingsHandler) UpdateQRIS(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		QRIS string `json:"qris" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var settings models.StandSettings
	if err := h.db.Where("stand_id = ?", standID).First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settings not found"})
		return
	}

	settings.QRIS = req.QRIS
	if err := h.db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update QRIS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QRIS updated successfully", "qris": settings.QRIS})
}

// UpdateStoreName updates only the store name
func (h *SettingsHandler) UpdateStoreName(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		StoreName string `json:"store_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var settings models.StandSettings
	if err := h.db.Where("stand_id = ?", standID).First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settings not found"})
		return
	}

	settings.StoreName = req.StoreName
	if err := h.db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store name updated successfully", "store_name": settings.StoreName})
}
