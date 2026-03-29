package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database")

	// Update owner with company_id
	log.Println("Updating OWNER user with company_id...")
	
	result := db.Exec(`
		UPDATE users 
		SET company_id = (
			SELECT id FROM companies WHERE owner_id = users.id LIMIT 1
		)
		WHERE email = 'owner@restaurant.com' AND company_id IS NULL
	`)

	if result.Error != nil {
		log.Fatal("Failed to update owner:", result.Error)
	}

	log.Printf("✓ Updated %d owner(s) with company_id", result.RowsAffected)

	// Verify
	type User struct {
		ID        string  `gorm:"column:id"`
		Email     string  `gorm:"column:email"`
		CompanyID *string `gorm:"column:company_id"`
		BranchID  *string `gorm:"column:branch_id"`
	}

	var owner User
	db.Table("users").Where("email = ?", "owner@restaurant.com").First(&owner)

	log.Println("\nOwner details:")
	log.Printf("- Email: %s", owner.Email)
	log.Printf("- ID: %s", owner.ID)
	if owner.CompanyID != nil {
		log.Printf("- Company ID: %s", *owner.CompanyID)
	} else {
		log.Printf("- Company ID: NULL")
	}
	if owner.BranchID != nil {
		log.Printf("- Branch ID: %s", *owner.BranchID)
	} else {
		log.Printf("- Branch ID: NULL (normal for OWNER)")
	}

	log.Println("\n✅ Update completed successfully!")
	log.Println("Please login again to get new token with company_id")
}
