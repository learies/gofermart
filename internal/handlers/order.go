package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/learies/gofermart/internal/config/logger"
	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/storage"
)

var ErrorStatusTooManyRequests = errors.New("no more than N requests per minute allowed")

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
		return order, nil
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

func ValidateOrderNumber(orderNumber string) bool {
	var sum int
	alternate := false

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil || !unicode.IsDigit(rune(orderNumber[i])) {
			return false
		}
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

func (h *Handler) CreateOrder(AccrualSystemAddress string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		orderNumber := strings.TrimSpace(string(body))
		if orderNumber == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		order := h.order.GetOrder(orderNumber)

		if !ValidateOrderNumber(orderNumber) {
			logger.Log.Error("Invalid order number", "order", orderNumber)
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		UserID, ok := r.Context().Value("userID").(int64)
		if !ok {
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		if order.OrderID != "" {
			if order.UserID != UserID {
				http.Error(w, "Order does not belong to user", http.StatusConflict)
				return
			}
		}

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

		orderInfo := models.Order{}

		select {
		case newOrder := <-orderChan:
			orderInfo = newOrder
			orderInfo.OrderID = orderNumber
			orderInfo.Status = "NEW"
			orderInfo.UserID = UserID
		case err := <-errChan:
			if errors.Is(err, ErrorStatusTooManyRequests) {
				http.Error(w, "No more than N requests per minute allowed", http.StatusTooManyRequests)
				return
			}
		}

		err = h.order.CreateOrder(orderInfo)
		if err != nil {
			if errors.Is(err, storage.ErrConflict) {
				http.Error(w, "We already have that order", http.StatusOK)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("New order number has been accepted for processing"))
	}
}

func (h *Handler) GetUserOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		UserID, ok := ctx.Value("userID").(int64)
		if !ok {
			logger.Log.Error("User is not authenticated")
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		userOrders, err := h.order.GetUserOrders(UserID)
		if err != nil {
			logger.Log.Error("Failed to get user orders", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(*userOrders) == 0 {
			logger.Log.Info("No user orders found")
			http.Error(w, "No user orders found", http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*userOrders)
	}
}

func (h *Handler) Withdraw(AccrualSystemAddress string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var withdraw models.WithdrawRequest

		if !ValidateOrderNumber(withdraw.OrderNumber) {
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		// Декодирование тела запроса
		if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
			logger.Log.Error("Failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверка аутентификации пользователя
		UserID, ok := r.Context().Value("userID").(int64)
		if !ok {
			logger.Log.Error("User is not authenticated")
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		// Получение информации о заказе в отдельной горутине
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
		}(withdraw.OrderNumber)

		orderInfo := models.Order{}

		select {
		case newOrder := <-orderChan:
			orderInfo = newOrder
			orderInfo.OrderID = withdraw.OrderNumber
			orderInfo.Withdrawn = withdraw.SumWithdrawn
			orderInfo.Status = "NEW"
			orderInfo.UserID = UserID
		case err := <-errChan:
			if errors.Is(err, ErrorStatusTooManyRequests) {
				http.Error(w, "No more than N requests per minute allowed", http.StatusTooManyRequests)
				return
			}
		}

		// Создание нового заказа
		if err := h.order.CreateOrder(orderInfo); err != nil {
			if errors.Is(err, storage.ErrConflict) {
				http.Error(w, "We already have that order", http.StatusOK)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Проверка возможности снятия средств с баланса
		if err := h.balance.CheckBalanceWithdrawal(UserID, withdraw.SumWithdrawn); err != nil {
			if errors.Is(err, storage.ErrInsufficientFunds) {
				http.Error(w, "Withdrawal amount exceeds the order accrual", http.StatusPaymentRequired)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Установка заголовков и ответ клиенту
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
