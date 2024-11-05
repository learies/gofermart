package app

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/user"
	"github.com/learies/gofermart/pkg/db"
)

type App struct {
	Routes *chi.Mux
	DB     *pgxpool.Pool
}

func NewApp() *App {
	dbPool, err := db.SetupDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	router := chi.NewRouter()
	userHandler := user.NewHandler(dbPool)

	router.Post("/api/user/register", userHandler.RegisterUser)

	return &App{
		Routes: router,
		DB:     dbPool,
	}
}

func (a *App) Close() {
	db.CloseDB(a.DB)
}
