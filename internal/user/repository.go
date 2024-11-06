package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents a user in the system
type repository interface {
	createUser(user user) error
	getUserByUsername(username string) (*user, error)
}

// PostgresRepository is a repository for user data in a PostgreSQL database
type postgresRepository struct {
	db      *pgxpool.Pool
	service service
}

// NewPostgresRepository creates a new PostgresRepository instance
func newPostgresRepository(dbPool *pgxpool.Pool) repository {
	return &postgresRepository{
		db:      dbPool,
		service: newUserService(),
	}
}

// CreateUser creates a new user in the database
func (repo *postgresRepository) createUser(user user) error {
	hashedPassword, err := repo.service.hashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByUsername retrieves a user from the database by username
func (repo *postgresRepository) getUserByUsername(username string) (*user, error) {
	row := repo.db.QueryRow(context.Background(),
		"SELECT username, password FROM users WHERE username=$1", username)

	var user user
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
