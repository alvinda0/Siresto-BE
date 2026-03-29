package main

import (
	"log"
	"project-name/config"
	"project-name/internal/entity"
	"project-name/internal/repository"
)

func main() {
	// Load config and connect to database
	config.LoadConfig()
	config.ConnectDB()

	db := config.DB

	log.Println("Starting recalculation of existing orders...")

	// Initialize repositories
	taxRepo := repository.NewTaxRepository(db)

	// Get all orders
	var orders []entity.Order
	if err := db.Preload("OrderItems").Find(&orders).Error; err != nil {
		log.Fatal("Failed to fetch orders:", err)
	}

	log.Printf("Found %d orders to recalculate\n", len(orders))

	successCount := 0
	errorCount := 0

	// Recalculate each order
	for _, order := range orders {
		// Calculate subtotal from order items
		var subtotal float64
		for _, item := range order.OrderItems {
			subtotal += item.Price * float64(item.Quantity)
		}

		// Get active taxes for this order's branch
		taxes, err := taxRepo.FindActiveTaxesByBranch(order.CompanyID, order.BranchID)
		if err != nil {
			log.Printf("Warning: Failed to get taxes for order %s: %v\n", order.ID, err)
			errorCount++
			continue
		}

		// Calculate taxes
		var totalTax float64
		currentAmount := subtotal

		for _, tax := range taxes {
			taxAmount := currentAmount * (tax.Presentase / 100)
			totalTax += taxAmount
			currentAmount += taxAmount
		}

		// Update order
		err = db.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"subtotal_amount": subtotal,
			"tax_amount":      totalTax,
			"total_amount":    subtotal + totalTax,
		}).Error

		if err != nil {
			log.Printf("Error updating order %s: %v\n", order.ID, err)
			errorCount++
		} else {
			log.Printf("✓ Updated order %s: subtotal=%.2f, tax=%.2f, total=%.2f\n", 
				order.ID, subtotal, totalTax, subtotal+totalTax)
			successCount++
		}
	}

	log.Println("\n========================================")
	log.Printf("Recalculation completed!")
	log.Printf("Success: %d orders", successCount)
	log.Printf("Errors: %d orders", errorCount)
	log.Println("========================================")
}
