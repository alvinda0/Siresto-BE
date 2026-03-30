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

	log.Println("Starting migration: add payment fields to orders table...")

	// Add payment fields to orders table
	err = db.Transaction(func(tx *gorm.DB) error {
		// Add payment_method column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50)
		`).Error; err != nil {
			return err
		}

		// Add payment_status column with default UNPAID
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS payment_status VARCHAR(50) DEFAULT 'UNPAID'
		`).Error; err != nil {
			return err
		}

		// Add paid_amount column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS paid_amount DECIMAL(15,2) DEFAULT 0
		`).Error; err != nil {
			return err
		}

		// Add change_amount column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS change_amount DECIMAL(15,2) DEFAULT 0
		`).Error; err != nil {
			return err
		}

		// Add payment_note column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS payment_note TEXT
		`).Error; err != nil {
			return err
		}

		// Add paid_at column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP
		`).Error; err != nil {
			return err
		}

		log.Println("✓ Payment fields added successfully")
		return nil
	})

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}
