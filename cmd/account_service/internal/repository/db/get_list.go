package accountRepositoryDb

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (r *repository) GetList(ctx context.Context) ([]models.Account, error) {
	accounts := make([]models.Account, 0)

	query := `
		SELECT id, balance, currency, is_locked, created_at, deleted_at 
		FROM accounts 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &accounts, query)
	if err != nil {
		r.log.Error("failed to fetch accounts list: %v", err)
		return nil, fmt.Errorf("failed to get accounts list: %w", err)
	}

	return accounts, nil
}
