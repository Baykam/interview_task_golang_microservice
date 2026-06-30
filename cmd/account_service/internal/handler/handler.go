package accountHandler

import (
	accountService "interview_task_golang_microservices/cmd/account_service/internal/service"
	"interview_task_golang_microservices/pkgs/config"
)

type Handler interface{}

type handler struct {
	accountSvc accountService.Service
	cfg        *config.Config
}

func Newhandler(
	accountSvc accountService.Service,
	cfg *config.Config,
) Handler {
	return &handler{accountSvc: accountSvc, cfg: cfg}
}
