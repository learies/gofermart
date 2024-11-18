package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/storage"
)

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (*models.User, error)
}

type userRepository struct {
	storage storage.UserStorage
}

func NewRepository(dbPool *pgxpool.Pool) UserRepository {
	return &userRepository{
		storage: storage.NewPostgresStorage(dbPool),
	}
}

func (repo *userRepository) CreateUser(user models.User) error {
	return repo.storage.CreateUser(user)
}

func (repo *userRepository) GetUserByUsername(username string) (*models.User, error) {
	return repo.storage.GetUserByUsername(username)
}
