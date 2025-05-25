package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Struct for comment input
type CommentInput struct {
	Content          string `json:"content"`
	ParentCommentID  *int   `json:"parent_comment_id"` // nullable
}

// POST /api/blogs/:id/comments
func CreateCommentHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("hello??")
		blgID, err := strconv.ParseFloat(c.Params("id"), 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid blog ID")
		}

		var blogID = int(blgID)

		userID :=  int(c.Locals("user_id").(float64))

		var input CommentInput
		if err := c.BodyParser(&input); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
		}

		depth := 0
		if input.ParentCommentID != nil {
			// Fetch parent comment to get depth
			err = db.QueryRow(context.Background(),
				"SELECT depth FROM comments WHERE id=$1", *input.ParentCommentID).Scan(&depth)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Parent comment not found")
			}
			depth++
		}

		_, err = db.Exec(context.Background(),
			`INSERT INTO comments (blog_id, user_id, parent_comment_id, content, depth)
			 VALUES ($1, $2, $3, $4, $5)`,
			blogID, userID, input.ParentCommentID, input.Content, depth)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create comment")
		}

		return c.SendStatus(fiber.StatusCreated)
	}
}



// type Comment struct {
// 	ID        int        `json:"id"`
// 	UserID    int        `json:"user_id"`
// 	Content   string     `json:"content"`
// 	ParentID  *int       `json:"parent_comment_id"`
// 	Depth     int        `json:"depth"`
// 	CreatedAt time.Time  `json:"created_at"`
// }

// func GetCommentsHandler(db *pgxpool.Pool) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		blogID, err := strconv.Atoi(c.Params("id"))
// 		if err != nil {
// 			return fiber.NewError(fiber.StatusBadRequest, "Invalid blog ID")
// 		}

// 		rows, err := db.Query(context.Background(),
// 			`SELECT id, user_id, content, parent_comment_id, depth, created_at
// 			 FROM comments WHERE blog_id=$1 ORDER BY created_at ASC`, blogID)
// 		if err != nil {
// 			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch comments")
// 		}
// 		defer rows.Close()

// 		var comments []Comment
// 		for rows.Next() {
// 			var cmt Comment
// 			if err := rows.Scan(&cmt.ID, &cmt.UserID, &cmt.Content, &cmt.ParentID, &cmt.Depth, &cmt.CreatedAt); err != nil {
// 				return fiber.NewError(fiber.StatusInternalServerError, "Scan error")
// 			}
// 			comments = append(comments, cmt)
// 		}
// 		return c.JSON(comments)
// 	}
// }



type Comment struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	UserName       string     `json:"user_name"`
	ProfilePicture string     `json:"profile_picture"`
	Content        string     `json:"content"`
	ParentID       *int       `json:"parent_comment_id"`
	Depth          int        `json:"depth"`
	CreatedAt      time.Time  `json:"created_at"`
	Children       []Comment  `json:"children"` // For nesting
}


func GetCommentsHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("heleoooooo?")
		blgID, err := strconv.ParseFloat(c.Params("id"), 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid blog ID")
		}

		var blogID = int(blgID)

		rows, err := db.Query(context.Background(), `
			SELECT c.id, c.user_id, u.name, p.profile_picture, c.content, c.parent_comment_id, c.depth, c.created_at
			FROM comments c
			JOIN users u ON c.user_id = u.id
			JOIN user_profiles p ON u.id = p.user_id
			WHERE c.blog_id = $1
			ORDER BY c.created_at ASC`, blogID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch comments")
		}
		defer rows.Close()

		// commentMap := make(map[int]*Comment)
		// fmt.Println(commentMap)
		// var rootComments []Comment
		// fmt.Println(rows)
		// for rows.Next() {
		// 	var cmt Comment
		// 	if err := rows.Scan(&cmt.ID, &cmt.UserID, &cmt.UserName, &cmt.ProfilePicture, &cmt.Content, &cmt.ParentID, &cmt.Depth, &cmt.CreatedAt); err != nil {
		// 		return fiber.NewError(fiber.StatusInternalServerError, "Scan error")
		// 	}
		// 	cmt.Children = []Comment{}
		// 	commentMap[cmt.ID] = &cmt

		// 	if cmt.ParentID == nil {
		// 		rootComments = append(rootComments, cmt)
		// 	} else {
		// 		parent := commentMap[*cmt.ParentID]
		// 		fmt.Println("parent", parent)
		// 		if parent != nil {
		// 			parent.Children = append(parent.Children, cmt)
		// 		}
		// 	}
		// }
		// fmt.Println("root ",rootComments)
		// return c.JSON(rootComments)
		commentMap := make(map[int]*Comment)
		var rootComments []*Comment

		for rows.Next() {
			var cmt Comment
			if err := rows.Scan(&cmt.ID, &cmt.UserID, &cmt.UserName, &cmt.ProfilePicture, &cmt.Content, &cmt.ParentID, &cmt.Depth, &cmt.CreatedAt); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Scan error")
			}
			cmt.Children = []Comment{}
			commentMap[cmt.ID] = &cmt

			if cmt.ParentID == nil {
				rootComments = append(rootComments, &cmt)
			} else {
				parent := commentMap[*cmt.ParentID]
				if parent != nil {
					parent.Children = append(parent.Children, cmt)
				}
			}
		}

		return c.JSON(rootComments)

	}
}




func DeleteCommentHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		commentID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid comment ID")
		}

		userID := c.Locals("user_id").(int)

		var commentOwner int
		err = db.QueryRow(context.Background(),
			"SELECT user_id FROM comments WHERE id=$1", commentID).Scan(&commentOwner)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Comment not found")
		}

		if commentOwner != userID {
			return fiber.NewError(fiber.StatusForbidden, "Not allowed")
		}

		_, err = db.Exec(context.Background(),
			"DELETE FROM comments WHERE id=$1", commentID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete comment")
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
