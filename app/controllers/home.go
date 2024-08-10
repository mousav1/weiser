package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/cache"
	"github.com/mousav1/weiser/app/http/request"
	"github.com/mousav1/weiser/app/views"
	"github.com/mousav1/weiser/facades"
)

// HomeController is responsible for showing the home page.
type HomeController struct {
	*BaseController
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// NewHomeController creates a new instance of HomeController.
func NewHomeController() *HomeController {
	return &HomeController{}
}

// Index shows the home page.
func (c *HomeController) Index(ctx *fiber.Ctx) error {
	person := Person{Name: "John", Age: 30}
	err := facades.Cache().Set("person", person, 5*time.Minute)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	var cachedPerson Person
	err = facades.Cache().Get("person", &cachedPerson)
	if err != nil {
		if err == cache.ErrCacheMiss {
			fmt.Println("Value not found in cache")
		}
		return ctx.SendString(err.Error())
	}
	cachedPersonJSON, err := json.Marshal(cachedPerson)
	if err != nil {
		return ctx.SendString(err.Error())
	}
	return ctx.SendString(string(cachedPersonJSON))
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

func (c *HomeController) ShowView(ctx *fiber.Ctx) error {
	data := views.ViewData{
		Title: "Home",
		Data:  "John Smith",
	}
	err := views.View(ctx, data, "test.html")
	if err != nil {
		return err
	}
	return nil
}
