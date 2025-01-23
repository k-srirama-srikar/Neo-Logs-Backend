package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// JWTMiddleware sets up the JWT authentication middleware
func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("your_secret_key"), // Replace "your_secret_key" with your actual secret key
		ErrorHandler: jwtError,                 // Custom error handler
	})
}

// jwtError handles JWT errors and returns a JSON response
func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}
