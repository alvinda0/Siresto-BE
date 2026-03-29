package main

import (
	"log"
	"project-name/config"

	"gorm.io/gorm"
)

func main() {
	// Load config and connect to database
	config.LoadConfig()
	config.ConnectDB()

	db := config.DB

	log.Println("Starting migration: Adding tax fields to orders table...")

	// Add new columns to orders table
	err := db.Transaction(func(tx *gorm.DB) error {
		// Add subtotal_amount column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS subtotal_amount DECIMAL(15,2) DEFAULT 0
		`).Error; err != nil {
			return err
		}
		log.Println("✓ Added subtotal_amount column")

		// Add tax_amount column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(15,2) DEFAULT 0
		`).Error; err != nil {
			return err
		}
		log.Println("✓ Added tax_amount column")

		// Update existing orders: set subtotal_amount = total_amount (karena dulu belum ada pajak)
		// dan tax_amount = 0
		if err := tx.Exec(`
			UPDATE orders 
			SET subtotal_amount = total_amount, 
			    tax_amount = 0 
			WHERE subtotal_amount = 0 OR subtotal_amount IS NULL
		`).Error; err != nil {
			return err
		}
		log.Println("✓ Updated existing orders with subtotal and tax amounts")

		return nil
	})

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("✓ Migration completed successfully!")
}
