package accountService

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/models"
)

func (s *service) CreateAccount(ctx context.Context, account *models.Account) error {
	// 0. Temel Validasyon (Mülakatta her zaman artı puandır)
	if account == nil {
		return fmt.Errorf("account data cannot be nil")
	}
	if account.Currency == nil || *account.Currency == "" {
		return fmt.Errorf("currency is required to create an account")
	}

	// 1. ADIM: Veri Tabanına Kaydet (Write to DB)
	// Önce veri tabanında kaydın oluşmasını sağlıyoruz. ID ve CreatedAt orada veya repo katmanında set edilecek.
	s.log.Info("Creating account in database for user: %v", account.UserID)
	err := s.accountRepo.DB.Create(ctx, account)
	if err != nil {
		s.log.Error("Failed to create account in database: %v", err)
		return fmt.Errorf("failed to create account in database: %w", err)
	}

	// 2. ADIM: Önbelleğe Yaz (Write to Cache)
	// DB işlemi başarılı olduysa, kullanıcı hemen arkasından GET isteği attığında
	// DB'yi yormasın diye veriyi anında Redis'e de yazıyoruz (Write-Through).
	s.log.Info("Caching newly created account with ID: %s", *account.ID)
	err = s.accountRepo.Cache.Add(ctx, account)
	if err != nil {
		// MÜLAKAT İÇİN KRİTİK NOT: Cache'e yazma hatası tüm akışı iptal etmemeli (Rollback yapılmamalı).
		// Çünkü veri zaten DB'ye başarıyla yazıldı. Cache başarısız olsa bile sistem çalışmaya devam edebilir.
		s.log.Info("Account created in DB but failed to cache with ID %s: %v", *account.ID, err)
	}

	return nil
}
