package transActionDb

import (
	"context"
	"errors"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (d *db) Create(ctx context.Context, data models.Transaction) (*string, error) {
	if data.TransactionType != models.TransactionTypeDeposit &&
		data.TransactionType != models.TransactionTypeWithdraw &&
		data.TransactionType != models.TransactionTypeTransfer {
		return nil, fmt.Errorf("transaction type '%s' is not supported", data.TransactionType)
	}

	if d.db == nil {
		d.logger.Error("KRİTİK HATA: sqlx.DB nesnesi NIL!")
		return nil, errors.New("database connection is nil")
	}

	d.logger.Info("Transaction başlatılıyor...")

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	switch data.TransactionType {
	case models.TransactionTypeDeposit:
		query := "UPDATE accounts SET balance = balance + $1 WHERE id = $2"
		_, err = tx.ExecContext(ctx, query, data.Amount, data.AccountId)
		if err != nil {
			return nil, fmt.Errorf("failed to update balance for deposit: %w", err)
		}

	case models.TransactionTypeWithdraw:
		query := "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1"
		res, err := tx.ExecContext(ctx, query, data.Amount, data.AccountId)
		if err != nil {
			return nil, fmt.Errorf("failed to update balance for withdraw: %w", err)
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return nil, errors.New("insufficient balance")
		}

	case models.TransactionTypeTransfer:
		if data.ToAccountId == nil || *data.ToAccountId == "" {
			return nil, errors.New("target account ID (to_account_id) is required for transfer")
		}

		deductQuery := "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1"
		res, err := tx.ExecContext(ctx, deductQuery, data.Amount, data.AccountId)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct balance from sender: %w", err)
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return nil, errors.New("sender has insufficient balance")
		}

		addQuery := "UPDATE accounts SET balance = balance + $1 WHERE id = $2"
		_, err = tx.ExecContext(ctx, addQuery, data.Amount, *data.ToAccountId)
		if err != nil {
			return nil, fmt.Errorf("failed to add balance to receiver: %w", err)
		}
	}

	var insertedID string
	insertQuery := `
		INSERT INTO transactions (account_id, to_account_id, amount, transaction_type, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id`

	err = tx.QueryRowContext(ctx, insertQuery,
		data.AccountId,
		data.ToAccountId,
		data.Amount,
		data.TransactionType,
	).Scan(&insertedID)

	if err != nil {
		return nil, fmt.Errorf("failed to insert transaction log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &insertedID, nil
}
