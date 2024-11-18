package repository

import (
	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/storage"
)

type OrderRepository interface {
	CreateOrder(order models.Order) error
	GetOrdersByUserID(userID int) ([]models.Order, error)
}

type orderRepository struct {
	storage storage.OrderStorage
}
