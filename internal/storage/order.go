package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
)

type OrderStorage interface {
	CreateOrder(order models.Order) error
	GetOrdersByUserID(userID int64) ([]models.Order, error)
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

func (store *orderStorage) GetOrdersByUserID(userID int64) ([]models.Order, error) {
	rows, err := store.db.Query(context.Background(),
		"SELECT id, user_id FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.OrderID, &order.UserID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
