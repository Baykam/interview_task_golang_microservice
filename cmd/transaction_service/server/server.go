package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	transActionRepository "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy"
	transActionService "interview_task_golang_microservices/cmd/transaction_service/internal/service"
	transactionConn "interview_task_golang_microservices/cmd/transaction_service/server/rabbit"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Server interface {
	Run() error
}

type server struct {
	cfg     *config.Config
	logger  logger.Logger
	sql     *sqlx.DB
	redis   *redis.Client
	rmq     *rabbitmq.RabbitMQ
	grpcS   *grpc.Server
	service *transActionService.Service
}

func NewServer(cfg *config.Config, logger logger.Logger, rmq *rabbitmq.RabbitMQ, sql *sqlx.DB, redis *redis.Client) Server {
	return &server{cfg: cfg, logger: logger, rmq: rmq}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	repos := transActionRepository.NewRepository(s.logger, s.redis, s.sql)
	s.service = transActionService.NewService(s.logger, repos)

	go func() {
		if err := s.startGRPC(); err != nil {
			s.logger.Error("error in grpc connection")
		}
	}()

	wg := &sync.WaitGroup{}
	transActionMessageProcessor := transactionConn.NewTransactionMessageProcessor(s.logger, s.cfg, s.service)
	transActionMessageProcessor.ProcessMessages(ctx, s.rmq, s.getConsumerGroupTopics(), transactionConn.PoolSize, wg)

	s.close(ctx)
	wg.Wait()
	return nil
}
