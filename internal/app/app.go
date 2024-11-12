package app

import (
	"log"
	"net/http"

	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/routes"
)

type App struct {
	Router *routes.Router
}

func NewApp(cfg *config.Config) *App {
	app := &App{
		Router: routes.NewRouter(),
	}
	app.Router.Initialize(cfg)
	return app
}

func (a *App) Run(cfg *config.Config) {
	log.Printf("Starting server on %s\n", cfg.RunAddress)
	if err := http.ListenAndServe(cfg.RunAddress, a.Router.Mux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
