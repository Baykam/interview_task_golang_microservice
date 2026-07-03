package accountService

import (
	"context"
	accountRepository "interview_task_golang_microservices/cmd/account_service/internal/repository"
	"interview_task_golang_microservices/models"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	DeleteAccount(ctx context.Context, id string) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	GetAccountsList(ctx context.Context) ([]models.Account, error)
}

type service struct {
	accountRepo *accountRepository.Repository
	rdb         *redis.Client
	rmq         *rabbitmq.RabbitMQ
	log         logger.Logger
}

func NewService(
	repo *accountRepository.Repository,
	rdb *redis.Client,
	rmq *rabbitmq.RabbitMQ,
	log logger.Logger,
) Service {
	return &service{accountRepo: repo, rdb: rdb, rmq: rmq, log: log}
}
