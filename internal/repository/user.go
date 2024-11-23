package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/storage"
)

type UserRepository interface {
	CreateUser(username, password string) (int64, error)
	GetUserByUsername(username string) (*models.User, error)
}

type userRepository struct {
	userStorage storage.UserStorage
}

func NewRepository(dbPool *pgxpool.Pool) UserRepository {
	return &userRepository{
		userStorage: storage.NewPostgresStorage(dbPool),
	}
}

func (repo *userRepository) CreateUser(username, password string) (int64, error) {
	return repo.userStorage.CreateUser(username, password)
}

func (repo *userRepository) GetUserByUsername(username string) (*models.User, error) {
	return repo.userStorage.GetUserByUsername(username)
}
