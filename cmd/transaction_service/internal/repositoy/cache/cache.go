package transactionCache

import (
	"context"
	"interview_task_golang_microservices/models"
	"interview_task_golang_microservices/pkgs/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ttl = 24 * time.Hour
)

type Cache interface {
	Add(ctx context.Context, transAction models.Transaction) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*models.Transaction, error)
}

type cache struct {
	log   logger.Logger
	redis *redis.Client
}

func NewCache(logger logger.Logger, redis *redis.Client) Cache {
	return &cache{log: logger, redis: redis}
}
