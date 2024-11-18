package models

type Order struct {
	OrderID int `db:"order_id" json:"order_id"`
	UserID  int `db:"user_id" json:"user_id"`
}
