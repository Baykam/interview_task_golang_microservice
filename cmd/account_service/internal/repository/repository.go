package accountRepository

import (
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/jmoiron/sqlx"
)

type Repository interface{}

type repos struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewRepository(
	db *sqlx.DB,
	log logger.Logger,
) Repository {
	return &repos{
		db:  db,
		log: log,
	}
}
