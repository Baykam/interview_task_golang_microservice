package transActionHandler

import (
	accountService "interview_task_golang_microservices/cmd/account_service/internal/service"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	accountProto "interview_task_golang_microservices/protos"
	"net/http"
)

type Handler interface {
	Deposit(w http.ResponseWriter, r *http.Request)
	GetTransactionHistory(w http.ResponseWriter, r *http.Request)
	Transfer(w http.ResponseWriter, r *http.Request)
	Withdraw(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	accountSvc accountService.Service
	grpcClient accountProto.AccountServiceClient
	publisher  rabbitmq.Publisher
	cfg        *config.Config
	logger     logger.Logger
}

func NewHandler(
	accountSvc accountService.Service,
	grpcClient accountProto.AccountServiceClient,
	publisher rabbitmq.Publisher,
	cfg *config.Config,
	logger logger.Logger,
) Handler {
	return &handler{
		accountSvc: accountSvc,
		grpcClient: grpcClient,
		publisher:  publisher,
		cfg:        cfg,
		logger:     logger,
	}
}
