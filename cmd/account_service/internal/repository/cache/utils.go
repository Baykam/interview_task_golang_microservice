package accountRepositoryCache

import (
	"fmt"
)

func redisKey(id string) string {
	redisKey := fmt.Sprintf("account:%s", id)
	return redisKey
}
