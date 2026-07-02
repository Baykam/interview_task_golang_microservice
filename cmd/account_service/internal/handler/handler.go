package handler

import (
	accountHandler "interview_task_golang_microservices/cmd/account_service/internal/handler/account"
	transActionHandler "interview_task_golang_microservices/cmd/account_service/internal/handler/transaction"
	accountService "interview_task_golang_microservices/cmd/account_service/internal/service"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	accountProto "interview_task_golang_microservices/protos"
)

type Handler struct {
	Account     accountHandler.Handler
	Transaction transActionHandler.Handler
}

func Newhandler(
	accountSvc accountService.Service,
	grpcClient accountProto.AccountServiceClient,
	publisher rabbitmq.Publisher,
	cfg *config.Config,
	logger logger.Logger,
) *Handler {
	return &Handler{
		Account:     accountHandler.Newhandler(accountSvc, grpcClient, publisher, cfg, logger),
		Transaction: transActionHandler.NewHandler(accountSvc, grpcClient, publisher, cfg, logger),
	}
}
