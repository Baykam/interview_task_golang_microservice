package accountRepositoryCache

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
	Get(ctx context.Context, id string) (*models.Account, error)
	Add(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id string) error
}

type cache struct {
	log   logger.Logger
	redis *redis.Client
}

func NewCache(
	log logger.Logger,
	redis *redis.Client,
) Cache {
	return &cache{
		log:   log,
		redis: redis,
	}
}
