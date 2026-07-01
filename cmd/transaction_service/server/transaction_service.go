package server

import (
	"context"
	"fmt"
	"time"

	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	pp "interview_task_golang_microservices/protos"

	"google.golang.org/grpc"
)

type transactionServiceServer struct {
	pp.UnimplementedAccountServiceServer
	cfg    *config.Config
	logger logger.Logger
}

func newTransactionServiceServer(cfg *config.Config, logger logger.Logger) *transactionServiceServer {
	return &transactionServiceServer{cfg: cfg, logger: logger}
}

func (s *transactionServiceServer) GetAccountBalance(ctx context.Context, req *pp.GetAccountBalanceRequest) (*pp.GetAccountBalanceResponse, error) {
	s.logger.Info("GetAccountBalance çağrısı alındı", "account_id", req.AccountId)
	return &pp.GetAccountBalanceResponse{
		AccountId: req.AccountId,
		Balance:   0,
		Currency:  "TRY",
	}, nil
}

func (s *transactionServiceServer) UpdateAccountBalance(ctx context.Context, req *pp.UpdateAccountBalanceRequest) (*pp.UpdateAccountBalanceResponse, error) {
	s.logger.Info("UpdateAccountBalance çağrısı alındı", "account_id", req.AccountId, "amount", req.Amount)
	return &pp.UpdateAccountBalanceResponse{
		Success:    true,
		NewBalance: req.Amount,
	}, nil
}

func (s *transactionServiceServer) CheckAccountExists(ctx context.Context, req *pp.CheckAccountExistsRequest) (*pp.CheckAccountExistsResponse, error) {
	s.logger.Info("CheckAccountExists çağrısı alındı", "account_id", req.AccountId)
	return &pp.CheckAccountExistsResponse{Exists: true}, nil
}

func dialAccountService(ctx context.Context, address string) (pp.AccountServiceClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("account service grpc dial hatasi: %w", err)
	}
	return pp.NewAccountServiceClient(conn), nil
}
