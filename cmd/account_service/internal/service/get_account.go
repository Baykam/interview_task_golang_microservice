package accountService

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

// 3. GetAccount: ТЗ 5.1'deki meşhur Cache-Aside mantığı
func (s *service) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	if id == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	s.log.Info("Checking cache for account %s...", id)
	cachedAccount, err := s.accountRepo.Cache.Get(ctx, id)
	if err != nil {
		s.log.Info("Cache error for account %s: %v. Continuing to DB.", id, err)
	}
	if cachedAccount != nil {
		s.log.Info("Cache HIT for account %s", id)
		return cachedAccount, nil
	}

	// Adım 2: Cache'de yoksa DB'den oku (Cache Miss durumu)
	s.log.Info("Cache MISS for account %s. Fetching from DB...", id)
	dbAccount, err := s.accountRepo.DB.GetById(ctx, id)
	if err != nil {
		s.log.Error("Database error fetching account %s: %v", id, err)
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}
	if dbAccount == nil {
		return nil, fmt.Errorf("account not found")
	}

	// Adım 3: DB'den okunan veriyi sonraki istekler için Redis'e yaz
	s.log.Info("Asynchronously populating cache for account %s...", id)
	if err := s.accountRepo.Cache.Add(ctx, dbAccount); err != nil {
		s.log.Info("Failed to write account %s to cache: %v", id, err)
	}

	return dbAccount, nil
}
