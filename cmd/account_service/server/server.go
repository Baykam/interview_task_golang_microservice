package server

import (
	"context"
	accountHandler "interview_task_golang_microservices/cmd/account_service/internal/handler"
	accountRepository "interview_task_golang_microservices/cmd/account_service/internal/repository"
	accountService "interview_task_golang_microservices/cmd/account_service/internal/service"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	accountProto "interview_task_golang_microservices/protos"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type server struct {
	cfg         *config.Config
	logger      logger.Logger
	db          *sqlx.DB
	rdb         *redis.Client
	rmq         *rabbitmq.RabbitMQ
	mux         *http.ServeMux
	httpHandler accountHandler.Handler
	grpcClient  accountProto.AccountServiceClient
	publisher   rabbitmq.Publisher
}

func NewServer(
	cfg *config.Config,
	logger logger.Logger,
	db *sqlx.DB,
	rdb *redis.Client,
	rmq *rabbitmq.RabbitMQ,
) *server {
	return &server{
		cfg:    cfg,
		logger: logger,
		db:     db,
		rdb:    rdb,
		rmq:    rmq,
		mux:    http.NewServeMux(),
	}
}

func (s *server) Run() error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ts, err := s.connectGRPCService(s.cfg.AccountService.Clients.TransactionServiceAddr)
	if err != nil {
		return err
	}
	defer ts.Close()
	s.grpcClient = ts.client

	prod, err := rabbitmq.NewPublisher(s.logger, s.cfg.RabbitMQ.URL)
	if err != nil {
		return err
	}
	s.publisher = prod

	// manage working functions
	{
		accountRepo := accountRepository.NewRepository(s.db, s.logger)
		accountSvc := accountService.NewService(accountRepo, s.rdb, s.rmq)
		s.httpHandler = accountHandler.Newhandler(accountSvc, s.grpcClient, s.publisher, s.cfg, s.logger)
	}

	httpServer := s.httpConnect()

	// GRACEFUL SHUTDOWN
	s.close(ctx, httpServer)
	return nil
}
