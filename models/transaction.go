package models

import "time"

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
)

type Transaction struct {
	Id              *string         `json:"id"`
	AccountId       *string         `json:"account_id"`
	Amount          *float64        `json:"amount"`
	TransactionType TransactionType `json:"type"`
	CreatedAt       *time.Time      `json:"created_at"`
}
