package main

import (
    "log"
    "os"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"  // Import the CORS middleware
    "backend/config"
    "backend/routes"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"context"
)


func runSQLScript(db *pgxpool.Pool, scriptPath string) error {
	// Open the SQL script file
	file, err := os.Open(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to open SQL script file: %w", err)
	}
	defer file.Close()

	// Read the SQL script file content
	sql, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read SQL script file: %w", err)
	}

	// Execute the script
	_, err = db.Exec(context.Background(), string(sql))
	if err != nil {
		return fmt.Errorf("failed to execute SQL script: %w", err)
	}

	return nil
}

func main() {
    // Load environment variables
    config.LoadEnv()

    // Connect to the database
    db, err := config.ConnectDB()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    fmt.Println("Database connection successful!")
    defer db.Close()


	// Path to the functions.sql file
	scriptPath := "sql/functions.sql"

	// Run the SQL script
	if err := runSQLScript(db, scriptPath); err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	} else {
		log.Println("SQL script executed successfully!")
	}



    // Initialize Fiber app
    app := fiber.New()

    // Enable CORS (you can configure the options if needed)
    app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000",  // Allow requests from frontend origin
        AllowMethods: "GET,POST,PUT,DELETE",   // Allow these HTTP methods
        AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Allow necessary headers
    }))

    // Set up routes
    routes.SetupRoutes(app, db)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    log.Fatal(app.Listen(":" + port))
}
