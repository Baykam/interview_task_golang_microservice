package commandService

import (
	transActionRepository "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy"
	"interview_task_golang_microservices/pkgs/logger"
)

type Service interface{}

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
