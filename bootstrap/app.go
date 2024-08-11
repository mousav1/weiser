package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/mousav1/weiser/app/cache"
	kernel "github.com/mousav1/weiser/app/http"
	middleware "github.com/mousav1/weiser/app/http/middlewares"
	"github.com/mousav1/weiser/app/session"
	"github.com/mousav1/weiser/app/storage"
	"github.com/mousav1/weiser/database"
	"github.com/mousav1/weiser/routes"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"

	log "github.com/sirupsen/logrus"
)

// SetupApp initializes and returns a configured Fiber app instance.
func SetupApp() (*fiber.App, *fasthttp.Server, error) {
	// setup logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)

	// set config
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get SQL database handle: %w", err)
	}
	defer sqlDB.Close()

	// Initialize the session manager
	if err := session.InitSessionManager(); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize session manager: %w", err)
	}
	go middleware.StartSessionCleaner()

	// Create a new in-memory cache with a default expiration of 1 minute
	err = cache.InitializeCache(time.Minute, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cache: %w", err)
	}

	// Initialize storage
	defaultDriverName := viper.GetString("storage.default")
	storageConfig := viper.Sub("storage.disks")
	if err := storage.InitializeStorage(defaultDriverName, storageConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Create the Fiber app
	app := fiber.New(fiber.Config{})

	// add middlewares
	for _, myMiddleware := range kernel.Middleware {
		app.Use(myMiddleware)
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
		return nil, nil, fmt.Errorf("failed to set up routes: %w", err)
	}

	// Start the server
	port := viper.GetString("server.port")
	if port == "" {
		port = "3000"
	}

	server := &fasthttp.Server{
		Handler: app.Handler(),
	}

	//Log the start of the application
	log.WithFields(log.Fields{
		"Server running on ": port,
	}).Info("Server running on")

	return app, server, nil
}

// ShutdownServer gracefully shuts down the server.
func ShutdownServer(server *fasthttp.Server) {
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
