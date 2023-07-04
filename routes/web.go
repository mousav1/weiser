package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/controllers"
	"github.com/mousav1/weiser/app/repositories"
	"github.com/mousav1/weiser/app/services"
	"gorm.io/gorm"
)

type App struct {
	db             *gorm.DB
	userController controllers.UserController
}

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	app.Get("/users/:id", userController.GetUserByID)
	app.Post("/users", userController.CreateUser)

	homeController := controllers.NewHomeController()

	app.Get("/", homeController.Index)
}
