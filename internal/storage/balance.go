package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
)

type BalanceStorage interface {
	GetBalanceByUserID(userID int64) (*models.Balance, error)
}

type balanceStorage struct {
	db *pgxpool.Pool
}

func NewBalanceStorage(dbPool *pgxpool.Pool) BalanceStorage {
	return &balanceStorage{
		db: dbPool,
	}
}

func (store *balanceStorage) GetBalanceByUserID(userID int64) (*models.Balance, error) {
	var balance models.Balance

	row := store.db.QueryRow(context.Background(),
		"SELECT user_id, current FROM balance WHERE user_id = $1", userID)

	row.Scan(&balance.UserID, &balance.Current)

	return &balance, nil
}
