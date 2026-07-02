package accountRepository

import (
	accountRepositoryCache "interview_task_golang_microservices/cmd/account_service/internal/repository/cache"
	accountRepositoryDb "interview_task_golang_microservices/cmd/account_service/internal/repository/db"
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Cache accountRepositoryCache.Cache
	DB    accountRepositoryDb.Repository
}

func NewRepository(
	db *sqlx.DB,
	redis *redis.Client,
	log logger.Logger,
) *Repository {
	return &Repository{
		Cache: accountRepositoryCache.NewCache(log, redis),
		DB:    accountRepositoryDb.NewRepository(db, log),
	}
}
