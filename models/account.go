package models

import "time"

type Account struct {
	ID        *string    `json:"id" db:"id"`
	UserID    *string    `json:"user_id" db:"user_id"`
	Balance   *int64     `json:"balance" db:"balance"`
	Currency  *string    `json:"currency" db:"currency"`
	IsLocked  *bool      `json:"is_locked" db:"is_locked"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
