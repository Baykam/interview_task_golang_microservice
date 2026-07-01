package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server interface {
	Run() error
}

type server struct {
	cfg    *config.Config
	logger logger.Logger
	rmq    *rabbitmq.RabbitMQ
	grpcS  *grpc.Server
}

func NewServer(cfg *config.Config, logger logger.Logger, rmq *rabbitmq.RabbitMQ) Server {
	return &server{cfg: cfg, logger: logger, rmq: rmq}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := s.startGRPC(cancel); err != nil {
		return err
	}

	var wg sync.WaitGroup
	if err := s.startConsumers(ctx, &wg); err != nil {
		return err
	}

	<-ctx.Done()
	s.logger.Info("Transaction Service kapanıyor...")
	s.grpcS.GracefulStop()
	s.rmq.Close()
	wg.Wait()
	return nil
}

func (s *server) startGRPC(cancel context.CancelFunc) error {
	grpcPort := s.cfg.TransactionService.GrpcPort
	if grpcPort == "" {
		grpcPort = "50052"
	}

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return fmt.Errorf("gRPC listener oluşturulamadı: %w", err)
	}

	s.grpcS = grpc.NewServer()
	registerTransactionServiceServer(s.grpcS, &transactionServiceGRPC{logger: s.logger})

	go func() {
		s.logger.Info("Transaction gRPC %s portunda dinliyor...", grpcPort)
		if err := s.grpcS.Serve(listener); err != nil {
			s.logger.Error("gRPC sunucu hatası", "error", err)
			cancel()
		}
	}()

	return nil
}

func (s *server) startConsumers(ctx context.Context, wg *sync.WaitGroup) error {
	queues := s.cfg.TransactionService.Queues.ListenQueues
	if len(queues) == 0 {
		return fmt.Errorf("transaction service için dinlenecek kuyruk bulunamadı")
	}

	for _, queueName := range queues {
		queueName := queueName
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.rmq.Consume(queueName, func(body []byte) {
				s.handleMessage(queueName, body)
			}); err != nil {
				s.logger.Error("RabbitMQ kuyruk dinleme hatası", "queue", queueName, "error", err)
				return
			}

			<-ctx.Done()
		}()
	}

	return nil
}

func (s *server) handleMessage(queueName string, body []byte) {
	s.logger.Info("RabbitMQ mesajı alındı", "queue", queueName, "size", len(body))
	// TODO: Burada mesaj parse edilip işlem yapılabilir.
}

type TransactionServiceServer interface {
	HealthCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

type transactionServiceGRPC struct {
	logger logger.Logger
}

func (t *transactionServiceGRPC) HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	t.logger.Info("HealthCheck çağrısı alındı")
	return &emptypb.Empty{}, nil
}

func _TransactionService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/HealthCheck"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TransactionService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "transaction.TransactionService",
	HandlerType: (*TransactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _TransactionService_HealthCheck_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func registerTransactionServiceServer(s *grpc.Server, srv TransactionServiceServer) {
	s.RegisterService(&_TransactionService_serviceDesc, srv)
}
