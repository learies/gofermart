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
	GetOrder(orderID string) *models.Order
	GetUserOrders(ctx context.Context, userID int64) (*[]models.OrderResponse, error)
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

func (store *orderStorage) GetOrder(orderID string) *models.Order {
	var order models.Order

	row := store.db.QueryRow(context.Background(),
		"SELECT id, user_id FROM orders WHERE id = $1", orderID)

	err := row.Scan(&order.OrderID, &order.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoData {
			logger.Log.Warn("Order not found", "orderID", orderID)
			return &order
		}
		logger.Log.Error("Error while scanning row", "error", err)
		return &order
	}

	return &order
}

func (store *orderStorage) GetUserOrders(ctx context.Context, userID int64) (*[]models.OrderResponse, error) {
	rows, err := store.db.Query(ctx,
		"SELECT id, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderResponse
	for rows.Next() {
		var order models.OrderResponse
		if err := rows.Scan(&order.OrderID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			logger.Log.Error("Error while scanning row", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	logger.Log.Info("User orders retrieved successfully", "orders", orders)

	if err := rows.Err(); err != nil {
		logger.Log.Error("GetUserOrders: Error while scanning rows", "error", err)
		return nil, err
	}

	return &orders, nil
}
