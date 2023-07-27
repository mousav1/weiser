package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/compress"
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

	log "github.com/sirupsen/logrus"
)

func main() {

	// setup logger
	// Set the log formatter to output colored logs
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	// Set the output of logs to stdout
	log.SetOutput(os.Stdout)
	// Set the log level to debug
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)

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

	// log each request
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		log.WithFields(log.Fields{
			"time":    start.Format(time.RFC3339Nano),
			"method":  c.Method(),
			"path":    c.Path(),
			"status":  c.Response().StatusCode,
			"latency": latency,
		}).Info("Request")
		return err
	})

	// Enable gzip compression
	app.Use(compress.New())

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

	//Log the start of the application
	log.WithFields(log.Fields{
		"Server running on ": port,
	}).Info("Server running on")

	fmt.Println("Press Ctrl+C to stop the server")

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
