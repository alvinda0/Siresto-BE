package main

import (
	"log"
	"project-name/config"
	"project-name/routes"

	"github.com/gin-contrib/cors"
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

	// Setup CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.GetAllowedOrigins()
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	routes.SetupRoutes(r)

	// Start server
	log.Println("Server running on port 8080")
	log.Println("Allowed origins:", config.GetAllowedOrigins())
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
