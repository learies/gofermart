package models

type Order struct {
	OrderID int   `db:"order_id" json:"order_id"`
	UserID  int64 `db:"user_id" json:"user_id"`
}
