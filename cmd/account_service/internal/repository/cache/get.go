package accountRepositoryCache

import (
	"context"
	"errors"
	"fmt"
	"interview_task_golang_microservices/dto"
	"interview_task_golang_microservices/models"

	"github.com/redis/go-redis/v9"
)

func (c *cache) Get(ctx context.Context, id string) (*models.Account, error) {
	if id == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	redisKey := redisKey(id)

	cachedData, err := c.redis.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account from redis: %w", err)
	}

	account, err := dto.ByteToAccount([]byte(cachedData))
	if err != nil {
		return nil, err
	}

	return account, nil
}
