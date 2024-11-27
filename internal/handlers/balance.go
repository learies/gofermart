package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID, ok := r.Context().Value("userID").(int64)
		if !ok {
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		balance, err := h.balance.GetBalanceByUserID(UserID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(balance)
	}
}