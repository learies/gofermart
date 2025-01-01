package models

type User struct {
	ID       int64  `db:"id" json:"id" `
	Username string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}
