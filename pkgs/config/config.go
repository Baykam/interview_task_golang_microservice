package config

import (
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"interview_task_golang_microservices/pkgs/redis"
	"interview_task_golang_microservices/pkgs/sql"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres sql.Config `yaml:"postgres"`

	Redis redis.Config `yaml:"redis"`

	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`

	// Account Service'e özel alanlar ve onun belirlediği kuyruklar
	AccountService struct {
		HttpPort               string `yaml:"http_port"`
		GrpcPort               string `yaml:"grpc_port"`
		TransactionServiceAddr string `yaml:"transaction_service_addr"`
		Queues                 struct {
			DepositQueue  string `yaml:"deposit_queue"`
			WithdrawQueue string `yaml:"withdraw_queue"`
			TransferQueue string `yaml:"transfer_queue"`
		} `yaml:"queues"`
	} `yaml:"account_service"`

	// Transaction Service'e özel alanlar ve onun dinleyeceği esnek kuyruk listesi
	TransactionService struct {
		GrpcPort           string `yaml:"grpc_port"`
		AccountServiceAddr string `yaml:"account_service_addr"`
		Queues             struct {
			ListenQueues []string `yaml:"listen_queues"`
		} `yaml:"queues"`
	} `yaml:"transaction_service"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
