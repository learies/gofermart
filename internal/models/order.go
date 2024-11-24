package models

type Order struct {
	OrderID string `db:"number" json:"number"`
	UserID  int64  `db:"user_id" json:"-"`
}
