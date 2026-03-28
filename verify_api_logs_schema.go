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
	fmt.Println("")

	// Check if table exists
	var exists bool
	err = db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'api_logs')").Scan(&exists).Error
	if err != nil {
		log.Fatal("Failed to check table:", err)
	}

	if !exists {
		fmt.Println("❌ Table api_logs does NOT exist")
		fmt.Println("   Please start your server to create the table")
		return
	}

	fmt.Println("✅ Table api_logs exists")
	fmt.Println("")

	// Check column types
	type ColumnInfo struct {
		ColumnName string
		DataType   string
	}

	var columns []ColumnInfo
	err = db.Raw(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'api_logs' 
		ORDER BY ordinal_position
	`).Scan(&columns).Error

	if err != nil {
		log.Fatal("Failed to get columns:", err)
	}

	fmt.Println("Table Schema:")
	fmt.Println("─────────────────────────────────────")
	for _, col := range columns {
		status := "✅"
		if col.ColumnName == "id" && col.DataType != "uuid" {
			status = "❌"
		}
		if col.ColumnName == "user_id" && col.DataType != "uuid" {
			status = "❌"
		}
		if col.ColumnName == "company_id" && col.DataType != "uuid" {
			status = "❌"
		}
		if col.ColumnName == "branch_id" && col.DataType != "uuid" {
			status = "❌"
		}
		fmt.Printf("%s %-20s %s\n", status, col.ColumnName, col.DataType)
	}
	fmt.Println("─────────────────────────────────────")
	fmt.Println("")

	// Check if id is UUID
	var idType string
	err = db.Raw("SELECT data_type FROM information_schema.columns WHERE table_name = 'api_logs' AND column_name = 'id'").Scan(&idType).Error
	if err != nil {
		log.Fatal("Failed to check id type:", err)
	}

	if idType == "uuid" {
		fmt.Println("✅ SUCCESS! Table api_logs is using UUID for id column")
		fmt.Println("   You can now use the API without scan errors")
	} else {
		fmt.Printf("❌ ERROR! id column is still %s (should be uuid)\n", idType)
		fmt.Println("   Please run: go run drop_api_logs.go")
		fmt.Println("   Then restart your server")
	}
}
