package main

import (
	"log"
	"project-name/config"
)

// File ini HANYA untuk development - akan menghapus semua data!
// Gunakan dengan hati-hati!
func main() {
	log.Println("⚠️  WARNING: This will DELETE ALL DATA in the database!")
	log.Println("Starting database reset...")

	config.LoadConfig()
	config.ConnectDB()

	// Drop all tables
	log.Println("Dropping existing tables...")
	tables := []string{"categories", "branches", "companies", "users", "roles"}
	for _, table := range tables {
		if err := config.DB.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			log.Printf("Warning: Failed to drop %s table: %v", table, err)
		}
	}

	// Recreate tables
	config.MigrateDB()

	// Seed data
	config.SeedDatabase()

	log.Println("✓ Database reset completed!")
	log.Println("You can now use the seeded accounts from TEST_ACCOUNTS.md")
}
