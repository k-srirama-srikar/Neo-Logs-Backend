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
	// authRoutes.Post("/follow", FollowUserHandler(db))
	// authRoutes.Post("/unfollow", UnfollowUserHandler(db))
	// app := fiber.New() nt needed
	
	// Public Routes
	app.Post("/api/login", handlers.LoginHandler(db))
	app.Post("/api/register", handlers.RegisterHandler(db))
	// app.Get("/api/users/:username", handlers.GetUserProfile(db))
	// backend/routes/user_routes.go
	// app.Get("/users/:username/followers", handlers.GetFollowersList(db))

	public := app.Group("/api")
	public.Use(middleware.ExtractUserID)
	public.Get("/users/:username/followers", handlers.GetFollowersList(db))
	public.Get("/users/:username/following", handlers.GetFollowingList(db))
	// Protected Routes
	authRoutes := app.Group("/api")
	authRoutes.Use(middleware.JWTMiddleware()) // Apply JWT middleware first
	authRoutes.Use(middleware.ExtractUserID)   // Extract user ID

	public.Get("/users/:username", handlers.GetUserProfile(db))
	authRoutes.Post("/follow", middleware.JWTMiddleware(), handlers.FollowUserHandler(db))
	authRoutes.Post("/unfollow", middleware.JWTMiddleware(), handlers.UnfollowUserHandler(db))



	// Add additional routes for posts, comments, etc.
}
