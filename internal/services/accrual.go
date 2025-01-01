package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/learies/gofermart/internal/config/logger"
	"github.com/learies/gofermart/internal/models"
)

var ErrorStatusTooManyRequests = errors.New("no more than N requests per minute allowed")
var ErrorOrderNotFound = errors.New("order not found")

type AccrualService struct{}

func NewAccrualService() AccrualService {
	return AccrualService{}
}

func (s *AccrualService) FetchOrder(AccrualSystemAddress, orderNumber string) (chan models.Order, chan error) {
	orderChan := make(chan models.Order)
	errChan := make(chan error)

	go func(orderNumber string) {
		defer close(orderChan)
		defer close(errChan)
		newOrder, err := fetchAccrualInfo(AccrualSystemAddress, orderNumber)
		if err != nil {
			errChan <- err
			return
		}
		orderChan <- newOrder
	}(orderNumber)

	return orderChan, errChan
}

func fetchAccrualInfo(AccrualSystemAddress, orderNumber string) (models.Order, error) {
	var order models.Order

	url := fmt.Sprintf("%s/api/orders/%s", AccrualSystemAddress, orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Error("Failed to fetch order", "error", err)
		return order, fmt.Errorf("failed to fetch order: %w", err)
	}
	defer resp.Body.Close()

	logger.Log.Info("Fetched order", "url", url)

	if resp.StatusCode == http.StatusNoContent {
		logger.Log.Info("Order not found")
		return order, ErrorOrderNotFound
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		logger.Log.Info("No more than N requests per minute allowed")
		return order, ErrorStatusTooManyRequests
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Error("Unexpected status code", "status", resp.StatusCode)
		return order, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Failed to read response body", "error", err)
		return order, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &order); err != nil {
		logger.Log.Error("Failed to unmarshal order", "error", err)
		return order, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	logger.Log.Info("Unmarshaled order", "order", order)
	return order, nil
}
