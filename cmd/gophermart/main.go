package main

import (
	"log"

	"github.com/learies/gofermart/internal/app"
	"github.com/learies/gofermart/internal/config"
)

func main() {
	cfg := config.NewConfig()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Could not create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
