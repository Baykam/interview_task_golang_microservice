package accountService

import (
	accountRepository "interview_task_golang_microservices/cmd/account_service/internal/repository"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"

	"github.com/redis/go-redis/v9"
)

type Service interface{}

type service struct {
	accountRepo accountRepository.Repository
	rdb         *redis.Client
	rmq         *rabbitmq.RabbitMQ
}

func NewService(
	repo accountRepository.Repository,
	rdb *redis.Client,
	rmq *rabbitmq.RabbitMQ,
) Service {
	return &service{accountRepo: repo, rdb: rdb, rmq: rmq}
}
