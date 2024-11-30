package models

import "time"

type Order struct {
	OrderID   string  `db:"id" json:"order"`
	Status    string  `db:"status" json:"status" default:"NEW"`
	Accrual   float32 `db:"accrual" json:"accrual,omitempty"`
	Withdrawn float32 `db:"withdrawn" json:"sum,omitempty"`
	UserID    int64   `db:"user_id" json:"user_id"`
}

type OrderResponse struct {
	OrderID    string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	Withdrawn  float32   `json:"sum,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
