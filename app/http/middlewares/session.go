package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/cookies"
	"github.com/mousav1/weiser/app/session"
)

func SessionMiddleware(c *fiber.Ctx) error {
	// Get the session manager
	manager := session.GetSessionManager()

	// Get session ID from the cookie
	sessionID, err := cookies.GetCookie(c, "weiser_session")
	if err != nil {
		// If session ID is not present, start a new session
		session := manager.StartSession(c)
		sessionID = session.ID
	}

	// Check if session ID is valid and has not expired
	if err := manager.CheckExpiration(sessionID.(string)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get session data from storage using the session ID from the cookie
	data, err := manager.GetDataBySessionID(sessionID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Add session data to the request context
	c.Locals("sessionData", data)

	// Check if the user is authorized to access the session data
	if !isAuthorized(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to access this session data",
		})
	}

	// Call the next handler with the modified request
	return c.Next()
}

func isAuthorized(c *fiber.Ctx) bool {
	// Check the user's role and permissions to determine if they are authorized to access the session data
	// For example, if the session data contains sensitive information, only certain users with a specific role may be authorized to access it.
	// You can use the `c.Locals()` method to retrieve the session data from the request context.
	return true
}

func DeleteExpiredSessions() {
	sessionManager := session.GetSessionManager()
	for {
		// Sleep for 24 hours before deleting all sessions
		time.Sleep(24 * time.Hour)

		// Loop through all the sessions in the storage and delete them
		for _, sessionID := range sessionManager.GetSessionIDs() {
			err := sessionManager.Clear(sessionID)
			if err != nil {
				log.Printf("Error deleting session: %v", err)
			} else {
				log.Printf("Deleted session with ID %v", sessionID)
			}
		}
	}
}
