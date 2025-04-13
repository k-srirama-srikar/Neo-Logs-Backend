// package middleware

// import (
// 	"github.com/gofiber/fiber/v2"
// 	jwtware "github.com/gofiber/jwt/v3"
// )

// // JWTMiddleware sets up the JWT authentication middleware
// func JWTMiddleware() fiber.Handler {
// 	return jwtware.New(jwtware.Config{
// 		SigningKey:   []byte("your_secret_key"), // Replace "your_secret_key" with your actual secret key
// 		ErrorHandler: jwtError,                 // Custom error handler
// 	})
// }

// // jwtError handles JWT errors and returns a JSON response
// func jwtError(c *fiber.Ctx, err error) error {
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
// 	}
// 	return c.Next()
// }

package middleware

import (
	"fmt"
	"strings"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// JWTMiddleware sets up the JWT authentication middleware
func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("your_secret_key"), // Ensure this matches your token generator
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer", // Required since frontend sends "Bearer <token>"
		ContextKey:   "user",   // This will store the decoded JWT in c.Locals("user")
		ErrorHandler: jwtError,
	})
}

// jwtError handles JWT errors and returns a JSON response
func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized 1"})
	}
	return c.Next()
}

// ExtractUserID Middleware (Optional, to get user_id in handlers)
func ExtractUserID(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		// No token sent → anonymous
		c.Locals("user_id", nil)
		return c.Next()
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Make sure token method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil || !token.Valid {
		c.Locals("user_id", nil) // treat as anonymous if invalid
		return c.Next()
	}

	// Extract user_id from claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, exists := claims["user_id"]; exists {
			fmt.Printf("DEBUG: Middleware set user_id: %v (Type: %T)\n", userID, userID)
			c.Locals("user_id", userID)
			return c.Next()
		}
	}
	fmt.Println("No token present... Going for anonymous access")
	c.Locals("user_id", nil) // default to anonymous
	return c.Next()
	// userToken := c.Locals("user")
	// // If no token is present, allow anonymous access
	// if userToken == nil {
	// 	fmt.Println("DEBUG: No JWT token found, proceeding as anonymous user")
	// 	c.Locals("user_id", nil)
	// 	return c.Next()
	// }

	// fmt.Println("User: " ,c.Locals("user"))
	// // if userToken == nil {
	// // 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized 23"})
	// // }

	// // Convert userToken to JWT claims
	// claims, ok := userToken.(*jwt.Token)
	// if !ok {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	// }

	// // Extract user_id from claims
	// mapClaims, ok := claims.Claims.(jwt.MapClaims)
	// if !ok || !claims.Valid {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	// }

	// fmt.Println("mapClaims", mapClaims)
	// fmt.Println("mapClaims", mapClaims["user_id"])

	// userID, exists := mapClaims["user_id"]
	// if !exists {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID missing in token"})
	// }

	// // Store user_id in context for handlers
	// fmt.Println(userID)
	// c.Locals("user_id", userID)
	// fmt.Printf("DEBUG: Middleware set user_id: %v (Type: %T)\n", userID, userID)
	// print(c.Locals("user_id"))
	// return c.Next()
}
