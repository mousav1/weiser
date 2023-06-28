package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/mousav1/weiser/app/database"
	"github.com/mousav1/weiser/web"
	"github.com/spf13/viper"
)

func main() {

	// set config
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("This is my first program in Go")
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	// Create the Fiber app
	app := fiber.New(fiber.Config{
		Views: html.New("./app/views", ".html"),
	})

	// app.Use(middleware.MiddlewareChain.Then)

	// set static directory
	app.Static("/static", "./static")

	// Register the routes
	// router, _ := web.NewApp(db)
	// router.SetupRoutes(app
	// wire.Build(providers.ProviderSet)
	// userController := InitializeUserController()
	web.SetupRoutes(app, db)

	// Start the server
	port := viper.GetString("server.port")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
