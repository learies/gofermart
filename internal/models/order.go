package models

import "time"

type Order struct {
	OrderID    string    `db:"id" json:"order"`
	Status     string    `db:"status" json:"status"`
	Accrual    float32   `db:"accrual,omitempty" json:"accrual,omitempty"`
	UploadedAt time.Time `db:"uploaded_at,omitempty" json:"uploaded_at,omitempty"`
	UserID     int64     `db:"user_id" json:"-"`
}
