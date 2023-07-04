package errorhandler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Handler(ctx *fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"url":    ctx.Path(),
				"method": ctx.Method(),
			}).Error("Internal Server Error")

			// Return error response
			res := ErrorResponse{
				Message: "Internal Server Error",
				Code:    fiber.StatusInternalServerError,
			}

			// Return response in the appropriate format
			if ctx.Accepts("json") != "" {
				ctx.Status(res.Code).JSON(res)
			} else {
				ctx.Status(res.Code).SendString(res.Message)
			}
		}
	}()

	// Proceed to next middleware
	return ctx.Next()
}

func ErrorHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if err := Handler(ctx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Internal Server Error"))
		}
		return ctx.Next()
	}
}
