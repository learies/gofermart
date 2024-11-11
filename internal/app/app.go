package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/handlers"
	"github.com/learies/gofermart/internal/routes"
)

type App struct {
	Router *routes.Router
}

func NewApp() *App {
	return &App{
		Router: routes.NewRouter(),
	}
}

func (a *App) Run(cfg *config.Config) {

	r := a.Router
	userHandlers := handlers.NewHandlers()

	r.Route("/", func(r chi.Router) {
		r.Get("/", userHandlers.RegisterUser)
	})

	log.Printf("Starting server on %s\n", cfg.RunAddress)
	if err := http.ListenAndServe(cfg.RunAddress, r.Mux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
