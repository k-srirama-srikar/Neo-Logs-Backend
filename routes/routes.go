package routes

import (
	"backend/handlers"
    "backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Public routes
	app.Post("/api/login", handlers.LoginHandler(db))
	app.Post("/api/register", handlers.RegisterHandler(db))

	// Protected routes (use JWT middleware)
	app.Use(handlers.JWTMiddleware())

    app.Get("/api/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "You have access!"})
	})

	// Post routes
	app.Get("/api/posts", handlers.GetAllPostsHandler(db))
	app.Post("/api/posts", handlers.CreatePostHandler(db))

	// Comment routes
	app.Post("/api/posts/:postId/comments", handlers.AddCommentHandler(db))
}
