package accountRepositoryDb

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (r *repository) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
	if account.Balance == nil {
		defaultBalance := int64(0)
		account.Balance = &defaultBalance
	}

	if account.IsLocked == nil {
		defaultLocked := false
		account.IsLocked = &defaultLocked
	}

	query := `
        INSERT INTO accounts (user_id, balance, currency, is_locked)
        VALUES (:user_id, :balance, :currency, :is_locked)
        RETURNING id, user_id, balance, currency, is_locked, created_at
    `

	rows, err := r.db.NamedQueryContext(ctx, query, account)
	if err != nil {
		r.log.Error("failed to insert account for user %v: %v", account.UserID, err)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	defer rows.Close()

	var insertedAccount models.Account

	if rows.Next() {
		if err := rows.StructScan(&insertedAccount); err != nil {
			r.log.Error("failed to scan inserted account for user %v: %v", account.UserID, err)
			return nil, fmt.Errorf("failed to scan inserted account: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &insertedAccount, nil
}
