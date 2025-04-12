package handlers

import (
	"context"
	"fmt"

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
		fmt.Println(c.Locals("user_id"))
		folID := c.Locals("user_id")
		fmt.Println("folid",folID.(float64))
		// followerIDFloat, ok := c.Locals("userID").(float64) // Get logged-in user ID from JWT the user id is stored as float
		// if !ok {
		// 	fmt.Printf("user_id type: %T, value: %v\n", c.Locals("user_id"), c.Locals("user_id")) // Debugging
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: user_id not found or invalid type"})
		// }
		
		var followerID = int(folID.(float64))
		var followingID int

		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE name = $1", req.Username).Scan(&followingID)
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

		fmt.Println(c.Locals("user_id"))
		folID := c.Locals("user_id")
		fmt.Println("folid",folID.(float64))


		// followerID, ok := c.Locals("user_id").(int)
		// if !ok {
		// 	fmt.Printf("user_id type: %T, value: %v\n", c.Locals("user_id"), c.Locals("user_id")) // Debugging
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: user_id not found or invalid type"})
		// }
		// var followingID int

		var followerID = int(folID.(float64))
		var followingID int

		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE name = $1", req.Username).Scan(&followingID)
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



// backend/handlers/user.go
func GetFollowersList(db *pgxpool.Pool) fiber.Handler {
    return func(c *fiber.Ctx) error {
        username := c.Params("username")
		fmt.Println("jgsdskgnrg ",username)

        query := `
            SELECT u.id, u.name, p.full_name, p.profile_picture
            FROM followers f
            JOIN users u ON u.id = f.follower_id
            LEFT JOIN user_profiles p ON p.user_id = u.id
            WHERE f.following_id = (SELECT id FROM users WHERE name = $1)
        `

        rows, err := db.Query(context.Background(), query, username)
        if err != nil {
			fmt.Println("some error")
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to fetch followers",
            })
        }
        defer rows.Close()

        var followers []map[string]interface{}
        for rows.Next() {
            var id int
            var name string
            var fullName, profilePic *string

            if err := rows.Scan(&id, &name, &fullName, &profilePic); err != nil {
                continue
            }

            followers = append(followers, fiber.Map{
                "id":               id,
                "username":         name,
                "full_name":        fullName,
                "profile_picture":  profilePic,
            })
        }

        return c.JSON(fiber.Map{
            "followers": followers,
        })
    }
}




// backend/handlers/user.go
func GetFollowingList(db *pgxpool.Pool) fiber.Handler {
    return func(c *fiber.Ctx) error {
        username := c.Params("username")
		fmt.Println("jgsdskgnrg ",username)

        query := `
            SELECT u.id, u.name, p.full_name, p.profile_picture
            FROM followers f
            JOIN users u ON u.id = f.following_id
            LEFT JOIN user_profiles p ON p.user_id = u.id
            WHERE f.follower_id = (SELECT id FROM users WHERE name = $1)
        `

        rows, err := db.Query(context.Background(), query, username)
        if err != nil {
			fmt.Println("some error")
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to fetch followers",
            })
        }
        defer rows.Close()

        var followers []map[string]interface{}
        for rows.Next() {
            var id int
            var name string
            var fullName, profilePic *string

            if err := rows.Scan(&id, &name, &fullName, &profilePic); err != nil {
                continue
            }

            followers = append(followers, fiber.Map{
                "id":               id,
                "username":         name,
                "full_name":        fullName,
                "profile_picture":  profilePic,
            })
        }

        return c.JSON(fiber.Map{
            "followers": followers,
        })
    }
}