package routes

import (
	"github.com/go-chi/chi"
)

type Router struct {
	*chi.Mux
}

func (r *Router) Initialize() {
	r.Mux = chi.NewRouter()
}

func NewRouter() *Router {
	return &Router{Mux: chi.NewRouter()}
}
