package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUsersTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	)`)

	return err
}

func CreateOrdersTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(255) PRIMARY KEY,
		status VARCHAR(10) DEFAULT 'NEW' CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
		accrual FLOAT DEFAULT 0.0,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER NOT NULL REFERENCES users(id)
	)`)

	return err
}

func SetupDB() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URI")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = CreateUsersTable(pool)
	if err != nil {
		return nil, err
	}

	err = CreateOrdersTable(pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// CloseDB closes the database connection
func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}
