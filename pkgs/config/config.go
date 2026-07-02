package config

import (
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"interview_task_golang_microservices/pkgs/redis"
	"interview_task_golang_microservices/pkgs/sql"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres sql.Config `yaml:"database"`

	Redis redis.Config `yaml:"redis"`

	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`

	// Account Service'e özel alanlar ve onun belirlediği kuyruklar
	AccountService struct {
		HttpPort string `yaml:"http_port"`
		Grpc     struct {
			Port    string `yaml:"port"`
			TimeOut string `yaml:"timeout"`
		} `yaml:"grpc"`
		Clients struct {
			TransactionServiceAddr string `yaml:"transaction_service_addr"`
		} `yaml:"clients"`
		Queues struct {
			DepositQueue  string `yaml:"deposit_queue"`
			WithdrawQueue string `yaml:"withdraw_queue"`
			TransferQueue string `yaml:"transfer_queue"`
		} `yaml:"queues"`
	} `yaml:"account_service"`

	// Transaction Service'e özel alanlar ve onun dinleyeceği esnek kuyruk listesi
	TransactionService struct {
		Grpc struct {
			Port    string `yaml:"port"`
			TimeOut string `yaml:"timeout"`
		} `yaml:"grpc"`
		Clients struct {
			AccountServiceAddr string `yaml:"account_service_addr"`
		} `yaml:"clients"`
		Queues struct {
			DepositQueue  string `yaml:"deposit_queue"`
			WithdrawQueue string `yaml:"withdraw_queue"`
			TransferQueue string `yaml:"transfer_queue"`
		} `yaml:"queues"`
	} `yaml:"transaction_service"`
}

func LoadConfig(path string) (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Eğer boşsa, yerelde elle çalıştırırken patlamasın diye varsayılan bir yol verelim
		configPath = "pkgs/config/config.yaml"
	}

	configPath = filepath.Clean(configPath)

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
