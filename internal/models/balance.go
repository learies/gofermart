package models

import "time"

type UserBalance struct {
	Current  float32 `json:"current"`
	Withdraw float32 `json:"withdrawn"`
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
