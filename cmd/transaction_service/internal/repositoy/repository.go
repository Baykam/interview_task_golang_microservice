package transActionRepository

import (
	transactionCache "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy/cache"
	transActionDb "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy/db"
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Cache transactionCache.Cache
	DB    transActionDb.DB
}

func NewRepository(logger logger.Logger, redis *redis.Client, sql *sqlx.DB) *Repository {
	return &Repository{
		Cache: transactionCache.NewCache(logger, redis),
		DB:    transActionDb.NewDb(logger, sql),
	}
}
