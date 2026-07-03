package commandService

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (s *service) CreateTransAction(ctx context.Context, data models.Transaction) error {
	s.logger.Info("Creating transaction record in DB", "account_id", data.AccountId, "type", data.TransactionType)

	insertedID, err := s.repos.DB.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("commandService.CreateTransAction failed: %w", err)
	}

	s.logger.Info("Transaction successfully created", "transaction_id", *insertedID)

	return nil
}
