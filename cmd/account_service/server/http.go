package server

import "net/http"

func (s *server) httpConnect() *http.Server {
	s.registerRoutes()

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

func (s *server) registerRoutes() {
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "UP", "service": "account-service"}`))
	})

	apiMux := http.NewServeMux()

	s.accountRoutes(apiMux)

	s.mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}

func (s *server) accountRoutes(apiMux *http.ServeMux) {
	accountMux := http.NewServeMux()

	accountMux.HandleFunc("POST /", s.httpHandler.Account.CreateAccount)
	accountMux.HandleFunc("GET /", s.httpHandler.Account.GetAccountsList)
	accountMux.HandleFunc("GET /{id}", s.httpHandler.Account.GetAccountById)
	accountMux.HandleFunc("DELETE /{id}", s.httpHandler.Account.DeleteAccount)

	accountMux.HandleFunc("POST /{id}/deposit", s.httpHandler.Transaction.Deposit)
	accountMux.HandleFunc("POST /{id}/withdraw", s.httpHandler.Transaction.Withdraw)
	accountMux.HandleFunc("POST /{id}/transfer", s.httpHandler.Transaction.Transfer)
	accountMux.HandleFunc("GET /{id}/transactions", s.httpHandler.Transaction.GetTransactionHistory)

	apiMux.Handle("/accounts/", http.StripPrefix("/accounts", accountMux))
}
