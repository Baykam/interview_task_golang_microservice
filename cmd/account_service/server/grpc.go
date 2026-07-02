package server

import (
	accountProto "interview_task_golang_microservices/protos"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	backOffLinear  = 100 * time.Millisecond
	backOffRetries = 3
)

type TransactionClient struct {
	client accountProto.AccountServiceClient
	conn   *grpc.ClientConn
}

func (s *server) connectGRPCService(target string) (*TransactionClient, error) {
	opt := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backOffLinear)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
		grpc_retry.WithMax(backOffRetries),
	}
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opt...)))
	if err != nil {
		return nil, errors.Wrap(err, "grpc.NewClient")
	}

	client := accountProto.NewAccountServiceClient(conn)
	return &TransactionClient{conn: conn, client: client}, nil
}

func (s *TransactionClient) Close() {
	s.conn.Close()
}
