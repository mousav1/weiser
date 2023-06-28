package middleware

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"net/http"
// 	"strings"
// )

// // AuthMiddleware is a middleware that verifies the user's authentication token.
// func AuthMiddleware(db *models.DB) MiddlewareFunc {
// 	return func(c *fiber.Ctx, next func() error) error {
// 		// Get the authentication token from the Authorization header.
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Missing Authorization header"})
// 		}

// 		// Split the authorization header into two parts: the scheme and the token.
// 		parts := strings.SplitN(authHeader, " ", 2)
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Authorization header"})
// 		}

// 		// Verify the authentication token using the database.
// 		token := parts[1]
// 		user, err := db.VerifyToken(token)
// 		if err != nil {
// 			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid authentication token"})
// 		}

// 		// Set the user information in the context for other handlers to use.
// 		c.Locals("user", user)

// 		// Call the next middleware/handler in the chain.
// 		return next()
// 	}
// }
