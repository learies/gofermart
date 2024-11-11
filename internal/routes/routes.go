package routes

import (
	"github.com/go-chi/chi"
	"github.com/learies/gofermart/internal/handlers"
)

type Router struct {
	*chi.Mux
}

func NewRouter() *Router {
	return &Router{Mux: chi.NewRouter()}
}

func (r *Router) Initialize() {
	routes := r.Mux
	userHandlers := handlers.NewHandlers()

	routes.Route("/", func(r chi.Router) {
		r.Get("/", userHandlers.RegisterUser)
	})
}
