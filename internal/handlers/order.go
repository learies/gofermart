package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/learies/gofermart/internal/config/logger"
	"github.com/learies/gofermart/internal/constants"
	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/services"
	"github.com/learies/gofermart/internal/storage"
)

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

		if !services.ValidateOrderNumber(orderNumber) {
			logger.Log.Error("Invalid order number", "order", orderNumber)
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		UserID, ok := r.Context().Value(constants.UserIDKey).(int64)
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

		orderChan, errChan := h.accrual.FetchOrder(AccrualSystemAddress, orderNumber)

		orderInfo := models.Order{}

		select {
		case newOrder := <-orderChan:
			orderInfo = newOrder
			orderInfo.OrderID = orderNumber
			orderInfo.UserID = UserID
		case err := <-errChan:
			if errors.Is(err, services.ErrorStatusTooManyRequests) {
				http.Error(w, "No more than N requests per minute allowed", http.StatusTooManyRequests)
				return
			}
			if errors.Is(err, services.ErrorOrderNotFound) {
				orderInfo.OrderID = orderNumber
				orderInfo.UserID = UserID
				orderInfo.Status = "NEW"
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

		UserID, ok := ctx.Value(constants.UserIDKey).(int64)
		if !ok {
			logger.Log.Error("User is not authenticated")
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		userOrders, err := h.order.GetUserOrders(ctx, UserID)
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

		if !services.ValidateOrderNumber(withdraw.OrderNumber) {
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
		UserID, ok := r.Context().Value(constants.UserIDKey).(int64)
		if !ok {
			logger.Log.Error("User is not authenticated")
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		// Получение информации о заказе в отдельной горутине
		orderChan, errChan := h.accrual.FetchOrder(AccrualSystemAddress, withdraw.OrderNumber)

		orderInfo := models.Order{}

		select {
		case newOrder := <-orderChan:
			orderInfo = newOrder
			orderInfo.OrderID = withdraw.OrderNumber
			orderInfo.Withdrawn = withdraw.SumWithdrawn
			orderInfo.UserID = UserID
		case err := <-errChan:
			if errors.Is(err, services.ErrorStatusTooManyRequests) {
				http.Error(w, "No more than N requests per minute allowed", http.StatusTooManyRequests)
				return
			}
			if errors.Is(err, services.ErrorOrderNotFound) {
				orderInfo.OrderID = withdraw.OrderNumber
				orderInfo.UserID = UserID
				orderInfo.Withdrawn = withdraw.SumWithdrawn
				orderInfo.Status = "NEW"
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
