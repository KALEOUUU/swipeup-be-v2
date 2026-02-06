package main

import (
	"log"
	"swipeup-admin-v2/internal/app/database"
	"swipeup-admin-v2/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db)

	// Start server
	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
