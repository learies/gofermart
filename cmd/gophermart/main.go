package main

import (
	"log"
	"net/http"

	"github.com/learies/gofermart/internal/app"
	"github.com/learies/gofermart/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize the application
	application := app.NewApp()
	defer application.Close()

	// Start the HTTP server
	log.Printf("Starting server on %s\n", cfg.RunAddress)
	if err := http.ListenAndServe(cfg.RunAddress, application.Routes); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
