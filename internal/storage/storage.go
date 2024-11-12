package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/services"
)

type Storage interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

type PostgresStorage struct {
	db   *pgxpool.Pool
	auth services.AuthService
}

func NewPostgresStorage(dbPool *pgxpool.Pool) Storage {
	return &PostgresStorage{
		db:   dbPool,
		auth: services.NewAuthService(),
	}
}

func (store *PostgresStorage) CreateUser(user models.User) error {
	hashedPassword, err := store.auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = store.db.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresStorage) GetUserByUsername(username string) (*models.User, error) {
	row := store.db.QueryRow(context.Background(),
		"SELECT username, password FROM users WHERE username=$1", username)

	var user models.User
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
