package app

import (
	"log"
	"net/http"

	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/routes"
)

type App struct {
	Router *routes.Router
	Config *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	router := routes.NewRouter()
	if err := router.Initialize(cfg); err != nil {
		return nil, err
	}

	return &App{
		Router: router,
		Config: cfg,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Starting server on %s\n", a.Config.RunAddress)
	return http.ListenAndServe(a.Config.RunAddress, a.Router.Mux)
}
