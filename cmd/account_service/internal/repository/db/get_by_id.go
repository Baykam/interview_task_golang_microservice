package accountRepositoryDb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (r *repository) GetById(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	query := `
		SELECT id, user_id, balance, currency, is_locked, created_at, deleted_at 
		FROM accounts 
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Info("account not found with id: %s", id)
			return nil, nil
		}

		r.log.Error("failed to get account by id %s: %v", id, err)
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	return &account, nil
}
