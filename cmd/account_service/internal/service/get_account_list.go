package accountService

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

// 4. GetAccountsList: Tüm aktif hesapları DB'den listeler
func (s *service) GetAccountsList(ctx context.Context) ([]models.Account, error) {
	s.log.Info("Fetching all active accounts from database...")
	accounts, err := s.accountRepo.DB.GetList(ctx)
	if err != nil {
		s.log.Error("Failed to fetch accounts list from DB: %v", err)
		return nil, fmt.Errorf("failed to fetch accounts list: %w", err)
	}

	return accounts, nil
}
