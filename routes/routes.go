package routes

import (
	"backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Auth routes
	app.Post("/api/login", handlers.LoginHandler(db))
	app.Post("/api/register", handlers.RegisterHandler(db))

	// Add additional routes for posts, comments, etc.
}
