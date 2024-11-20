package routes

import (
	"github.com/go-chi/chi"

	"github.com/learies/gofermart/internal/config"
	"github.com/learies/gofermart/internal/config/logger"
	"github.com/learies/gofermart/internal/handlers"
	internalMiddleware "github.com/learies/gofermart/internal/middleware"
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
		logger.Log.Error("Unable to connect to database", "error", err)
	}

	routes := r.Mux
	routes.Use(internalMiddleware.JWTMiddleware)
	routes.Use(internalMiddleware.WithLogging)

	userHandlers := handlers.NewHandler(dbPool)

	routes.Route("/api/user", func(r chi.Router) {
		r.Post("/register", userHandlers.RegisterUser)
		r.Post("/login", userHandlers.LoginUser)
		r.Post("/orders", userHandlers.CreateOrder)
		r.Get("/orders", userHandlers.GetOrders)
	})

	return nil
}
