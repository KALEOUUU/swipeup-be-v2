package siswa

import (
	"fmt"
	"net/http"
	"swipeup-admin-v2/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderHandler handles order-related requests for students
type OrderHandler struct {
	db *gorm.DB
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// CreateOrder creates a new order for the current student
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
		Items         []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment method
	if req.PaymentMethod != "card" && req.PaymentMethod != "cash" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	// Group items by stand
	standItems := make(map[uint][]models.OrderItem)
	totalByStand := make(map[uint]float64)

	for _, item := range req.Items {
		// Get product details
		var product models.Product
		if err := h.db.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found: " + fmt.Sprintf("%d", item.ProductID)})
			return
		}

		// Check if product is active and has stock
		if !product.IsActive {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not available: " + product.Name})
			return
		}
		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product: " + product.Name})
			return
		}

		// Calculate price (with discount)
		price := product.Price
		if product.Discount > 0 {
			price = product.Price * (1 - product.Discount/100)
		}

		// Create order item
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
			Subtotal:  price * float64(item.Quantity),
		}

		// Group by stand
		standItems[product.StandID] = append(standItems[product.StandID], orderItem)
		totalByStand[product.StandID] += orderItem.Subtotal
	}

	// Create orders for each stand
	var createdOrders []models.Order
	for standID, items := range standItems {
		// Generate order number
		orderNumber := fmt.Sprintf("ORD-%d-%d-%d", userID, standID, time.Now().Unix())

		// Create order
		order := models.Order{
			OrderNumber:   orderNumber,
			UserID:        userID.(uint),
			TotalAmount:   totalByStand[standID],
			Status:        "request",
			PaymentMethod: req.PaymentMethod,
			StandID:       standID,
			OrderItems:    items,
		}

		if err := h.db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		createdOrders = append(createdOrders, order)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Orders created successfully",
		"orders":  createdOrders,
	})
}

// GetOrders returns all orders for the current student
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var orders []models.Order
	if err := h.db.Where("user_id = ?", userID).Preload("OrderItems.Product").Preload("Stand").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// DeleteOrder deletes an order for the current student (soft delete)
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Only allow deletion of orders that are still pending or in request status
	if order.Status != "payment_pending" && order.Status != "request" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only cancel orders that are pending or in request status"})
		return
	}

	// Soft delete the order
	if err := h.db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}