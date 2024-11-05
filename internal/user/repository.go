package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(user User) error
	GetUserByUsername(username string) (*User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(dbPool *pgxpool.Pool) Repository {
	return &PostgresRepository{db: dbPool}
}

func (repo *PostgresRepository) CreateUser(user User) error {
	_, err := repo.db.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) GetUserByUsername(username string) (*User, error) {
	row := repo.db.QueryRow(context.Background(),
		"SELECT username, password FROM users WHERE username=$1", username)

	var user User
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
