package models

import "time"

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	Id              *string         `json:"id" db:"id"`
	AccountId       *string         `json:"account_id" db:"account_id"`
	ToAccountId     *string         `json:"to_account_id" db:"to_account_id"`
	Amount          *int64          `json:"amount" db:"amount"`
	TransactionType TransactionType `json:"transaction_type" db:"transaction_type"`
	CreatedAt       *time.Time      `json:"created_at" db:"created_at"`
}
