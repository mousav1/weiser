package middleware

// import (
// 	"github.com/gofiber/fiber/v2"
// )

// // MiddlewareFunc represents a function that takes a *fiber.Ctx and a next function.
// type MiddlewareFunc func(*fiber.Ctx, func() error) error

// // MiddlewareChain is a chain of middleware functions to be executed in order.
// type MiddlewareChain []MiddlewareFunc

// // Then chains the middleware functions and returns a single fiber.Handler.
// func (chain MiddlewareChain) Then(handler fiber.Handler) fiber.Handler {
// 	if handler == nil {
// 		panic("handler cannot be nil")
// 	}

// 	// Create a new fiber.Handler that executes the middleware chain and the final handler.
// 	return func(c *fiber.Ctx) error {
// 		// Create a function that executes the final handler.
// 		next := func() error {
// 			return handler(c)
// 		}

// 		// Execute the middleware chain.
// 		for i := len(chain) - 1; i >= 0; i-- {
// 			nextMiddleware := chain[i]
// 			next = func(middleware MiddlewareFunc, oldNext func() error) func() error {
// 				return func() error {
// 					return middleware(c, oldNext)
// 				}
// 			}(nextMiddleware, next)
// 		}

// 		// Call the final handler.
// 		return next()
// 	}
// }

// // Define the middleware chain.
// var MiddlewareChain = MiddlewareChain{
// 	LoggerMiddleware,
// 	// Add more middleware here as needed.

// 	AuthMiddleware(models.NewDB()), // Add the AuthMiddleware with the DB instance.
// 	// Add more middleware here as needed.
// }
