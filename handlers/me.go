package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// MeHandler returns the authenticated user's info
func MeHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Extract info from claims
	userID := claims["user_id"]
	username := claims["username"]
	email := claims["email"]

	return c.JSON(fiber.Map{
		"user_id":  userID,
		"username": username,
		"email":    email,
	})
}
