package accountService

import (
	"context"
	"fmt"
)

// 2. DeleteAccount: ТЗ 5.1'e göre hem DB'den soft delete yapar hem de Redis'ten siler (Cache Invalidation)
func (s *service) DeleteAccount(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	// Önce DB'de hesabın var olup olmadığını kontrol etmek iyi bir pratiktir
	account, err := s.accountRepo.DB.GetById(ctx, id)
	if err != nil {
		s.log.Error("Failed to check account before deletion: %v", err)
		return fmt.Errorf("failed to delete account: %w", err)
	}
	if account == nil {
		return fmt.Errorf("account not found")
	}

	// DB'de yumuşak silme (Soft Delete) yapıyoruz
	s.log.Info("Soft deleting account %s from database...", id)
	err = s.accountRepo.DB.Delete(ctx, account)
	if err != nil {
		s.log.Error("Database failed to delete account %s: %v", id, err)
		return fmt.Errorf("failed to delete account from DB: %w", err)
	}

	// DB silme işlemi başarılıysa, eski veri kalmasın diye Redis'ten tamamen uçuruyoruz
	s.log.Info("Invalidating cache for account %s...", id)
	if err := s.accountRepo.Cache.Delete(ctx, id); err != nil {
		s.log.Info("Failed to remove account %s from cache: %v", id, err)
	}

	return nil
}
