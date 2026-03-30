package main

import (
	"fmt"
	"log"
	"project-name/config"

	_ "github.com/lib/pq"
)

func main() {
	// Load config
	config.LoadConfig()

	// Update PREPARING to PROCESSING
	result := config.DB.Exec(`
		UPDATE orders 
		SET status = 'PROCESSING' 
		WHERE status = 'PREPARING'
	`)

	if result.Error != nil {
		log.Fatalf("Failed to update order status: %v", result.Error)
	}

	rowsAffected := result.RowsAffected
	fmt.Printf("✓ Successfully updated %d orders from PREPARING to PROCESSING\n", rowsAffected)

	if rowsAffected == 0 {
		fmt.Println("ℹ No orders with PREPARING status found")
	}
}
