package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/http/request"
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

func (c *HomeController) SetSessionData(ctx *fiber.Ctx) error {
	req, err := request.New(ctx)
	req.Setsession("name", "moahmmad")
	if err != nil {
		// If there was an error encoding the session data, send a 500 Internal Server Error response
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Set("Content-Type", "application/json")
	return ctx.SendString("ok")
}

func (c *HomeController) GetSessionData(ctx *fiber.Ctx) error {
	req, err := request.New(ctx)
	sessionData := req.Getsession("name")

	// Return the session data as JSON
	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		// If there was an error encoding the session data, send a 500 Internal Server Error response
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Set("Content-Type", "application/json")
	return ctx.Send(jsonData)
}
