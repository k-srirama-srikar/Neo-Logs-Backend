package routes

import (
	"handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Public routes
	app.Post("/api/login", handlers.LoginHandler(db))
	app.Post("/api/register", handlers.RegisterHandler(db))

	// Protected routes (use JWT middleware)
	app.Use(handlers.JWTMiddleware())

	// Post routes
	app.Get("/api/posts", handlers.GetAllPostsHandler(db))
	app.Post("/api/posts", handlers.CreatePostHandler(db))

	// Comment routes
	app.Post("/api/posts/:postId/comments", handlers.AddCommentHandler(db))
}
