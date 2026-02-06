package stand

import (
	"net/http"
	"swipeup-admin-v2/internal/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductHandler handles product-related requests for stand admins
type ProductHandler struct {
	db *gorm.DB
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// GetProducts returns all products for the current stand
func (h *ProductHandler) GetProducts(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var products []models.Product
	if err := h.db.Where("stand_id = ?", standID).Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProduct returns a single product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var product models.Product
	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).Preload("Category").First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Calculate discounted price
	if product.Discount > 0 {
		product.DiscountedPrice = product.Price - (product.Price * product.Discount / 100)
	} else {
		product.DiscountedPrice = product.Price
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		CategoryID  uint    `json:"category_id" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Stock       int     `json:"stock" binding:"required"`
		ImageBase64 string  `json:"image_base64"`
		QRISBase64  string  `json:"qris_base64"`
		Discount    float64 `json:"discount" binding:"min=0,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageBase64: req.ImageBase64,
		QRISBase64:  req.QRISBase64,
		Discount:     req.Discount,
		IsActive:    true,
		StandID:      standID.(uint),
		Status:       "request", // Default status for new products
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var product models.Product
	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CategoryID  uint    `json:"category_id"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		ImageBase64 string  `json:"image_base64"`
		QRISBase64  string  `json:"qris_base64"`
		Discount    float64 `json:"discount" binding:"min=0,max=100"`
		IsActive    *bool   `json:"is_active"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.CategoryID != 0 {
		product.CategoryID = req.CategoryID
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.ImageBase64 != "" {
		product.ImageBase64 = req.ImageBase64
	}
	if req.QRISBase64 != "" {
		product.QRISBase64 = req.QRISBase64
	}
	if req.Discount >= 0 {
		product.Discount = req.Discount
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.Status != nil {
		product.Status = *req.Status
	}

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).Delete(&models.Product{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// UpdateProductStatus updates product status
func (h *ProductHandler) UpdateProductStatus(c *gin.Context) {
	id := c.Param("id")
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"payment_pending": true,
		"request":       true,
		"cooking":       true,
		"done":          true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var product models.Product
	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Status = req.Status
	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product status updated successfully", "status": product.Status})
}
