package accountRepositoryDb

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
	"time"
)

func (r *repository) Create(ctx context.Context, account *models.Account) error {

	now := time.Now()
	if account.CreatedAt == nil {
		account.CreatedAt = &now
	}

	if account.Balance == nil {
		defaultBalance := int64(0)
		account.Balance = &defaultBalance
	}

	if account.IsLocked == nil {
		defaultLocked := false
		account.IsLocked = &defaultLocked
	}

	query := `
		INSERT INTO accounts (id, user_id, balance, currency, is_locked, created_at, deleted_at)
		VALUES (:id, :user_id, :balance, :currency, :is_locked, :created_at, :deleted_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, account)
	if err != nil {
		r.log.Error("failed to insert account for user %v: %v", account.UserID, err)
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}
