package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/learies/gofermart/internal/models"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	order := string(body)
	if order == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	OrderID, err := strconv.Atoi(order)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	UserID := r.Context().Value("userID").(int)

	newOrder := models.Order{
		OrderID: OrderID,
		UserID:  UserID,
	}

	if err := h.order.CreateOrder(newOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	UserID := r.Context().Value("userID").(int)

	orders, err := h.order.GetOrdersByUserID(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
