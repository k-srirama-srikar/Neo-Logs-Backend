package main

import (
	"log"
	"os"

	"/config"
	"/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New()

	// Set up routes
	routes.SetupRoutes(app, db)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(app.Listen(":" + port))
}
