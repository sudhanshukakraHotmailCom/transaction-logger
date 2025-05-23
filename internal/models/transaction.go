package models

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	SenderAccount   string    `json:"sender_account"`
	ReceiverAccount string    `json:"receiver_account"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	TransactionType string    `json:"transaction_type"`
	Status          string    `json:"status"`
	UserID          string    `json:"user_id"`
}

type CreateTransactionRequest struct {
	SenderAccount   string  `json:"sender_account" validate:"required"`
	ReceiverAccount string  `json:"receiver_account" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Currency        string  `json:"currency" validate:"required,oneof=USD EUR GBP"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=Transfer Deposit Withdrawal"`
	UserID          string  `json:"-"` // Not exposed in JSON, used internally
}
