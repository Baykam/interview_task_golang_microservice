package transactionCache

import (
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/redis/go-redis/v9"
)

type Cache interface{}

type cache struct{}

func NewCache(logger logger.Logger, redis *redis.Client) Cache {
	return &cache{}
}
