package main

import (
	"interview_task_golang_microservices/cmd/account_service/server"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"interview_task_golang_microservices/pkgs/redis"
	"interview_task_golang_microservices/pkgs/sql"
)

func main() {

	logger := logger.NewLogger()
	logger.Info("Uygulama main.go üzerinden başlatılıyor...")

	// 2. CONFIG YÜKLEME
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatalf("Konfigürasyon yükleme hatası: %v", err)
	}
	logger.Info("Konfigürasyon (config.yaml) başarıyla okundu.")

	// 3. POSTGRESQL BAĞLANTISI (CONNECTION POOL)
	db, err := sql.NewPostgresDB(sql.Config{
		DSN:             cfg.Postgres.DSN,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
	})
	if err != nil {
		logger.Fatalf("PostgreSQL bağlantı hatası: %v", err)
	}
	logger.Info("PostgreSQL bağlantı havuzu başarıyla oluşturuldu.")

	// 4. REDIS BAĞLANTISI (CLIENT POOL)
	rdb, err := redis.NewRedisClient(redis.Config{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  cfg.Redis.DialTimeout,
	})
	if err != nil {
		logger.Fatalf("Redis bağlantı hatası: %v", err)
	}
	logger.Info("Redis önbellek istemcisi başarıyla bağlandı.")

	// 5. RABBITMQ BAĞLANTISI (CONNECTION & CHANNEL)
	rmq, err := rabbitmq.NewRabbitMQ(rabbitmq.Config{
		URL:         cfg.RabbitMQ.URL,
		DialTimeout: cfg.RabbitMQ.DialTimeout,
	})
	if err != nil {
		logger.Fatalf("RabbitMQ bağlantı hatası: %v", err)
	}

	server := server.NewServer(cfg, logger, db, rdb, rmq)
	if err := server.Run(); err != nil {
		logger.Fatalf("error for: %v", err)
	}
}
