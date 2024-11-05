package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupDB sets up the database connection
func SetupDB() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// CloseDB closes the database connection
func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}
