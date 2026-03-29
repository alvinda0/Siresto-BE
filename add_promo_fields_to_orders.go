package main

import (
	"log"
	"project-name/config"

	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	db := config.DB

	log.Println("Adding promo fields to orders table...")

	// Add new columns
	err := db.Transaction(func(tx *gorm.DB) error {
		// Add promo_id column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS promo_id UUID REFERENCES promos(id) ON DELETE SET NULL
		`).Error; err != nil {
			return err
		}

		// Add discount_amount column
		if err := tx.Exec(`
			ALTER TABLE orders 
			ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(15,2) DEFAULT 0
		`).Error; err != nil {
			return err
		}

		// Update existing orders to have discount_amount = 0 if NULL
		if err := tx.Exec(`
			UPDATE orders 
			SET discount_amount = 0 
			WHERE discount_amount IS NULL
		`).Error; err != nil {
			return err
		}

		log.Println("✓ Added promo_id column")
		log.Println("✓ Added discount_amount column")
		log.Println("✓ Updated existing orders")

		return nil
	})

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("✓ Migration completed successfully!")
	log.Println("\nNew order calculation formula:")
	log.Println("  Total = ((Subtotal - Discount) + Tax Priority 1) + Tax Priority 2")
	log.Println("\nNew fields in orders table:")
	log.Println("  - promo_id: UUID (nullable, references promos table)")
	log.Println("  - discount_amount: DECIMAL(15,2) (default 0)")
}
