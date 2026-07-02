package main

import (
	"interview_task_golang_microservices/cmd/transaction_service/server"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"interview_task_golang_microservices/pkgs/redis"
	"interview_task_golang_microservices/pkgs/sql"
)

func main() {
	logger := logger.NewLogger()
	logger.Info("Transaction Service başlatılıyor...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatalf("Konfigürasyon yükleme hatası: %v", err)
	}

	sql, err := sql.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logger.Fatalf("Postgres DB bağlantı hatası: %v", err)
	}

	redis, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatalf("Redis bağlantı hatası: %v", err)
	}

	rmq, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		logger.Fatalf("RabbitMQ bağlantı hatası: %v", err)
	}

	srv := server.NewServer(cfg, logger, rmq, sql, redis)
	if err := srv.Run(); err != nil {
		logger.Fatalf("Transaction Service çalıştırılamadı: %v", err)
	}
}
