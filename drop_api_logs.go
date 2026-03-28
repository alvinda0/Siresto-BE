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

	// Connect to database
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to database successfully")

	// Drop api_logs table
	fmt.Println("Dropping api_logs table...")
	if err := db.Exec("DROP TABLE IF EXISTS api_logs CASCADE").Error; err != nil {
		log.Fatal("Failed to drop table:", err)
	}

	fmt.Println("✅ Table api_logs dropped successfully!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("1. Restart your server: go run cmd/server/main.go")
	fmt.Println("2. The table will be recreated with UUID schema automatically")
}
