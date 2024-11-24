package models

import "time"

type Order struct {
	OrderID    string    `db:"number" json:"number"`
	Status     string    `db:"status" json:"status"`
	Accrual    float64   `db:"accrual,omitempty" json:"accrual,omitempty"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
	UserID     int64     `db:"user_id" json:"-"`
}
