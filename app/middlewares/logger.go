package middleware

// import (
// 	"fmt"
// 	"github.com/gofiber/fiber/v2"
// )

// // LoggerMiddleware logs the request and response.
// func LoggerMiddleware(c *fiber.Ctx, next func() error) error {
// 	// Log the request.
// 	fmt.Printf("Request: %s %s\n", c.Method(), c.Path())

// 	// Call the next middleware/handler in the chain.
// 	err := next()

// 	// Log the response.
// 	fmt.Printf("Response: %d\n", c.Response().StatusCode())

// 	return err
// }
