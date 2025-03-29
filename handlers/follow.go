package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FollowRequest struct {
	Username string `json:"username"`
}

// Follow a user
func FollowUserHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req FollowRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		followerID := c.Locals("userID").(int) // Get logged-in user ID from JWT
		var followingID int

		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", req.Username).Scan(&followingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		_, err = db.Exec(context.Background(), "INSERT INTO followers (follower_id, following_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", followerID, followingID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error following user"})
		}

		return c.JSON(fiber.Map{"message": "Followed successfully"})
	}
}

// Unfollow a user
func UnfollowUserHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req FollowRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		followerID := c.Locals("userID").(int)
		var followingID int

		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", req.Username).Scan(&followingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		_, err = db.Exec(context.Background(), "DELETE FROM followers WHERE follower_id = $1 AND following_id = $2", followerID, followingID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error unfollowing user"})
		}

		return c.JSON(fiber.Map{"message": "Unfollowed successfully"})
	}
}
