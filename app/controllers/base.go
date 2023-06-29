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
