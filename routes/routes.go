package routes

import (
	"backend/handlers"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Auth routes
	// app.Post("/api/login", handlers.LoginHandler(db))
	// app.Post("/api/register", handlers.RegisterHandler(db))
	// // app.Get("/api/users/:username", handlers.GetUserProfileHandler(db))
	// app.Get("/api/users/:username", handlers.GetUserProfile(db))
	// authRoutes.Post("/follow", FollowUserHandler(db))
	// authRoutes.Post("/unfollow", UnfollowUserHandler(db))
	// app := fiber.New() nt needed

	// Public Routes
	app.Post("/api/login", handlers.LoginHandler(db))
	app.Post("/api/register", handlers.RegisterHandler(db))

	// Protected Routes
	authRoutes := app.Group("/api")
	authRoutes.Use(middleware.JWTMiddleware()) // Apply JWT middleware first
	authRoutes.Use(middleware.ExtractUserID)   // Extract user ID

	authRoutes.Get("/users/:username", handlers.GetUserProfile(db))
	authRoutes.Post("/follow", handlers.FollowUserHandler(db))
	authRoutes.Post("/unfollow", handlers.UnfollowUserHandler(db))


	// Add additional routes for posts, comments, etc.
}
