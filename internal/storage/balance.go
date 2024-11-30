package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/config/logger"
	"github.com/learies/gofermart/internal/models"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type BalanceStorage interface {
	GetBalanceByUserID(userID int64) (*models.UserBalance, error)
	CheckBalanceWithdrawal(userID int64, amount float32) error
	GetWithdrawalsByUserID(userID int64) (*[]models.UserWithdrawal, error)
}

type balanceStorage struct {
	db *pgxpool.Pool
}

func NewBalanceStorage(dbPool *pgxpool.Pool) BalanceStorage {
	return &balanceStorage{
		db: dbPool,
	}
}

func (store *balanceStorage) GetBalanceByUserID(userID int64) (*models.UserBalance, error) {
	var userBalance models.UserBalance

	row := store.db.QueryRow(context.Background(),
		"SELECT SUM(accrual) - SUM(withdrawn) AS accrual, SUM(withdrawn) withdrawn FROM orders WHERE user_id = $1", userID)

	row.Scan(&userBalance.Current, &userBalance.Withdraw)

	return &userBalance, nil
}

func (store *balanceStorage) CheckBalanceWithdrawal(userID int64, amount float32) error {

	userBalance, err := store.GetBalanceByUserID(userID)
	if err != nil {
		return err
	}

	if userBalance.Current < 0 {
		logger.Log.Error("Insufficient funds", "balance", userBalance.Current, "withdrawal", amount)
		return ErrInsufficientFunds
	}

	return nil
}

func (store *balanceStorage) GetWithdrawalsByUserID(userID int64) (*[]models.UserWithdrawal, error) {
	var userWithdrawals []models.UserWithdrawal

	rows, err := store.db.Query(context.Background(),
		"SELECT id, withdrawn, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userWithdrawal models.UserWithdrawal
		err := rows.Scan(&userWithdrawal.OrderNumber, &userWithdrawal.Withdrawn, &userWithdrawal.UploadedAt)
		if err != nil {
			return nil, err
		}

		userWithdrawals = append(userWithdrawals, userWithdrawal)
	}

	return &userWithdrawals, nil
}
