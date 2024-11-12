package main

import (
	"github.com/learies/gofermart/internal/app"
	"github.com/learies/gofermart/internal/config"
)

func main() {
	cfg := config.NewConfig()

	application := app.NewApp(cfg)
	application.Run(cfg)
}
