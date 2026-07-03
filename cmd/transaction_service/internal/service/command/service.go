package commandService

import (
	"context"
	transActionRepository "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy"
	"interview_task_golang_microservices/models"
	"interview_task_golang_microservices/pkgs/logger"
)

type Service interface {
	CreateTransAction(ctx context.Context, data models.Transaction) error
}

type service struct {
	logger logger.Logger
	repos  *transActionRepository.Repository
}

func NewService(logger logger.Logger, repos *transActionRepository.Repository) Service {
	return &service{
		logger: logger,
		repos:  repos,
	}
}
