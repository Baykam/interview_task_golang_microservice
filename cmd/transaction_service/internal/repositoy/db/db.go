package transActionDb

import (
	"context"
	"interview_task_golang_microservices/models"
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	GetListByAccountID(ctx context.Context, accountID string) ([]models.Transaction, error)
	Create(ctx context.Context, data models.Transaction) (*string, error)
}

type db struct {
	logger logger.Logger
	db     *sqlx.DB
}

func NewDb(logger logger.Logger, sql *sqlx.DB) DB {
	return &db{logger: logger, db: sql}
}
