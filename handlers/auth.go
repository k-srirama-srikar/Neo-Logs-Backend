package handlers

import (
	"backend/models"
	"context"
	"fmt"
	"regexp"
	"time"
	// "net/http"
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



// GetUserProfileHandler retrieves user profile based on username
// func GetUserProfileHandler(db *pgxpool.Pool) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		username := c.Params("username")

// 		// Query to fetch user details from users and user_profiles tables
// 		var user struct {
// 			ID             int     `json:"id"`
// 			Username       string  `json:"name"`
// 			Email          string  `json:"email"`
// 			FullName       *string `json:"full_name"`
// 			Bio            *string `json:"bio"`
// 			ProfilePicture *string `json:"profile_picture"`
// 			Public         bool    `json:"public"`
// 		}

// 		query := `
// 			SELECT u.id, u.name, u.email, p.full_name, p.bio, p.profile_picture, p.public
// 			FROM users u
// 			LEFT JOIN user_profiles p ON u.id = p.user_id
// 			WHERE u.name = $1
// 		`

// 		err := db.QueryRow(context.Background(), query, username).Scan(
// 			&user.ID, &user.Username, &user.Email, &user.FullName,
// 			&user.Bio, &user.ProfilePicture, &user.Public,
// 		)

// 		if err != nil {
// 			fmt.Printf("Error fetching user profile: %v", err)
// 			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
// 		}

// 		// Return user profile data as JSON
// 		return c.JSON(user)
// 	}
// }



// func GetUserProfile(db *pgxpool.Pool) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		username := c.Params("username")
// 		authUserID, _ := c.Locals("user_id").(int) // Get logged-in user ID (if any)

// 		var user struct {
// 			ID        int    `json:"id"`
// 			Username  string `json:"username"`
// 			FullName  string `json:"full_name"`
// 			Bio       string `json:"bio"`
// 			Followers int    `json:"followers"`
// 			Following int    `json:"following"`
// 		}

// 		query := `
// 		SELECT u.id, u.name, up.full_name, up.bio,
// 		       (SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers,
// 		       (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following
// 		FROM users u
// 		LEFT JOIN user_profiles up ON u.id = up.user_id
// 		WHERE u.name = $1`
// 		err := db.QueryRow(context.Background(), query, username).Scan(
// 			&user.ID, &user.Username, &user.FullName, &user.Bio, &user.Followers, &user.Following,
// 		)
// 		if err != nil {
// 			fmt.Printf("Error fetching user profile: %v\n", err)
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
// 		}

// 		// Check if logged-in user is viewing their own profile
// 		isOwner := authUserID == user.ID

// 		return c.JSON(fiber.Map{
// 			"user":     user,
// 			"is_owner": isOwner,
// 		})
// 	}
// }



// func GetUserProfileHandler(db *pgxpool.Pool) fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         username := c.Params("username")
//         userID := c.Locals("userID").(int) // Extract from JWT

//         var profile UserProfile
//         err := db.QueryRow(context.Background(), `
//             SELECT u.id, u.username, p.full_name, p.bio, p.profile_picture,
//                    (SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers,
//                    (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following,
//                    EXISTS (SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = u.id) AS is_following
//             FROM users u
//             LEFT JOIN user_profiles p ON u.id = p.user_id
//             WHERE u.username = $2`, userID, username).
//             Scan(&profile.ID, &profile.Username, &profile.FullName, &profile.Bio, &profile.ProfilePicture, &profile.Followers, &profile.Following, &profile.IsFollowing)

//         if err != nil {
//             return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
//         }

//         return c.JSON(fiber.Map{"user": profile})
//     }
// }



// func GetUserProfile(db *pgxpool.Pool) fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         username := c.Params("username")
//         userID, ok := c.Locals("userID").(int) // Extract from JWT
// 		fmt.Println("userID")
//         if !ok {
//             return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
//         }

//         var profile struct {
//             ID             int     `json:"id"`
//             Username       string  `json:"username"`
//             FullName       *string `json:"full_name"`
//             Bio            *string `json:"bio"`
//             ProfilePicture *string `json:"profile_picture"`
//             Followers      int     `json:"followers"`
//             Following      int     `json:"following"`
//             IsFollowing    bool    `json:"is_following"`
//         }

//         err := db.QueryRow(context.Background(), `
//             SELECT u.id, u.name, p.full_name, p.bio, p.profile_picture,
//                    (SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers,
//                    (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following,
//                    EXISTS (SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = u.id) AS is_following
//             FROM users u
//             LEFT JOIN user_profiles p ON u.id = p.user_id
//             WHERE u.name = $2`, userID, username).
//             Scan(&profile.ID, &profile.Username, &profile.FullName, &profile.Bio, &profile.ProfilePicture, &profile.Followers, &profile.Following, &profile.IsFollowing)

//         if err != nil {
// 			fmt.Printf("Error fetching user profile: %v\n", err)
//             return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
//         }

//         // Check if the logged-in user is viewing their own profile
//         isOwner := userID == profile.ID

//         return c.JSON(fiber.Map{
//             "user":     profile,
//             "is_owner": isOwner,  // true if the logged-in user is accessing their own profile
//         })
//     }
// }


func GetUserProfile(db *pgxpool.Pool) fiber.Handler {
    return func(c *fiber.Ctx) error {
        username := c.Params("username")

        // Get userID from JWT (if available)
        // userID, _ := c.Locals("userID").(int) // Now optional
		// Get userID from JWT (if available)
        var userID int
		uid := c.Locals("user_id")
		if uid!=nil{
		userID = int(uid.(float64))}else{
			userID=0
		}
        // if uid, ok := c.Locals("userID").(float64); ok { // JWT usually stores float64
        //     convertedID := int(uid)
        //     userID = &convertedID
        // }

        var profile struct {
            ID             int     `json:"id"`
            Username       string  `json:"username"`
            FullName       *string `json:"full_name"`
            Bio            *string `json:"bio"`
            ProfilePicture *string `json:"profile_picture"`
            Followers      int     `json:"followers"`
            Following      int     `json:"following"`
            IsFollowing    *bool   `json:"is_following,omitempty"` // Pointer to allow null
        }

        query := `
            SELECT u.id, u.name, p.full_name, p.bio, p.profile_picture,
                   (SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers,
                   (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following
            FROM users u
            LEFT JOIN user_profiles p ON u.id = p.user_id
            WHERE u.name = $1`
        
        err := db.QueryRow(context.Background(), query, username).
            Scan(&profile.ID, &profile.Username, &profile.FullName, &profile.Bio, &profile.ProfilePicture, &profile.Followers, &profile.Following)
        
        if err != nil {
            fmt.Printf("Error fetching user profile: %v\n", err)
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
        }

        // If logged in, check follow status
        if  userID != 0 {
            err = db.QueryRow(context.Background(), `
                SELECT EXISTS (SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)`,
                userID, profile.ID).Scan(&profile.IsFollowing)
            if err != nil {
                fmt.Printf("Error checking follow status: %v\n", err)
                profile.IsFollowing = nil // Set to nil if query fails
            }
        }

		fmt.Printf("User ID: %v, Profile ID: %d, Is Following: %v\n", userID, profile.ID, profile.IsFollowing)


        // Determine if logged-in user is viewing their own profile
        isOwner :=  userID == profile.ID

        return c.JSON(fiber.Map{
            "user":     profile,
            "is_owner": isOwner,  // true if the logged-in user is accessing their own profile
        })
    }
}
