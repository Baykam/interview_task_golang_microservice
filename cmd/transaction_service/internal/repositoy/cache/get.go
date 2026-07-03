package transactionCache

import (
	"context"
	"errors"
	"fmt"
	"interview_task_golang_microservices/dto"
	"interview_task_golang_microservices/models"

	"github.com/redis/go-redis/v9"
)

func (c *cache) Get(ctx context.Context, id string) (*models.Transaction, error) {
	if id == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	redisKey := c.redisKey(id)

	data, err := c.redis.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account from redis: %w", err)
	}

	transaction := dto.ByteToTransAction([]byte(data))
	if transaction == nil {
		return nil, fmt.Errorf("Redis is not have this data: %s", id)
	}

	return transaction, nil
}
