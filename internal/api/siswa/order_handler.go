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

// GetOrdersByMonth returns orders for the current student filtered by month
func (h *OrderHandler) GetOrdersByMonth(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	year := c.Query("year")
	month := c.Query("month")

	if year == "" || month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Year and month parameters are required"})
		return
	}

	// Build date range for the month
	startDate := fmt.Sprintf("%s-%s-01", year, month)
	endDate := fmt.Sprintf("%s-%s-31", year, month) // Simplified, assumes 31 days

	var orders []models.Order
	query := h.db.Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate)
	if err := query.Preload("OrderItems.Product").Preload("Stand").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Calculate monthly summary
	var totalOrders int64
	var totalAmount float64
	for _, order := range orders {
		totalOrders++
		totalAmount += order.TotalAmount
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"summary": gin.H{
			"year": year,
			"month": month,
			"total_orders": totalOrders,
			"total_amount": totalAmount,
		},
	})
}

// GetOrderReceipt generates a printable receipt for an order
func (h *OrderHandler) GetOrderReceipt(c *gin.Context) {
	orderID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", orderID, userID).Preload("OrderItems.Product").Preload("Stand").First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Generate HTML receipt
	receiptHTML := h.generateReceiptHTML(order)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, receiptHTML)
}

// generateReceiptHTML creates HTML receipt for printing
func (h *OrderHandler) generateReceiptHTML(order models.Order) string {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Receipt - %s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 400px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; border-bottom: 2px solid #000; padding-bottom: 10px; margin-bottom: 20px; }
        .order-info { margin-bottom: 20px; }
        .items { margin-bottom: 20px; }
        .item { display: flex; justify-content: space-between; margin-bottom: 5px; }
        .total { border-top: 1px solid #000; padding-top: 10px; font-weight: bold; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
        @media print { body { margin: 0; } }
    </style>
</head>
<body>
    <div class="header">
        <h2>%s</h2>
        <p>Order Receipt</p>
    </div>
    
    <div class="order-info">
        <p><strong>Order Number:</strong> %s</p>
        <p><strong>Date:</strong> %s</p>
        <p><strong>Payment Method:</strong> %s</p>
        <p><strong>Status:</strong> %s</p>
    </div>
    
    <div class="items">
        <h3>Items:</h3>`,
		order.OrderNumber,
		order.Stand.Name,
		order.OrderNumber,
		order.CreatedAt.Format("2006-01-02 15:04:05"),
		order.PaymentMethod,
		order.Status)

	for _, item := range order.OrderItems {
		html += fmt.Sprintf(`
        <div class="item">
            <span>%s (x%d)</span>
            <span>Rp %.0f</span>
        </div>`, item.Product.Name, item.Quantity, item.Subtotal)
	}

	html += fmt.Sprintf(`
    </div>
    
    <div class="total">
        <div class="item">
            <span>Total Amount:</span>
            <span>Rp %.0f</span>
        </div>
    </div>
    
    <div class="footer">
        <p>Thank you for your order!</p>
        <p>Generated on %s</p>
    </div>
</body>
</html>`, order.TotalAmount, time.Now().Format("2006-01-02 15:04:05"))

	return html
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