package models

import "time"

type Account struct {
	ID        *string    `json:"id"`
	UserID    *int64     `json:"user_id"`
	Balance   *float64   `json:"balance"`
	Currency  *string    `json:"currency"`
	IsLocked  *bool      `json:"is_locked"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
