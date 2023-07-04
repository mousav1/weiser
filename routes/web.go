package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/controllers"
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

	return nil
}
