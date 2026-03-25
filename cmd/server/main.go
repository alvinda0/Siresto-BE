package main

import (
	"log"
	"project-name/config"
	"project-name/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()
	config.ConnectDB()
	config.MigrateDB()
	config.SeedDatabase()

	// Setup router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Start server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
