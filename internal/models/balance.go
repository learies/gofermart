package models

type Balance struct {
	UserID  int64  `db:"user_id" json:"-"`
	Current uint32 `db:"current" json:"current"`
}
