package transActionDb

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (d *db) GetListByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// sqlx etiketleri (db:"column_name") struct üzerinde tanımlıysa direkt SelectContext kullanabiliriz
	query := `
        SELECT id, account_id, to_account_id, amount, transaction_type, created_at
        FROM transactions
        WHERE account_id = $1 OR to_account_id = $2
        ORDER BY created_at DESC`

	err := d.db.SelectContext(ctx, &transactions, query, accountID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to select transactions via sqlx: %w", err)
	}

	return transactions, nil
}
