package querieService

import (
	"context"
	"fmt"
	accountProto "interview_task_golang_microservices/protos"
)

func (s *service) GetTransactionHistory(ctx context.Context, in *accountProto.GetTransactionHistoryRequest) (*accountProto.GetTransactionHistoryResponse, error) {
	s.logger.Info("Fetching transaction history for account", "account_id", in.AccountId)

	// 1. Daha önce yazdığımız DB fonksiyonunu çağırıyoruz
	dbTransactions, err := s.repos.DB.GetListByAccountID(ctx, in.AccountId)
	if err != nil {
		s.logger.Error("Failed to fetch transaction history from DB", "error", err)
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	// 2. DB'den gelen []models.Transaction slice'ını []*accountProto.TransactionMessage slice'ına mapliyoruz
	var protoTransactions []*accountProto.TransactionMessage
	for _, tx := range dbTransactions {
		// Pointer olan to_account_id alanını güvenli bir şekilde string'e çeviriyoruz (nil kontrolü)
		toAccountStr := ""
		if tx.ToAccountId != nil {
			toAccountStr = *tx.ToAccountId
		}

		protoTx := &accountProto.TransactionMessage{
			Id:              *tx.Id,
			AccountId:       *tx.AccountId,
			ToAccountId:     toAccountStr,
			Amount:          *tx.Amount,
			TransactionType: string(tx.TransactionType),                 // Özel enum tipini string'e cast ediyoruz
			CreatedAt:       tx.CreatedAt.Format("2006-01-02 15:04:05"), // Standart tarih formatı
		}
		protoTransactions = append(protoTransactions, protoTx)
	}

	// 3. Proto formatındaki response'u dönüyoruz
	return &accountProto.GetTransactionHistoryResponse{
		Transactions: protoTransactions,
	}, nil
}
