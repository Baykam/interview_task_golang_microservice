package accountService

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (s *service) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	if account == nil {
		return nil, fmt.Errorf("account data cannot be nil")
	}
	if account.Currency == nil || *account.Currency == "" {
		return nil, fmt.Errorf("currency is required to create an account")
	}

	s.log.Info("Creating account in database for user: %v", account.UserID)
	newAccount, err := s.accountRepo.DB.Create(ctx, account)
	if err != nil {
		s.log.Error("Failed to create account in database: %v", err)
		return nil, fmt.Errorf("failed to create account in database: %w", err)
	}

	s.log.Info("Caching newly created account with ID: %s", "accountId", *newAccount.ID)
	err = s.accountRepo.Cache.Add(ctx, newAccount)
	if err != nil {
		s.log.Info("Account created in DB but failed to cache with ID %s: %v", *account.ID, err)
	}

	return newAccount, nil
}
