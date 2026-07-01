package main

import (
	"interview_task_golang_microservices/cmd/transaction_service/server"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
)

func main() {
	logger := logger.NewLogger()
	logger.Info("Transaction Service başlatılıyor...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatalf("Konfigürasyon yükleme hatası: %v", err)
	}

	rmq, err := rabbitmq.NewRabbitMQ(rabbitmq.Config{
		URL:         cfg.RabbitMQ.URL,
		DialTimeout: cfg.RabbitMQ.DialTimeout,
	})
	if err != nil {
		logger.Fatalf("RabbitMQ bağlantı hatası: %v", err)
	}
	defer rmq.Close()

	srv := server.NewServer(cfg, logger, rmq)
	if err := srv.Run(); err != nil {
		logger.Fatalf("Transaction Service çalıştırılamadı: %v", err)
	}
}
