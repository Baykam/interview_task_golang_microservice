package querieService

import (
	"context"
	accountProto "interview_task_golang_microservices/protos"
)

func (s *service) CheckAccountExists(ctx context.Context, in *accountProto.CheckAccountExistsRequest) (*accountProto.CheckAccountExistsResponse, error) {
	return &accountProto.CheckAccountExistsResponse{}, nil
}
