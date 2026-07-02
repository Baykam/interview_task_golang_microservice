package accountRepositoryCache

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/dto"
	"interview_task_golang_microservices/models"
)

func (c *cache) Add(ctx context.Context, account *models.Account) error {
	if account == nil || account.ID == nil {
		return fmt.Errorf("account or account ID cannot be nil")
	}

	redisKey := redisKey(*account.ID)

	accountData := dto.AccountToByte(account)

	err := c.redis.Set(ctx, redisKey, accountData, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set account in redis: %w", err)
	}

	return nil
}
