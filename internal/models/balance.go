package models

type Balance struct {
	Current  float32 `db:"current" json:"current"`
	Withdraw float32 `db:"withdrawn" json:"withdrawn"`
}

type Withdraw struct {
	UserID   int64   `json:"-"`
	OrderID  string  `json:"order"`
	Withdraw float32 `json:"sum"`
}
