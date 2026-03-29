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
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	} else {
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

	log.Println("Checking promo data...")

	// Check promos table structure
	var columns []struct {
		ColumnName string
		DataType   string
	}
	
	err = db.Raw(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'promos' 
		ORDER BY ordinal_position
	`).Scan(&columns).Error
	
	if err != nil {
		log.Fatal("Failed to get columns:", err)
	}

	log.Println("\n=== Promos Table Structure ===")
	for _, col := range columns {
		log.Printf("  %s: %s", col.ColumnName, col.DataType)
	}

	// Check if promo_category column exists
	var hasPromoCategory bool
	err = db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'promos' AND column_name = 'promo_category'
		)
	`).Scan(&hasPromoCategory).Error
	
	if err != nil {
		log.Fatal("Failed to check promo_category:", err)
	}

	log.Printf("\n=== promo_category column exists: %v ===\n", hasPromoCategory)

	// Check promo data
	var promos []map[string]interface{}
	err = db.Raw(`
		SELECT id, name, code, promo_category, company_id, branch_id 
		FROM promos 
		LIMIT 5
	`).Scan(&promos).Error
	
	if err != nil {
		log.Fatal("Failed to get promos:", err)
	}

	log.Println("=== Sample Promos ===")
	for i, promo := range promos {
		log.Printf("%d. ID: %v, Name: %v, Category: %v", i+1, promo["id"], promo["name"], promo["promo_category"])
	}

	// Check promo_products table
	var hasPromoProducts bool
	err = db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'promo_products'
		)
	`).Scan(&hasPromoProducts).Error
	
	if err != nil {
		log.Fatal("Failed to check promo_products:", err)
	}

	log.Printf("\n=== promo_products table exists: %v ===\n", hasPromoProducts)

	// Check promo_bundles table
	var hasPromoBundles bool
	err = db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'promo_bundles'
		)
	`).Scan(&hasPromoBundles).Error
	
	if err != nil {
		log.Fatal("Failed to check promo_bundles:", err)
	}

	log.Printf("=== promo_bundles table exists: %v ===\n", hasPromoBundles)

	log.Println("\n✅ Check completed!")
}
