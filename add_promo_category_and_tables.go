package main

import (
	"fmt"
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

	var dsn string
	
	// Check if DATABASE_URL exists (preferred)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	} else {
		// Fallback to individual env vars
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Starting promo category migration...")

	// 1. Add promo_category column to promos table
	if err := db.Exec(`
		ALTER TABLE promos 
		ADD COLUMN IF NOT EXISTS promo_category VARCHAR(20) NOT NULL DEFAULT 'normal'
	`).Error; err != nil {
		log.Fatal("Failed to add promo_category column:", err)
	}
	log.Println("✓ Added promo_category column to promos table")

	// 2. Create promo_products table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS promo_products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			promo_id UUID NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(promo_id, product_id)
		)
	`).Error; err != nil {
		log.Fatal("Failed to create promo_products table:", err)
	}
	log.Println("✓ Created promo_products table")

	// 3. Create promo_bundles table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS promo_bundles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			promo_id UUID NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			quantity INT NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(promo_id, product_id)
		)
	`).Error; err != nil {
		log.Fatal("Failed to create promo_bundles table:", err)
	}
	log.Println("✓ Created promo_bundles table")

	// 4. Create indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_promo_products_promo_id ON promo_products(promo_id)
	`).Error; err != nil {
		log.Fatal("Failed to create index on promo_products:", err)
	}
	log.Println("✓ Created index on promo_products")

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_promo_products_product_id ON promo_products(product_id)
	`).Error; err != nil {
		log.Fatal("Failed to create index on promo_products:", err)
	}
	log.Println("✓ Created index on promo_products (product_id)")

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_promo_bundles_promo_id ON promo_bundles(promo_id)
	`).Error; err != nil {
		log.Fatal("Failed to create index on promo_bundles:", err)
	}
	log.Println("✓ Created index on promo_bundles")

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_promo_bundles_product_id ON promo_bundles(product_id)
	`).Error; err != nil {
		log.Fatal("Failed to create index on promo_bundles:", err)
	}
	log.Println("✓ Created index on promo_bundles (product_id)")

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_promos_promo_category ON promos(promo_category)
	`).Error; err != nil {
		log.Fatal("Failed to create index on promos:", err)
	}
	log.Println("✓ Created index on promos (promo_category)")

	log.Println("✅ Migration completed successfully!")
	log.Println("\nPromo categories:")
	log.Println("  - normal: Promo umum untuk semua produk")
	log.Println("  - product: Promo untuk produk tertentu")
	log.Println("  - bundle: Promo bundle (beli produk A + B dapat diskon)")
}
