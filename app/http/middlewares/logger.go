package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	// Create a new instance of logrus logger
	logger := logrus.New()

	// Set the output file for the logger
	file, err := os.OpenFile("./storage/logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Error("Failed to open log file: ", err)
	} else {
		logger.SetOutput(file)
	}

	// Log the request
	logger.WithFields(logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
	}).Info("Request received")

	// Call the next handler
	return c.Next()
}
