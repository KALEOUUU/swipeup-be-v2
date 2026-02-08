package siswa

import (
	"fmt"
	"net/http"
	"strings"
	"swipeup-admin-v2/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CartHandler handles cart-related requests for students
type CartHandler struct {
	db *gorm.DB
}

// NewCartHandler creates a new CartHandler instance
func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

// GetCart returns the current user's cart with all items
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).Preload("CartItems.Product").First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create empty cart if not exists
			cart = models.Cart{
				UserID:     userID.(uint),
				TotalItems: 0,
				TotalPrice: 0,
				CartItems:  []models.CartItem{},
			}
			c.JSON(http.StatusOK, cart)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart adds a product to the cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get product details
	var product models.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if product is active and has stock
	if !product.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not available"})
		return
	}
	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Calculate price (with discount)
	price := product.Price
	if product.Discount > 0 {
		price = product.Price * (1 - product.Discount/100)
	}

	// Find or create cart
	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new cart
			cart = models.Cart{
				UserID:     userID.(uint),
				TotalItems: 0,
				TotalPrice: 0,
			}
			if err := h.db.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find cart"})
			return
		}
	}

	// Check if cart already has items from different stand (single stand policy)
	if cart.ID != 0 {
		var existingItems []models.CartItem
		if err := h.db.Where("cart_id = ?", cart.ID).Find(&existingItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing cart items"})
			return
		}

		// If cart has items from different stand, reject
		if len(existingItems) > 0 && existingItems[0].StandID != product.StandID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot add products from different stands. Please checkout or clear cart first.",
				"current_stand": existingItems[0].StandID,
				"requested_stand": product.StandID,
			})
			return
		}
	}

	// Check if product already in cart
	var existingItem models.CartItem
	if err := h.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&existingItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Add new item to cart
			cartItem := models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
				Price:     price,
				Subtotal:  price * float64(req.Quantity),
				StandID:   product.StandID,
			}

			if err := h.db.Create(&cartItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check cart item"})
			return
		}
	} else {
		// Update existing item quantity
		existingItem.Quantity += req.Quantity
		existingItem.Subtotal = price * float64(existingItem.Quantity)

		if err := h.db.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
	}

	// Update cart totals
	h.updateCartTotals(cart.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// UpdateCartItem updates the quantity of a cart item
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID := c.Param("id")

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find cart
	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Find cart item
	var cartItem models.CartItem
	if err := h.db.Where("id = ? AND cart_id = ?", itemID, cart.ID).Preload("Product").First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	// Check stock
	if cartItem.Product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Update quantity and subtotal
	cartItem.Quantity = req.Quantity
	cartItem.Subtotal = cartItem.Price * float64(req.Quantity)

	if err := h.db.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	// Update cart totals
	h.updateCartTotals(cart.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

// RemoveFromCart removes an item from the cart
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID := c.Param("id")

	// Find cart
	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Find and delete cart item
	if err := h.db.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	// Update cart totals
	h.updateCartTotals(cart.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

// ClearCart removes all items from the cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Find cart
	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Delete all cart items
	if err := h.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	// Reset cart totals
	cart.TotalItems = 0
	cart.TotalPrice = 0
	if err := h.db.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// Checkout creates orders from cart items
func (h *CartHandler) Checkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		PaymentMethod string  `json:"payment_method" binding:"required"`
		CashAmount    float64 `json:"cash_amount,omitempty"` // Required for cash payment
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment method
	if req.PaymentMethod != "card" && req.PaymentMethod != "cash" && req.PaymentMethod != "qris" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method. Use 'card', 'cash', or 'qris'"})
		return
	}

	// Validate cash amount for cash payment
	if req.PaymentMethod == "cash" && req.CashAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cash amount is required for cash payment"})
		return
	}

	// Get cart with items
	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).Preload("CartItems.Product").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found or empty"})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Get stand ID from first cart item (all items should be from same stand)
	standID := cart.CartItems[0].StandID
	var orderItems []models.OrderItem
	var totalAmount float64

	for _, cartItem := range cart.CartItems {
		// Check stock again
		if cartItem.Product.Stock < cartItem.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product: " + cartItem.Product.Name})
			return
		}

		// Create order item
		orderItem := models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Price,
			Subtotal:  cartItem.Subtotal,
		}

		orderItems = append(orderItems, orderItem)
		totalAmount += cartItem.Subtotal

		// Reduce stock
		cartItem.Product.Stock -= cartItem.Quantity
		if err := h.db.Save(&cartItem.Product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d-%d", userID, standID, time.Now().Unix())

	// Determine initial status and validate payment based on method
	initialStatus := "payment_pending"
	var cashAmount float64
	var paymentProofURL string

	switch req.PaymentMethod {
	case "cash":
		// Validate cash amount is sufficient
		if req.CashAmount < totalAmount {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Insufficient cash. Required: %.0f, Provided: %.0f", totalAmount, req.CashAmount),
			})
			return
		}
		initialStatus = "request"
		cashAmount = req.CashAmount

	case "qris":
		// QRIS requires payment proof upload, so status remains payment_pending
		initialStatus = "payment_pending"

	case "card":
		// Assume card payment is instant, proceed to request
		initialStatus = "request"
	}

	// Create order
	order := models.Order{
		OrderNumber:    orderNumber,
		UserID:         userID.(uint),
		TotalAmount:    totalAmount,
		Status:         initialStatus,
		PaymentMethod:  req.PaymentMethod,
		StandID:        standID,
		OrderItems:     orderItems,
		CashAmount:     cashAmount,
		PaymentProofURL: paymentProofURL,
	}

	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Clear cart after successful checkout
	if err := h.db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		// Log error but don't fail the checkout
		fmt.Printf("Warning: Failed to clear cart items: %v\n", err)
	}

	// Reset cart totals
	cart.TotalItems = 0
	cart.TotalPrice = 0
	if err := h.db.Save(&cart).Error; err != nil {
		// Log error but don't fail the checkout
		fmt.Printf("Warning: Failed to reset cart totals: %v\n", err)
	}

	// Prepare response based on payment method
	response := gin.H{
		"message": "Checkout successful",
		"order":   order, // Single order instead of array
	}

	// For QRIS payment, include QRIS code for the stand
	if req.PaymentMethod == "qris" {
		var settings models.StandSettings
		if err := h.db.Where("stand_id = ?", standID).First(&settings).Error; err != nil {
			// Log error but continue
			fmt.Printf("Warning: Failed to get QRIS for stand %d: %v\n", standID, err)
		} else if settings.QRIS != "" {
			response["qris_code"] = gin.H{
				"stand_id":   standID,
				"qris_code": settings.QRIS,
				"store_name": settings.StoreName,
			}
			response["message"] = "Checkout successful. Please scan QRIS code to complete payment."
		}
	}

	c.JSON(http.StatusCreated, response)
}

// GetQRISCode returns QRIS code for a specific stand
func (h *CartHandler) GetQRISCode(c *gin.Context) {
	standIDStr := c.Param("stand_id")
	standID := parseUint(standIDStr)
	if standID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stand ID"})
		return
	}

	// Get stand settings for QRIS
	var settings models.StandSettings
	if err := h.db.Where("stand_id = ?", standID).First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stand settings not found"})
		return
	}

	if settings.QRIS == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "QRIS code not configured for this stand"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stand_id": standID,
		"qris_code": settings.QRIS,
		"store_name": settings.StoreName,
	})
}

// GetQRISByOrder returns QRIS code for a specific order (alternative to stand_id)
func (h *CartHandler) GetQRISByOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orderIDStr := c.Param("order_id")
	orderID := parseUint(orderIDStr)
	if orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get order and verify ownership
	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Get stand settings for QRIS
	var settings models.StandSettings
	if err := h.db.Where("stand_id = ?", order.StandID).First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stand settings not found"})
		return
	}

	if settings.QRIS == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "QRIS code not configured for this stand"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"stand_id": order.StandID,
		"qris_code": settings.QRIS,
		"store_name": settings.StoreName,
		"total_amount": order.TotalAmount,
	})
}

// UploadPaymentProof uploads payment proof for QRIS payment
func (h *CartHandler) UploadPaymentProof(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orderIDStr := c.Param("order_id")
	orderID := parseUint(orderIDStr)
	if orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get order and verify ownership
	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if order is pending QRIS payment
	if order.PaymentMethod != "qris" || order.Status != "payment_pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is not eligible for payment proof upload"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("payment_proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment proof file is required"})
		return
	}

	// Validate file type (image only)
	if !isValidImageFile(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("payment_proof_%d_%d_%d.png", orderID, userID, time.Now().Unix())
	filepath := "uploads/payment_proofs/" + filename

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment proof"})
		return
	}

	// Update order with payment proof URL and change status
	order.PaymentProofURL = filepath
	order.Status = "request" // Now ready for stand processing

	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment proof uploaded successfully",
		"order": order,
	})
}

// updateCartTotals recalculates and updates cart totals
func (h *CartHandler) updateCartTotals(cartID uint) error {
	var cart models.Cart
	if err := h.db.Preload("CartItems").First(&cart, cartID).Error; err != nil {
		return err
	}

	totalItems := 0
	totalPrice := 0.0

	for _, item := range cart.CartItems {
		totalItems += item.Quantity
		totalPrice += item.Subtotal
	}

	cart.TotalItems = totalItems
	cart.TotalPrice = totalPrice

	return h.db.Save(&cart).Error
}

// Helper functions
func parseUint(s string) uint {
	var result uint
	fmt.Sscanf(s, "%d", &result)
	return result
}

func isValidImageFile(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range validExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}