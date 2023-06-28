package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/views"
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
	v := views.NewView("base", "user/edit")
	if err := v.Render(ctx, data); err != nil {
		return err
	}

	return nil
}
