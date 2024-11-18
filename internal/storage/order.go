package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
)

type OrderStorage interface {
	CreateOrder(order models.Order) error
}

type orderStorage struct {
	db *pgxpool.Pool
}

func NewOrderStorage(dbPool *pgxpool.Pool) OrderStorage {
	return &orderStorage{
		db: dbPool,
	}
}

func (store *orderStorage) CreateOrder(order models.Order) error {
	_, err := store.db.Exec(context.Background(),
		"INSERT INTO orders (id, user_id) VALUES ($1, $2)",
		order.OrderID, order.UserID)
	if err != nil {
		return err
	}

	return nil
}
