package transactionCache

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/dto"
	"interview_task_golang_microservices/models"
)

func (c *cache) Add(ctx context.Context, transAction models.Transaction) error {
	if transAction.Id == nil || *transAction.Id == "" {
		return fmt.Errorf("Id must have")
	}

	rKey := c.redisKey(*transAction.Id)

	data := dto.TransActionToByte(transAction)

	err := c.redis.Set(ctx, rKey, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set account in redis: %w", err)
	}

	return nil
}
