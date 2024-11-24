package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/learies/gofermart/internal/models"
	"github.com/learies/gofermart/internal/storage"
)

func (h *Handler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		orderNumber := string(body)
		if orderNumber == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		order := h.order.GetOrderByOrderID(orderNumber)

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

		newOrder := models.Order{
			OrderID: orderNumber,
			UserID:  UserID,
		}

		err = h.order.CreateOrder(newOrder)
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

func (h *Handler) GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		UserID, ok := r.Context().Value("userID").(int64)
		if !ok {
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		orders, err := h.order.GetOrdersByUserID(UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(orders) == 0 {
			http.Error(w, "No orders found", http.StatusAccepted)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orders)
	}
}
