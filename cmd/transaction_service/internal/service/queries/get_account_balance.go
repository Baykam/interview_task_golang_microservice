package querieService

import (
	"context"
	accountProto "interview_task_golang_microservices/protos"
)

func (s *service) GetAccountBalance(ctx context.Context, in *accountProto.GetAccountBalanceRequest) (*accountProto.GetAccountBalanceResponse, error) {
	return &accountProto.GetAccountBalanceResponse{}, nil
}
