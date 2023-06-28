package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/views"
)

// HomeController is responsible for showing the home page.
type HomeController struct {
	*BaseController
}

// NewHomeController creates a new instance of HomeController.
func NewHomeController(base *BaseController) *HomeController {
	return &HomeController{
		BaseController: base,
	}
}

// Index shows the home page.
func (c *HomeController) Index(ctx *fiber.Ctx) error {
	// Render the home page template
	v := views.NewView("base", "user/edit")
	if err := v.Render(ctx, nil); err != nil {
		return err
	}
	return nil
}