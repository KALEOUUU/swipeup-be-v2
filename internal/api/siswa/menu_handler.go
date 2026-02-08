package siswa

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuHandler handles menu/product-related requests for students
type MenuHandler struct {
	db *gorm.DB
}

// NewMenuHandler creates a new MenuHandler instance
func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

// ProductWithStand represents a product with stand information
type ProductWithStand struct {
	models.Product
	StandName string `json:"stand_name"`
}

// GetProducts returns all active products from all active stands
func (h *MenuHandler) GetProducts(c *gin.Context) {
	var products []models.Product
	
	// Get all active products from active stands
	if err := h.db.
		Joins("JOIN stand_settings ON products.stand_id = stand_settings.stand_id AND stand_settings.is_active = ? AND stand_settings.deleted_at IS NULL", true).
		Where("products.is_active = ?", true).
		Preload("Category", "is_active = ?", true).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Build response with stand information
	type ProductResponse struct {
		models.Product
		StandName string `json:"stand_name"`
	}

	var response []ProductResponse
	for _, p := range products {
		// Get stand name
		var standSettings models.StandSettings
		h.db.Where("stand_id = ?", p.StandID).First(&standSettings)
		
		response = append(response, ProductResponse{
			Product:   p,
			StandName: standSettings.StoreName,
		})
	}

	c.JSON(http.StatusOK, response)
}