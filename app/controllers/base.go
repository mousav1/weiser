package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// BaseController is a base controller for all controllers in the application.
type BaseController struct {
}

// NewBaseController creates a new instance of BaseController.
func NewBaseController() *BaseController {
	return &BaseController{}
}

// Render renders the template.
func (c *BaseController) Render(ctx *fiber.Ctx, data interface{}) error {
	// Render the template
	// data := ViewData{
	// 	"Title": "Home",
	// 	"Name":  "John Smith",
	// }
	// err := view(w, data, "index.html")
	// if err != nil {
	// 	return err
	// }

	return nil
}

// func (c *BaseController) Render(ctx *fiber.Ctx, data interface{}, view string) error {
// 	err := view(ctx, data, view+".html")
// 	if err != nil {
// 		return handleError(ctx, err)
// 	}

// 	return nil
// }

// handleError handles the error and returns the appropriate response to the client.
func handleError(ctx *fiber.Ctx, err error) error {
	// Log the error
	// ctx.App().Logger.Error("error message")

	// Return error response
	res := fiber.Map{
		"error": err.Error(),
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(res)
}
