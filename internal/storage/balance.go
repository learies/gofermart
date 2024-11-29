package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type BalanceStorage interface {
	GetBalanceByUserID(userID int64) (*models.Balance, error)
	CheckBalanceWithdrawal(userID int64, amount float32) error
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
		"SELECT SUM(accrual) - SUM(withdrawn) AS accrual, SUM(withdrawn) withdrawn FROM orders WHERE user_id = $1", userID)

	row.Scan(&balance.Current, &balance.Withdraw)

	return &balance, nil
}

func (store *balanceStorage) CheckBalanceWithdrawal(userID int64, amount float32) error {

	balance, err := store.GetBalanceByUserID(userID)
	if err != nil {
		return err
	}

	if balance.Current < amount {
		return ErrInsufficientFunds
	}

	return nil
}
