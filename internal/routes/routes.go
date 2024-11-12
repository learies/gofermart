package routes

import (
	"log"

	"github.com/go-chi/chi"

	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/handlers"
	"github.com/learies/gofermart/internal/storage/postgres"
)

type Router struct {
	*chi.Mux
}

func NewRouter() *Router {
	return &Router{Mux: chi.NewRouter()}
}

func (r *Router) Initialize(cfg *config.Config) error {

	dbPool, err := postgres.SetupDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	routes := r.Mux
	userHandlers := handlers.NewHandler(dbPool)

	routes.Route("/api/user", func(r chi.Router) {
		r.Get("/register", userHandlers.RegisterUser)
		r.Post("/login", userHandlers.LoginUser)
	})

	return nil
}
