package transActionDb

import (
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/jmoiron/sqlx"
)

type DB interface{}

type db struct{}

func NewDb(logger logger.Logger, sql *sqlx.DB) DB {
	return &db{}
}
