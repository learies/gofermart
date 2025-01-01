package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/services"
)

var ErrConflict = errors.New("data conflict")

type UserStorage interface {
	CreateUser(username, password string) (int64, error)
	GetUserByUsername(username string) (*models.User, error)
}

type userStorage struct {
	db   *pgxpool.Pool
	auth services.AuthService
}

func NewPostgresStorage(dbPool *pgxpool.Pool) UserStorage {
	return &userStorage{
		db:   dbPool,
		auth: services.NewAuthService(),
	}
}

func (store *userStorage) CreateUser(username, password string) (int64, error) {

	hashedPassword, err := store.auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	row := store.db.QueryRow(context.Background(),
		"INSERT INTO users (username, password) VALUES($1, $2) RETURNING id",
		username, hashedPassword)

	var userID int64
	err = row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = ErrConflict
		}
		return userID, err
	}

	return userID, nil
}

func (store *userStorage) GetUserByUsername(username string) (*models.User, error) {
	row := store.db.QueryRow(context.Background(),
		"SELECT id, username, password FROM users WHERE username=$1", username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
