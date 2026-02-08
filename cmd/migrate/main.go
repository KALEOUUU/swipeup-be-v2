package main

import (
	"fmt"
	"log"
	"os"

	"swipeup-admin-v2/internal/app/database"
	"swipeup-admin-v2/internal/app/models"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Create database if it doesn't exist
	err = createDatabaseIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = migrateDatabase(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully!")
}

func migrateDatabase(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// Auto migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Transaction{},
		&models.StandSettings{},
		&models.GlobalSettings{},
		&models.Cart{},
		&models.CartItem{},
	)

	if err != nil {
		return err
	}

	log.Println("Successfully migrated tables:")
	log.Println("  - users")
	log.Println("  - categories")
	log.Println("  - products")
	log.Println("  - orders")
	log.Println("  - order_items")
	log.Println("  - transactions")
	log.Println("  - stand_settings")
	log.Println("  - global_settings")
	log.Println("  - carts")
	log.Println("  - cart_items")

	// Insert default global settings if they don't exist
	var settingsCount int64
	db.Model(&models.GlobalSettings{}).Count(&settingsCount)

	if settingsCount == 0 {
		log.Println("Inserting default global settings...")
		defaultSettings := []models.GlobalSettings{
			{Key: "discount_rate", Value: "0"},
			{Key: "school_name", Value: "Swipeup School"},
			{Key: "currency", Value: "IDR"},
		}
		if err := db.Create(&defaultSettings).Error; err != nil {
			log.Printf("Warning: Failed to insert default settings: %v", err)
		} else {
			log.Println("  - Default global settings inserted")
		}
	}

	return nil
}

func createDatabaseIfNotExists() error {
	log.Println("Checking if database exists...")
	
	// Get database name from environment variable
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "swipeup" // default database name
	}
	
	// Connect to MySQL without specifying a database
	// Use default values if environment variables are not set
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3307"
	}
	
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = ""
	}
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort)
	
	// Create a new connection to MySQL server
	dbWithoutDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	
	// Create database if it doesn't exist
	log.Printf("Creating database '%s' if it doesn't exist...", dbName)
	err = dbWithoutDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)).Error
	if err != nil {
		return err
	}
	
	log.Printf("Database '%s' created successfully", dbName)
	
	return nil
}
