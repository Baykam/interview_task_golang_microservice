package transactionCache

import (
	"context"
	"fmt"
)

func (c *cache) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	redisKey := c.redisKey(id)

	err := c.redis.Del(ctx, redisKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete account from redis: %w", err)
	}

	return nil
}
