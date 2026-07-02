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

	apiMux := http.NewServeMux()

	s.accountRoutes(apiMux)

	s.mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}

func (s *server) accountRoutes(apiMux *http.ServeMux) {

	accountMux := http.NewServeMux()

	accountMux.HandleFunc("POST /api/accounts", s.httpHandler.Account.CreateAccount)
	accountMux.HandleFunc("GET /api/accounts", s.httpHandler.Account.GetAccountsList)
	accountMux.HandleFunc("GET /api/accounts/{id}", s.httpHandler.Account.GetAccountById)
	accountMux.HandleFunc("DELETE /api/accounts/{id}", s.httpHandler.Account.DeleteAccount)

	accountMux.HandleFunc("POST /api/accounts/{id}/deposit", s.httpHandler.Transaction.Deposit)
	accountMux.HandleFunc("POST /api/accounts/{id}/withdraw", s.httpHandler.Transaction.Withdraw)
	accountMux.HandleFunc("POST /api/accounts/{id}/transfer", s.httpHandler.Transaction.Transfer)
	accountMux.HandleFunc("GET /api/accounts/{id}/transactions", s.httpHandler.Transaction.GetTransactionHistory)

	apiMux.Handle("/api/", http.StripPrefix("/api", accountMux))
}
