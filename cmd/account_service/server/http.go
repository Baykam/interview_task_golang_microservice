package server

import "net/http"

func (s *server) httpConnect() *http.Server {
	s.mux.HandleFunc("/api/accounts", s.httpHandler.GetAccount)

	httpServer := &http.Server{
		Addr:    ":" + s.cfg.AccountService.HttpPort,
		Handler: s.mux,
	}

	go func() {
		s.logger.Info("HTTP API Gateway %s portunda dinlemeye başladı...", "port_number", s.cfg.AccountService.HttpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Debug("HTTP Sunucu durdurma hatası: %v", err)
		}
	}()

	return httpServer
}
