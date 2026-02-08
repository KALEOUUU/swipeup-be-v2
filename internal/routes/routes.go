package routes

import (
	"swipeup-admin-v2/internal/api/admin"
	"swipeup-admin-v2/internal/api/auth"
	"swipeup-admin-v2/internal/api/siswa"
	"swipeup-admin-v2/internal/api/stand"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all the application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize handlers
	authHandler := auth.NewAuthHandler(db)
	
		// Admin handlers
	adminUserHandler := admin.NewUserHandler(db)
	adminCategoryHandler := admin.NewCategoryHandler(db)
	adminProductHandler := admin.NewProductHandler(db)
	
	// Student handlers
	siswaUserHandler := siswa.NewUserHandler(db)
	siswaOrderHandler := siswa.NewOrderHandler(db)
	siswaMenuHandler := siswa.NewMenuHandler(db)
	siswaCartHandler := siswa.NewCartHandler(db)
	
	// Stand handlers
	standProductHandler := stand.NewProductHandler(db)
	standOrderHandler := stand.NewOrderHandler(db)
	standSettingsHandler := stand.NewSettingsHandler(db)
	standCategoryHandler := stand.NewCategoryHandler(db)
	
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "Swipeup API is running",
			})
		})
		
		// Auth routes (public)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.POST("/register", adminUserHandler.CreateUser)
		}
		
		// Student routes (protected)
		siswaGroup := v1.Group("/siswa")
		siswaGroup.Use(AuthMiddleware())
		{
			siswaGroup.GET("/profile", siswaUserHandler.GetProfile)
			siswaGroup.GET("/balance", siswaUserHandler.GetBalance)
			siswaGroup.GET("/orders", siswaOrderHandler.GetOrders)
			siswaGroup.GET("/orders/monthly", siswaOrderHandler.GetOrdersByMonth)
			siswaGroup.GET("/orders/:id/receipt", siswaOrderHandler.GetOrderReceipt)
			siswaGroup.POST("/orders", siswaOrderHandler.CreateOrder)
			siswaGroup.DELETE("/orders/:id", siswaOrderHandler.DeleteOrder)
			siswaGroup.GET("/transactions", siswaUserHandler.GetTransactions)
			siswaGroup.GET("/products", siswaMenuHandler.GetProducts)
			
			// Cart management
			cart := siswaGroup.Group("/cart")
			{
				cart.GET("", siswaCartHandler.GetCart)
				cart.POST("/items", siswaCartHandler.AddToCart)
				cart.PUT("/items/:id", siswaCartHandler.UpdateCartItem)
				cart.DELETE("/items/:id", siswaCartHandler.RemoveFromCart)
				cart.DELETE("", siswaCartHandler.ClearCart)
				cart.POST("/checkout", siswaCartHandler.Checkout)
				cart.GET("/qris/:stand_id", siswaCartHandler.GetQRISCode)
				cart.GET("/orders/:order_id/qris", siswaCartHandler.GetQRISByOrder)
				cart.POST("/orders/:order_id/payment-proof", siswaCartHandler.UploadPaymentProof)
			}
		}
		
		// Stand admin routes (protected + stand role)
		standGroup := v1.Group("/stand")
		standGroup.Use(AuthMiddleware(), StandMiddleware())
		{
			// Product management
			products := standGroup.Group("/products")
			{
				products.GET("", standProductHandler.GetProducts)
				products.GET("/:id", standProductHandler.GetProduct)
				products.POST("", standProductHandler.CreateProduct)
				products.PUT("/:id", standProductHandler.UpdateProduct)
				products.DELETE("/:id", standProductHandler.DeleteProduct)
				products.PUT("/:id/status", standProductHandler.UpdateProductStatus)
			}
			
			// Order management
			orders := standGroup.Group("/orders")
			{
				orders.GET("", standOrderHandler.GetOrders)
				orders.GET("/pending", standOrderHandler.GetPendingOrders)
				orders.GET("/:id", standOrderHandler.GetOrder)
				orders.POST("", standOrderHandler.CreateOrder)
				orders.PUT("/:id/status", standOrderHandler.UpdateOrderStatus)
				orders.DELETE("/:id", standOrderHandler.DeleteOrder)
				orders.GET("/monthly", standOrderHandler.GetOrdersByMonth)
				orders.GET("/revenue/monthly", standOrderHandler.GetMonthlyRevenueRecap)
			}
			
			// Category management
			categories := standGroup.Group("/categories")
			{
				categories.GET("", standCategoryHandler.GetCategories)
			}
			
			// Settings management
			settings := standGroup.Group("/settings")
			{
				settings.GET("", standSettingsHandler.GetSettings)
				settings.PUT("", standSettingsHandler.UpdateSettings)
				settings.PUT("/qris", standSettingsHandler.UpdateQRIS)
				settings.PUT("/store-name", standSettingsHandler.UpdateStoreName)
			}
		}
		
		// Admin routes (protected + admin role)
		adminGroup := v1.Group("/admin")
		adminGroup.Use(AuthMiddleware(), AdminMiddleware())
		{
			// User management
			users := adminGroup.Group("/users")
			{
				users.GET("", adminUserHandler.GetUsers)
				users.GET("/:id", adminUserHandler.GetUser)
				users.POST("", adminUserHandler.CreateUser)
				users.PUT("/:id", adminUserHandler.UpdateUser)
				users.DELETE("/:id", adminUserHandler.DeleteUser)
				users.POST("/:id/topup", adminUserHandler.TopUpBalance)
			}
			
			// Category management
			categories := adminGroup.Group("/categories")
			{
				categories.GET("", adminCategoryHandler.GetCategories)
				categories.GET("/:id", adminCategoryHandler.GetCategory)
				categories.POST("", adminCategoryHandler.CreateCategory)
				categories.PUT("/:id", adminCategoryHandler.UpdateCategory)
				categories.DELETE("/:id", adminCategoryHandler.DeleteCategory)
			}
			
			// Stand canteen management
			standCanteens := adminGroup.Group("/stand-canteens")
			{
				standCanteens.GET("", adminProductHandler.GetStandCanteens)
				standCanteens.GET("/:id", adminProductHandler.GetStandCanteen)
				standCanteens.POST("", adminProductHandler.CreateStandCanteen)
				standCanteens.PUT("/:id", adminProductHandler.UpdateStandCanteen)
				standCanteens.DELETE("/:id", adminProductHandler.DeleteStandCanteen)
			}
			
			// Global settings management
			globalSettings := adminGroup.Group("/global-settings")
			{
				globalSettings.GET("", adminCategoryHandler.GetGlobalSettings)
				globalSettings.PUT("/:key", adminCategoryHandler.UpdateGlobalSetting)
			}
		}
	}
}
