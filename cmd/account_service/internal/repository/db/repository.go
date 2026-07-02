package accountRepositoryDb

import (
	"context"
	"interview_task_golang_microservices/models"
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, account *models.Account) error
	GetList(ctx context.Context) ([]models.Account, error)
	GetById(ctx context.Context, id string) (*models.Account, error)
	Delete(ctx context.Context, account *models.Account) error
}

type repository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewRepository(
	db *sqlx.DB,
	log logger.Logger) Repository {
	return &repository{
		db: db, log: log,
	}
}
