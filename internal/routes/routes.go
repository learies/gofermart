package routes

import (
	"net/http"

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
		r.Post("/register", userHandlers.RegisterUser())
		r.Post("/login", userHandlers.LoginUser())
		r.Post("/orders", userHandlers.CreateOrder(cfg.AccrualSystemAddress))
		r.Get("/orders", userHandlers.GetOrders())
		r.Get("/balance", userHandlers.GetBalance())
		r.Post("/balance/withdraw", userHandlers.Withdraw(cfg.AccrualSystemAddress))
		r.Get("/withdrawals", userHandlers.GetUserWithdrawals())
		r.MethodNotAllowed(methodNotAllowedHandler)
	})

	return nil
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusBadRequest)
}
