package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/services"
	"github.com/learies/gofermart/internal/storage"
)

type Handler struct {
	user  storage.UserStorage
	auth  services.AuthService
	jwt   services.JWTService
	order storage.OrderStorage
}

func NewHandler(dbPool *pgxpool.Pool) *Handler {
	return &Handler{
		user:  storage.NewPostgresStorage(dbPool),
		auth:  services.NewAuthService(),
		jwt:   services.NewJWTService(),
		order: storage.NewOrderStorage(dbPool),
	}
}
