package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
)

type OrderStorage interface {
	CreateOrder(order models.Order) error
	GetOrderByOrderID(orderID int) (models.Order, error)
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
	row := store.db.QueryRow(context.Background(),
		"INSERT INTO orders (id, user_id) VALUES ($1, $2) RETURNING id",
		order.OrderID, order.UserID)

	var number int
	err := row.Scan(&number)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = ErrConflict
		}
		return err
	}

	return nil
}

func (store *orderStorage) GetOrderByOrderID(orderID int) (models.Order, error) {
	var order models.Order

	row := store.db.QueryRow(context.Background(),
		"SELECT id, user_id FROM orders WHERE id = $1", orderID)

	err := row.Scan(&order.OrderID, &order.UserID)
	if err != nil {
		return order, err
	}

	return order, nil
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
