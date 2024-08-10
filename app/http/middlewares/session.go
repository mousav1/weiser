package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/cookies"
	"github.com/mousav1/weiser/app/session"
)

func SessionMiddleware(c *fiber.Ctx) error {
	manager := session.GetSessionManager()
	sessionID, err := cookies.GetCookie(c, "weiser_session")
	if err != nil {
		// Start a new session if no session ID is present
		sess := manager.StartSession(c)
		sessionID = sess.ID
	}

	// Check if the session is valid
	if err := manager.CheckExpiration(sessionID); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retrieve session data
	data, err := manager.GetDataBySessionID(sessionID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Store session data in the context
	c.Locals("sessionData", data)

	if !isAuthorized(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to access this session data",
		})
	}

	return c.Next()
}

func isAuthorized(c *fiber.Ctx) bool {
	// Implement your authorization logic here
	return true
}
