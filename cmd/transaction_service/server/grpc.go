package server

import (
	accountProto "interview_task_golang_microservices/protos"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxConnectionIdle = 5 * time.Minute
	gRPCTimeout       = 15 * time.Second
	maxConnectionAge  = 5 * time.Minute
	gRPCTime          = 10 * time.Minute
)

func (s *server) startGRPC() error {
	grpcPort := s.cfg.TransactionService.Grpc.Port
	if grpcPort == "" {
		grpcPort = "50052"
	}
	lis, err := net.Listen("tcp", ":"+grpcPort)

	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle,
			Timeout:           gRPCTimeout,
			MaxConnectionAge:  maxConnectionAge,
			Time:              gRPCTime,
		}),
	)

	accountProto.RegisterAccountServiceServer(grpcServer, s.service.Queries)

	return grpcServer.Serve(lis)
}
