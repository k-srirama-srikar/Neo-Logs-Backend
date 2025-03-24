package handlers

import (
	"backend/models"
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key") // Replace with a secure value



func isValidEmail(email string) bool {
	// Simple regex for email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}



func isUniqueConstraintViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505" // PostgreSQL unique_violation error code
	}
	return false
}






// LoginHandler authenticates a user and generates a JWT token.
func LoginHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Identifier    string `json:"identifier"` // can be email or username
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		// fmt.Println("hello...")
		// Fetch user from the database
		var id int
		var name, email, password string
		//err := db.QueryRow(context.Background(), "SELECT password FROM users WHERE email = $1", req.Email).Scan(&storedPassword)
		
		query := `SELECT id, name, email, password FROM get_user_by_email($1)`
		err := db.QueryRow(context.Background(), query, req.Identifier).Scan(&id, &name, &email, &password)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Compare passwords
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials Pass"})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":id,
			"email": email,
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
		}

		return c.JSON(fiber.Map{
			"message": "Login Successful",
			"token": tokenString,
			"user": fiber.Map{"id":id, "name":name, "email":email,},})
	}
}

// // RegisterHandler registers a new user in the system.
// func RegisterHandler(db *pgxpool.Pool) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		var user models.User
// 		if err := c.BodyParser(&user); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
// 		}

// 		// Hash the password
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
// 		}

// 		// Insert the user into the database
// 		query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
// 		_, err = db.Exec(context.Background(), query, user.Name, user.Email, string(hashedPassword))
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving user"})
// 		}

// 		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
// 	}
// }


func RegisterHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Struct to parse request body
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Parse and validate the input
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		if req.Name == "" || req.Email == "" || req.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
		}
		if !isValidEmail(req.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		}
		if len(req.Password) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
		}

		// Call InsertUser from the models package
		err = models.InsertUser(db, req.Name, req.Email, string(hashedPassword))
		if err != nil {
			// Check for unique constraint violation (assuming it's on the email)
			if isUniqueConstraintViolation(err) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving user"})
		}

		// Return success response
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
			"user":    req.Name,
			"email":   req.Email,
		})
	}
}
