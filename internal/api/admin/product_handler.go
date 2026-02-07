package admin

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductHandler handles product-related requests for admin
type ProductHandler struct {
	db *gorm.DB
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}


// CreateStandCanteen creates a new stand canteen
func (h *ProductHandler) CreateStandCanteen(c *gin.Context) {
	var req struct {
		StandID   uint   `json:"stand_id" binding:"required"`
		StoreName string `json:"store_name" binding:"required"`
		QRIS      string `json:"qris"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists and has stand_admin role
	var user models.User
	if err := h.db.First(&user, req.StandID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	if user.Role != "stand_admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User must have stand_admin role"})
		return
	}

	// Check if stand settings already exist for this user
	var existing models.StandSettings
	if err := h.db.Where("stand_id = ?", req.StandID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stand canteen already exists for this user"})
		return
	}

	// Create stand canteen settings
	settings := models.StandSettings{
		StandID:   req.StandID,
		StoreName: req.StoreName,
		QRIS:       req.QRIS,
		IsActive:   true,
	}

	if err := h.db.Create(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stand canteen: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, settings)
}

// UpdateStandCanteen updates an existing stand canteen
func (h *ProductHandler) UpdateStandCanteen(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		StoreName string `json:"store_name"`
		QRIS      string `json:"qris"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var settings models.StandSettings
	if err := h.db.First(&settings, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stand canteen not found"})
		return
	}

	if req.StoreName != "" {
		settings.StoreName = req.StoreName
	}
	if req.QRIS != "" {
		settings.QRIS = req.QRIS
	}
	if req.IsActive != nil {
		settings.IsActive = *req.IsActive
	}

	if err := h.db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stand canteen"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// DeleteStandCanteen deletes a stand canteen by ID
func (h *ProductHandler) DeleteStandCanteen(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.StandSettings{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stand canteen"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stand canteen deleted successfully"})
}

// GetStandCanteens returns all stand canteens
func (h *ProductHandler) GetStandCanteens(c *gin.Context) {
	var standCanteens []models.StandSettings
	if err := h.db.Preload("Stand").Find(&standCanteens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stand canteens"})
		return
	}
	c.JSON(http.StatusOK, standCanteens)
}

// GetStandCanteen returns a single stand canteen by ID
func (h *ProductHandler) GetStandCanteen(c *gin.Context) {
	id := c.Param("id")
	var standCanteen models.StandSettings
	if err := h.db.Preload("Stand").First(&standCanteen, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stand canteen not found"})
		return
	}
	c.JSON(http.StatusOK, standCanteen)
}
