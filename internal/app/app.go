package app

import (
	"net/http"

	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/config/logger"
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
	logger.Log.Info("Starting server", "address", a.Config.RunAddress)
	return http.ListenAndServe(a.Config.RunAddress, a.Router.Mux)
}
