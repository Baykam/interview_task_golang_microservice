package transActionService

import (
	transActionRepository "interview_task_golang_microservices/cmd/transaction_service/internal/repositoy"
	commandService "interview_task_golang_microservices/cmd/transaction_service/internal/service/command"
	querieService "interview_task_golang_microservices/cmd/transaction_service/internal/service/queries"
	"interview_task_golang_microservices/pkgs/logger"
)

type Service struct {
	Queries querieService.Service
	Command commandService.Service
}

func NewService(logger logger.Logger, repos *transActionRepository.Repository) *Service {
	return &Service{
		Queries: querieService.NewService(logger, repos),
		Command: commandService.NewService(logger, repos),
	}
}
