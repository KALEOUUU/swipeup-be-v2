package stand

import (
	"fmt"
	"net/http"
	"swipeup-admin-v2/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderHandler handles order-related requests for stand admins
type OrderHandler struct {
	db *gorm.DB
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// GetOrders returns all orders for the current stand
func (h *OrderHandler) GetOrders(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var orders []models.Order
	if err := h.db.Where("stand_id = ?", standID).Preload("User").Preload("OrderItems.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetOrdersByMonth returns orders for the current stand filtered by month
func (h *OrderHandler) GetOrdersByMonth(c *gin.Context) {
	standID, exists := c.Get("user_id")
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
	query := h.db.Where("stand_id = ? AND created_at >= ? AND created_at <= ?", standID, startDate, endDate)
	if err := query.Preload("User").Preload("OrderItems.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Calculate monthly summary
	var totalOrders int64
	var totalRevenue float64
	var completedOrders int64
	var pendingOrders int64

	for _, order := range orders {
		totalOrders++
		totalRevenue += order.TotalAmount
		if order.Status == "done" {
			completedOrders++
		} else if order.Status == "payment_pending" || order.Status == "request" {
			pendingOrders++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"summary": gin.H{
			"year": year,
			"month": month,
			"total_orders": totalOrders,
			"completed_orders": completedOrders,
			"pending_orders": pendingOrders,
			"total_revenue": totalRevenue,
		},
	})
}

// GetMonthlyRevenueRecap returns monthly revenue recap for the current stand
func (h *OrderHandler) GetMonthlyRevenueRecap(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	year := c.Query("year")
	if year == "" {
		year = fmt.Sprintf("%d", time.Now().Year())
	}

	// Get revenue data for each month of the year
	var monthlyData []gin.H

	for month := 1; month <= 12; month++ {
		startDate := fmt.Sprintf("%s-%02d-01", year, month)
		endDate := fmt.Sprintf("%s-%02d-31", year, month)

		var result struct {
			TotalOrders   int64   `json:"total_orders"`
			TotalRevenue  float64 `json:"total_revenue"`
			CompletedOrders int64 `json:"completed_orders"`
		}

		h.db.Model(&models.Order{}).
			Where("stand_id = ? AND created_at >= ? AND created_at <= ?", standID, startDate, endDate).
			Select("COUNT(*) as total_orders, SUM(total_amount) as total_revenue, SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END) as completed_orders").
			Scan(&result)

		monthlyData = append(monthlyData, gin.H{
			"month": month,
			"month_name": time.Month(month).String(),
			"total_orders": result.TotalOrders,
			"completed_orders": result.CompletedOrders,
			"total_revenue": result.TotalRevenue,
		})
	}

	// Calculate yearly totals
	var yearlyTotal struct {
		TotalOrders   int64   `json:"total_orders"`
		TotalRevenue  float64 `json:"total_revenue"`
		CompletedOrders int64 `json:"completed_orders"`
	}

	startYear := fmt.Sprintf("%s-01-01", year)
	endYear := fmt.Sprintf("%s-12-31", year)

	h.db.Model(&models.Order{}).
		Where("stand_id = ? AND created_at >= ? AND created_at <= ?", standID, startYear, endYear).
		Select("COUNT(*) as total_orders, SUM(total_amount) as total_revenue, SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END) as completed_orders").
		Scan(&yearlyTotal)

	c.JSON(http.StatusOK, gin.H{
		"year": year,
		"monthly_data": monthlyData,
		"yearly_summary": gin.H{
			"total_orders": yearlyTotal.TotalOrders,
			"completed_orders": yearlyTotal.CompletedOrders,
			"total_revenue": yearlyTotal.TotalRevenue,
		},
	})
}

// GetOrder returns a single order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).Preload("User").Preload("OrderItems.Product").First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order (for stand admin to create orders on behalf of students)
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		UserID        uint     `json:"user_id" binding:"required"`
		PaymentMethod string   `json:"payment_method" binding:"required"`
		Items         []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	var user models.User
	if err := h.db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Start transaction
	tx := h.db.Begin()

	// Determine initial status based on payment method
	initialStatus := "payment_pending"
	if req.PaymentMethod == "cash" {
		initialStatus = "request"
	}

	// Create order
	order := models.Order{
		OrderNumber:   "ORD-" + time.Now().Format("20060102150405"),
		UserID:        req.UserID,
		Status:        initialStatus,
		PaymentMethod: req.PaymentMethod,
		StandID:       standID.(uint),
	}

	// Process order items
	var totalAmount float64
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		var product models.Product
		if err := tx.Where("id = ? AND stand_id = ?", item.ProductID, standID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product: " + product.Name})
			return
		}

		// Calculate discounted price
		price := product.Price
		if product.Discount > 0 {
			price = product.Price - (product.Price * product.Discount / 100)
		}

		subtotal := float64(item.Quantity) * price
		totalAmount += subtotal

		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
			Subtotal:  subtotal,
		}
		orderItems = append(orderItems, orderItem)

		// Update product stock
		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}
	}

	// Check user balance (only for card payments)
	if req.PaymentMethod == "card" && user.Balance < totalAmount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Update order total
	order.TotalAmount = totalAmount

	// Create order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}
	}

	// Deduct user balance (only for card payments)
	if req.PaymentMethod == "card" {
		balanceBefore := user.Balance
		user.Balance -= totalAmount
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user balance"})
			return
		}

		// Create transaction record
		transaction := models.Transaction{
			TransactionNumber: "PUR-" + time.Now().Format("20060102150405"),
			UserID:           user.ID,
			Type:             "purchase",
			Amount:           totalAmount,
			BalanceBefore:    balanceBefore,
			BalanceAfter:     user.Balance,
			Description:      "Purchase: " + order.OrderNumber,
			OrderID:          &order.ID,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, order)
}

// UpdateOrderStatus updates order status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
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
		"request":         true,
		"cooking":         true,
		"done":            true,
		"cancelled":       true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = req.Status
	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "status": order.Status})
}

// DeleteOrder deletes an order for the current stand (soft delete)
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND stand_id = ?", id, standID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Only allow deletion of orders that haven't been completed or cancelled
	if order.Status == "done" || order.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete completed or cancelled orders"})
		return
	}

	// Soft delete the order
	if err := h.db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// GetPendingOrders returns all pending orders for the current stand
func (h *OrderHandler) GetPendingOrders(c *gin.Context) {
	standID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var orders []models.Order
	if err := h.db.Where("stand_id = ? AND status IN ?", standID, []string{"payment_pending", "request", "cooking"}).Preload("User").Preload("OrderItems.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
