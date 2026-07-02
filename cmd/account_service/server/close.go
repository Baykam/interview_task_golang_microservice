package server

import (
	"context"
	"net/http"
)

func (s *server) close(ctx context.Context, httpServer *http.Server) {

	<-ctx.Done()

	if err := httpServer.Shutdown(ctx); err != nil {
		s.logger.Debug("HTTP Sunucu kapatılırken hata oluştu: %v", err)
	}

	s.rmq.Close()
	s.rdb.Close()
	s.db.Close()

	s.logger.Debug("Tüm bağlantılar ve havuzlar temizlendi. Uygulama güvenli bir şekilde sonlandırıldı (Safe Exit).")

}
