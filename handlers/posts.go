package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// POST /api/blogs
func CreateBlogHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user_id from middleware
		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		userID, ok := userIDRaw.(float64) // JWT stores as float64
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		type BlogInput struct {
			Title      string   `json:"title"`
			Content    string   `json:"content"`
			Tags       []string `json:"tags"`
			Visibility bool     `json:"visibility"` // true: public, false: private
			Status     string   `json:"status"`     // 'draft' or 'published'
		}

		var blog BlogInput
		if err := c.BodyParser(&blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		query := `
			INSERT INTO blogs (user_id, title, content, tags, visibility, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
		var blogID int
		err := db.QueryRow(
			context.Background(),
			query,
			int(userID), blog.Title, blog.Content, blog.Tags, blog.Visibility, blog.Status,
		).Scan(&blogID)

		if err != nil {
			fmt.Println("Insert error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create blog"})
		}

		return c.JSON(fiber.Map{"message": "Blog created", "blog_id": blogID})
	}
}


// BlogResponse for structured JSON
type BlogResponse struct {
	ID         int      `json:"id"`
	UserID     int      `json:"user_id"`
	Username   string   `json:"username"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	CreatedAt  time.Time   `json:"created_at"`
	Visibility bool     `json:"visibility"`
	Status 	   string   `json:"status"`
}

func GetAllBlogs(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(), `
			SELECT blogs.id, blogs.user_id, users.name, blogs.title, blogs.content, blogs.tags, blogs.created_at, blogs.visibility
			FROM blogs
			JOIN users ON blogs.user_id = users.id
			WHERE blogs.visibility = TRUE AND blogs.status = 'published'
			ORDER BY blogs.created_at DESC;
		`)
		if err != nil {
			fmt.Println("erroeref")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch blogs"})
		}
		defer rows.Close()

		var blogs []BlogResponse
		for rows.Next() {
			var blog BlogResponse
			var tags []string
			err := rows.Scan(&blog.ID, &blog.UserID, &blog.Username, &blog.Title, &blog.Content, &tags, &blog.CreatedAt, &blog.Visibility)
			if err != nil {
				fmt.Println("blog",blog)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing blog data"})
			}
			blog.Tags = tags
			blogs = append(blogs, blog)
		}
		fmt.Println("Blogs: ",blogs)
		return c.JSON(blogs)
	}
}




func GetBlogByID(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blogID := c.Params("id")
		userID, _ := c.Locals("user_id").(float64) // Will be nil if not logged in

		var blog BlogResponse
		var tags []string
		var ownerID int

		err := db.QueryRow(context.Background(), `
			SELECT blogs.id, blogs.user_id, users.name, blogs.title, blogs.content, blogs.tags, blogs.created_at, blogs.visibility
			FROM blogs
			JOIN users ON blogs.user_id = users.id
			WHERE blogs.id = $1;
		`, blogID).Scan(&blog.ID, &ownerID, &blog.Username, &blog.Title, &blog.Content, &tags, &blog.CreatedAt, &blog.Visibility)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}

		// Check visibility if not owner
		if !blog.Visibility && int(userID) != ownerID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized to view this blog"})
		}

		blog.UserID = ownerID
		blog.Tags = tags
		return c.JSON(blog)
	}
}



func GetBlogsByUsername(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		requesterID, _ := c.Locals("user_id").(float64)

		var ownerID int
		err := db.QueryRow(context.Background(), `SELECT id FROM users WHERE username = $1`, username).Scan(&ownerID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		showAll := requesterID == float64(ownerID)

		query := `
			SELECT blogs.id, blogs.user_id, users.name, blogs.title, blogs.content, blogs.tags, blogs.created_at, blogs.visibility
			FROM blogs
			JOIN users ON blogs.user_id = users.id
			WHERE users.username = $1`

		if !showAll {
			query += " AND blogs.visibility = TRUE AND blogs.status = 'published'"
		}

		rows, err := db.Query(context.Background(), query, username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching blogs"})
		}
		defer rows.Close()

		var blogs []BlogResponse
		for rows.Next() {
			var blog BlogResponse
			var tags []string
			err := rows.Scan(&blog.ID, &blog.UserID, &blog.Username, &blog.Title, &blog.Content, &tags, &blog.CreatedAt, &blog.Visibility)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing blog"})
			}
			blog.Tags = tags
			blogs = append(blogs, blog)
		}

		return c.JSON(blogs)
	}
}



func UpdateBlog(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blogID := c.Params("id")
		userID, _ := c.Locals("user_id").(float64)

		var existingUserID int
		err := db.QueryRow(context.Background(), `SELECT user_id FROM blogs WHERE id = $1`, blogID).Scan(&existingUserID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}

		if int(userID) != existingUserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only edit your own blog"})
		}

		var update BlogResponse
		if err := c.BodyParser(&update); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		_, err = db.Exec(context.Background(), `
			UPDATE blogs SET title = $1, content = $2, tags = $3, visibility = $4, status = $5, updated_at = CURRENT_TIMESTAMP
			WHERE id = $6;
		`, update.Title, update.Content, update.Tags, update.Visibility, update.Status, blogID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update blog"})
		}

		return c.JSON(fiber.Map{"message": "Blog updated successfully"})
	}
}




func DeleteBlog(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blogID := c.Params("id")
		userID, _ := c.Locals("user_id").(float64)

		var ownerID int
		err := db.QueryRow(context.Background(), `SELECT user_id FROM blogs WHERE id = $1`, blogID).Scan(&ownerID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}

		if int(userID) != ownerID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized to delete this blog"})
		}

		_, err = db.Exec(context.Background(), `DELETE FROM blogs WHERE id = $1`, blogID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete blog"})
		}

		return c.JSON(fiber.Map{"message": "Blog deleted successfully"})
	}
}

