package handlers

import "net/http"

type Handler struct {
}

func NewHandlers() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
