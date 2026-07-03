package transActionDb

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (d *db) GetListByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `
		SELECT id, account_id, to_account_id, amount, transaction_type, created_at
		FROM transactions
		WHERE account_id = $1 OR to_account_id = $2
		ORDER BY created_at DESC`

	rows, err := d.db.QueryContext(ctx, query, accountID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tx models.Transaction

		err := rows.Scan(
			&tx.Id,
			&tx.AccountId,
			&tx.ToAccountId,
			&tx.Amount,
			&tx.TransactionType,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}

		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return transactions, nil
}
