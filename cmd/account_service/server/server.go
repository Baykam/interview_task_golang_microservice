package server

import (
	"context"
	accountHandler "interview_task_golang_microservices/cmd/account_service/internal/handler"
	accountRepository "interview_task_golang_microservices/cmd/account_service/internal/repository"
	accountService "interview_task_golang_microservices/cmd/account_service/internal/service"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
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

	// manage working functions
	{
		accountRepo := accountRepository.NewRepository(s.db, s.logger)
		accountSvc := accountService.NewService(accountRepo, s.rdb, s.rmq)
		s.httpHandler = accountHandler.Newhandler(accountSvc, s.cfg)
	}

	httpServer := s.httpConnect()

	// GRACEFUL SHUTDOWN
	s.close(ctx, httpServer)
	return nil
}
