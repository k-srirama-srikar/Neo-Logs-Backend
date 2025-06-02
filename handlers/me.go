package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"fmt"
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





func UpdateUserProfile(db *pgxpool.Pool) fiber.Handler {
    return func(c *fiber.Ctx) error {
        username := c.Params("username")

        // Optional: log to help trace behavior
        fmt.Printf("Updating profile for username: %s\n", username)

        // Parse request body early
        var input struct {
            FullName       string `json:"full_name"`
            Bio            string `json:"bio"`
            Overview       string `json:"overview"`
            ProfilePicture string `json:"profile_picture"`
        }

        if err := c.BodyParser(&input); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Invalid input body")
        }

        // Fetch user ID from username
        var userID string
        err := db.QueryRow(context.Background(), `SELECT id FROM users WHERE name=$1`, username).Scan(&userID)
        if err != nil {
            return fiber.NewError(fiber.StatusNotFound, "User not found")
        }

        // Perform update
        cmdTag, err := db.Exec(context.Background(), `
            UPDATE user_profiles 
            SET full_name=$1, bio=$2, overview=$3, profile_picture=$4, updated_at=NOW()
            WHERE user_id=$5
        `, input.FullName, input.Bio, input.Overview, input.ProfilePicture, userID)

        if err != nil {
            fmt.Printf("DB error: %v\n", err)
            return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile")
        }

        // Optional: ensure something was updated
        if cmdTag.RowsAffected() == 0 {
            return fiber.NewError(fiber.StatusNotFound, "Profile not found for update")
        }

        return c.JSON(fiber.Map{"message": "Profile updated successfully"})
    }
}

