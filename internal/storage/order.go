package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/config/logger"
	"github.com/learies/gofermart/internal/models"
)

type OrderStorage interface {
	CreateOrder(order models.Order) error
	GetOrderByOrderID(orderID string) models.Order
	GetOrdersByUserID(userID int64) ([]models.OrderResponse, error)
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
	if store.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	row := store.db.QueryRow(context.Background(),
		"INSERT INTO orders (id, status, accrual, withdrawn, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		order.OrderID, order.Status, order.Accrual, order.Withdrawn, order.UserID)

	var number string
	err := row.Scan(&number)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = ErrConflict
		}
		return err
	}
	logger.Log.Info("Order created successfully", "order", order)
	return nil
}

func (store *orderStorage) GetOrderByOrderID(orderID string) models.Order {
	var order models.Order

	row := store.db.QueryRow(context.Background(),
		"SELECT id, user_id FROM orders WHERE id = $1", orderID)

	err := row.Scan(&order.OrderID, &order.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoData {
			return models.Order{}
		}
		return models.Order{}
	}

	return order
}

func (store *orderStorage) GetOrdersByUserID(userID int64) ([]models.OrderResponse, error) {
	rows, err := store.db.Query(context.Background(),
		"SELECT id, status, accrual, withdrawn, uploaded_at, user_id FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderResponse
	for rows.Next() {
		var order models.OrderResponse
		if err := rows.Scan(&order.OrderID, &order.Status, &order.Accrual, &order.Withdrawn, &order.UploadedAt, &order.UserID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error while scanning rows", "error", err)
		return nil, err
	}

	return orders, nil
}
