package querieService

import (
	"context"
	accountProto "interview_task_golang_microservices/protos"
)

func (s *service) UpdateAccountBalance(ctx context.Context, in *accountProto.UpdateAccountBalanceRequest) (*accountProto.UpdateAccountBalanceResponse, error) {
	return &accountProto.UpdateAccountBalanceResponse{}, nil
}
