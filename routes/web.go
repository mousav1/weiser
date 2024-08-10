package routes

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/controllers"
	"github.com/mousav1/weiser/app/cookies"
	"github.com/mousav1/weiser/app/repositories"
	"github.com/mousav1/weiser/app/services"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type App struct {
	db             *gorm.DB
	userController controllers.UserController
}

func SetupRoutes(app *fiber.App, db *gorm.DB) error {
	userRepo := repositories.NewUserRepository(db)
	if userRepo == nil {
		return errors.New("failed to create user repository")
	}
	userService := services.NewUserService(userRepo)
	if userService == nil {
		return errors.New("failed to create user service")
	}
	userController := controllers.NewUserController(userService)
	if userController == nil {
		return errors.New("failed to create user controller")
	}

	app.Get("/users/:id", userController.GetUserByID)
	app.Post("/users", userController.CreateUser)

	homeController := controllers.NewHomeController()
	if homeController == nil {
		return errors.New("failed to create home controller")
	}

	app.Get("/", homeController.Index)
	app.Get("/show", homeController.ShowView)
	app.Get("/set-session", homeController.SetSessionData)
	app.Get("/get-session", homeController.GetSessionData)

	app.Get("/set", func(c *fiber.Ctx) error {
		username := "123"
		expire := time.Now().Add(24 * time.Hour)
		cookies.SetCookie(c, "username", username, expire)
		return c.SendString("Cookie has been set!")
	})

	app.Get("/get", func(c *fiber.Ctx) error {
		cookie, err := cookies.GetCookie(c, "username")
		if err != nil {
			return c.SendString("Cookie not found")
		}

		return c.SendString(fmt.Sprintf("Cookie value is: %s\n", cookie))
	})

	return nil
}
