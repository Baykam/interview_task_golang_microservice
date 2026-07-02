package server

import "context"

func (s *server) close(ctx context.Context) {
	<-ctx.Done()
	s.logger.Info("Transaction Service kapanıyor...")
	s.rmq.Close()
	s.grpcS.GracefulStop()
}
