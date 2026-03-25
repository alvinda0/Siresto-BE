package main

import (
	"log"
	"project-name/config"
	"project-name/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	hashPassword := func(password string) string {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashed)
	}

	// Update semua password seeder
	passwords := map[string]string{
		"superadmin@siresto.com":   "admin123",
		"support@siresto.com":      "support123",
		"finance@siresto.com":      "finance123",
		"owner@restaurant.com":     "owner123",
		"admin@restaurant.com":     "admin123",
		"cashier@restaurant.com":   "cashier123",
		"kitchen@restaurant.com":   "kitchen123",
		"waiter@restaurant.com":    "waiter123",
	}

	for email, password := range passwords {
		hashedPassword := hashPassword(password)
		result := config.DB.Model(&entity.User{}).Where("email = ?", email).Update("password", hashedPassword)
		if result.Error != nil {
			log.Printf("Error updating %s: %v", email, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("✓ Updated password for: %s", email)
		} else {
			log.Printf("✗ User not found: %s", email)
		}
	}

	log.Println("Password update completed!")
}
