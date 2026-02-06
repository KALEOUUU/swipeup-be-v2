package main

import (
	"log"

	"swipeup-admin-v2/internal/app/database"
	"swipeup-admin-v2/internal/app/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed users
	err = seedUsers(db)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	log.Println("Seeding completed successfully!")
}

func seedUsers(db *gorm.DB) error {
	log.Println("Seeding users...")

	// Define users to seed
	usersToSeed := []struct {
		Name      string
		Email     string
		Phone     string
		Role      string
		Class     string
		Password  string
	}{
		{
			Name:      "Administrator",
			Email:     "admin@swipeup.com",
			Phone:     "081234567890",
			Role:      "admin",
			Class:     "",
			Password:  "Password.1",
		},
		{
			Name:      "Kale Student",
			Email:     "kale@example.com",
			Phone:     "081234567891",
			Role:      "student",
			Class:     "XII RPL 1",
			Password:  "Password.1",
		},
		{
			Name:      "Kale Gmail",
			Email:     "kale@gmail.com",
			Phone:     "081234567892",
			Role:      "student",
			Class:     "XII RPL 2",
			Password:  "Password.1",
		},
		{
			Name:      "Stand Owner",
			Email:     "stand@example.com",
			Phone:     "081234567893",
			Role:      "stand_admin",
			Class:     "",
			Password:  "Password.1",
		},
	}

	for _, userData := range usersToSeed {
		// Check if user already exists
		var existingUser models.User
		result := db.Where("email = ?", userData.Email).First(&existingUser)
		
		if result.Error == nil {
			log.Printf("  ⏭️  User %s (%s) already exists, skipping", userData.Name, userData.Email)
			continue
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("  ❌ Failed to hash password for %s: %v", userData.Name, err)
			continue
		}

		// Create user
		user := models.User{
			Name:      userData.Name,
			Email:     userData.Email,
			Phone:     userData.Phone,
			Role:      userData.Role,
			Class:     userData.Class,
			Balance:   0,
			IsActive:  true,
			RFIDCard:  "",
			Password:  string(hashedPassword),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("  ❌ Failed to create user %s: %v", userData.Name, err)
			continue
		}

		log.Printf("  ✅ Created user: %s (ID: %d, Role: %s)", userData.Name, user.ID, userData.Role)
	}

	log.Println("\n📋 Default credentials for all users:")
	log.Println("  Password: Password.1")
	log.Println("\n🔐 Login credentials:")
	for _, userData := range usersToSeed {
		log.Printf("  - %s: %s / %s (Role: %s)", userData.Name, userData.Email, userData.Password, userData.Role)
	}

	return nil
}
