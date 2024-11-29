package models

import "time"

type Balance struct {
	Current  float32 `db:"current" json:"current"`
	Withdraw float32 `db:"withdrawn" json:"withdrawn"`
}

type WithdrawRequest struct {
	OrderNumber  string  `json:"order"`
	SumWithdrawn float32 `json:"sum"`
}

type WithdrawalsResponse struct {
	OrderNumber string    `json:"order"`
	Withdrawn   float32   `json:"sum"`
	UploadedAt  time.Time `json:"processed_at"`
}
