package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/cache"
	kernel "github.com/mousav1/weiser/app/http"
	middleware "github.com/mousav1/weiser/app/http/middlewares"
	"github.com/mousav1/weiser/app/session"
	"github.com/mousav1/weiser/database"
	"github.com/mousav1/weiser/routes"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func main() {

	// set config
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read configuration file: %s", err)
	}

	// Connect to the database
	db, err := database.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("failed to close database connection: %s", err)
		}
	}()

	// Initialize the session manager
	if err := session.InitSessionManager(); err != nil {
		log.Fatalf("failed to initialize session manager: %s", err)
	}
	go middleware.DeleteExpiredSessions()

	// Create a new in-memory cache with a default expiration of 1 minute
	cache.NewCache(time.Minute, nil)

	// Create the Fiber app
	app := fiber.New(fiber.Config{})

	// add middlewares
	for _, middleware := range kernel.Middleware {
		app.Use(middleware)
	}

	// set static directory
	app.Static("/static", "./static")

	// Register the routes
	if err := routes.SetupRoutes(app, db); err != nil {
		log.Fatalf("failed to set up routes: %s", err)
	}

	// Start the server
	port := viper.GetString("server.port")
	if port == "" {
		port = "3000"
	}

	// define your routes and middleware here
	server := &fasthttp.Server{
		Handler: app.Handler(),
	}

	go func() {
		if err := server.ListenAndServe(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("failed to start server: %s", err)
		}
	}()

	// Wait for signals to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("signal received: %v, shutting down server...\n", sig)

	// Create a context with a timeout of 5 seconds
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown server: %s", err)
	}
}
