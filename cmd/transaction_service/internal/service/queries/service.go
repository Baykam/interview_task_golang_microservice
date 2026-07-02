package querieService

import (
	"context"
	transActionRepository "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy"
	"interview_task_golang_microservices/pkgs/logger"
	accountProto "interview_task_golang_microservices/protos"
)

type Service interface {
	accountProto.AccountServiceServer
	GetAccountBalance(ctx context.Context, in *accountProto.GetAccountBalanceRequest) (*accountProto.GetAccountBalanceResponse, error)
	UpdateAccountBalance(ctx context.Context, in *accountProto.UpdateAccountBalanceRequest) (*accountProto.UpdateAccountBalanceResponse, error)
	CheckAccountExists(ctx context.Context, in *accountProto.CheckAccountExistsRequest) (*accountProto.CheckAccountExistsResponse, error)
}

type service struct {
	accountProto.UnimplementedAccountServiceServer
	logger logger.Logger
	repos  *transActionRepository.Repository
}

func NewService(logger logger.Logger, repos *transActionRepository.Repository) Service {
	return &service{logger: logger, repos: repos}
}
