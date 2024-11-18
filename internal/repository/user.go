package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/storage"
)

type Repository interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

type repository struct {
	storage storage.Storage
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{
		storage: storage.NewPostgresStorage(dbPool),
	}
}

func (repo *repository) CreateUser(user models.User) error {
	return repo.storage.CreateUser(user)
}

func (repo *repository) GetUserByUsername(username string) (*models.User, error) {
	return repo.storage.GetUserByUsername(username)
}
