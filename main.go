package main

import (
	"fmt"
	"log"

	"github.com/mousav1/weiser/bootstrap"
	"github.com/spf13/viper"
)

func main() {
	app, server, err := bootstrap.SetupApp()
	if err != nil {
		log.Fatalf("failed to set up the application: %s", err)
	}

	// You can optionally use app to check if it's correctly set up
	if app == nil {
		log.Fatalf("Failed to initialize the Fiber app")
	}

	fmt.Println("Press Ctrl+C to stop the server")

	go func() {
		if err := server.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("server.port"))); err != nil {
			log.Fatalf("failed to start server: %s", err)
		}
	}()

	// Gracefully shutdown the server
	bootstrap.ShutdownServer(server)
}
