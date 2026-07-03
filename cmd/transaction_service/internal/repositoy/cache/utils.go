package transactionCache

import "fmt"

func (c *cache) redisKey(id string) string {
	return fmt.Sprintf("transaction:%s", id)
}
