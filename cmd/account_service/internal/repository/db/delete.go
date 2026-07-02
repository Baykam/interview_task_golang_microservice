package accountRepositoryDb

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
	"time"
)

func (r *repository) Delete(ctx context.Context, account *models.Account) error {
	if account == nil || account.ID == nil {
		return fmt.Errorf("account or account ID cannot be nil")
	}

	now := time.Now()
	account.DeletedAt = &now

	query := `
		UPDATE accounts 
		SET deleted_at = :deleted_at 
		WHERE id = :id AND deleted_at IS NULL
	`

	result, err := r.db.NamedExecContext(ctx, query, account)
	if err != nil {
		r.log.Error("failed to soft delete account with id %s: %v", *account.ID, err)
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found or already deleted")
	}

	return nil
}
