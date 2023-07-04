package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/utils"
)

// HomeController is responsible for showing the home page.
type HomeController struct {
	*BaseController
}

// NewHomeController creates a new instance of HomeController.
func NewHomeController() *HomeController {
	return &HomeController{}
}

// Index shows the home page.
func (c *HomeController) Index(ctx *fiber.Ctx) error {
	// Render the home page template
	// data := ViewData{
	// 	"Title": "Home",
	// 	"Name":  "John Smith",
	// }
	// err := view(w, data, "index.html")
	// if err != nil {
	// 	return err
	// }
	utils.Info("Hello, World!")
	return ctx.SendString("Hello, World!")

}
