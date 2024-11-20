package main

import (
	"github.com/learies/gofermart/internal/app"
	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/config/logger"
)

func main() {
	cfg := config.NewConfig()

	err := logger.NewLogger("info")

	application, err := app.NewApp(cfg)
	if err != nil {
		logger.Log.Error("Could not create app", "error", err)
	}

	if err := application.Run(); err != nil {
		logger.Log.Error("Could not start server", "error", err)
	}
}
