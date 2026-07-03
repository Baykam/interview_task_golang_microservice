package transactionConn

import (
	"context"
	"fmt"
	transActionService "interview_task_golang_microservices/cmd/transaction_service/internal/service"
	"interview_task_golang_microservices/dto"
	"interview_task_golang_microservices/models"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"strings"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PoolSize int = 5

type transactionMessageProcessor struct {
	log logger.Logger
	cfg *config.Config
	ts  *transActionService.Service
}

func NewTransactionMessageProcessor(
	log logger.Logger,
	cfg *config.Config,
	ts *transActionService.Service,
) *transactionMessageProcessor {
	return &transactionMessageProcessor{log: log, cfg: cfg, ts: ts}
}

// ProcessMessages artık ana başlatıcı (orchestrator) görevini görüyor
func (t *transactionMessageProcessor) ProcessMessages(ctx context.Context, rmq *rabbitmq.RabbitMQ, queues []string, poolSize int, wg *sync.WaitGroup) {

	// ch := rmq.Channel()

	for _, queueName := range queues {
		for i := 1; i <= poolSize; i++ {
			wg.Add(1)
			go t.worker(ctx, rmq, queueName, i, wg)
		}
	}
}

func (t *transactionMessageProcessor) worker(ctx context.Context, rmq *rabbitmq.RabbitMQ, queueName string, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	ch, err := rmq.NewChannel()
	if err != nil {
		t.log.Error("workerId: %v, %s için kanal açılamadı: %v", "workerId", workerId, "queueName", queueName, "error", err)
		return
	}
	defer ch.Close()

	err = ch.Qos(
		1,     // Prefetch Count: Her worker aynı anda sadece 1 mesaj alabilir
		0,     // Prefetch Size
		false, // Global: Sadece bu kanal için geçerli
	)
	if err != nil {
		return
	}

	_, err = ch.QueueDeclare(
		queueName, // Dinlenecek kuyruğun adı
		true,      // Durable: Sunucu kapansa da kuyruk kaybolmasın (Config'indeki yapıya göre)
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		t.log.Error("Worker-%d [%s] kuyruğu oluşturulamadı: %v", "workerId", workerId, "queueName", queueName, "error", err)
		return
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.log.Error("workerId: %v, %s kuyruğu dinlenemedi: %v", "workerId", workerId, "queueName", queueName, err)
		return
	}

	t.log.Info("Worker-%v [%s] kuyruğunu dinlemeye başladı.", "workerId", workerId, "queueName", queueName)

	for {
		select {
		case <-ctx.Done():
			t.log.Info("Worker-%v [%s] kapatılıyor...", "workerId", workerId, "queueName", queueName)
			return
		case message, ok := <-msgs:
			if !ok {
				t.log.Info("Worker-%v, %s kuyruk kanalı kapandı.", "workerId", workerId, "queueName", queueName)
				return
			}

			t.handleMessageByQueue(ctx, queueName, message, workerId)
		}
	}
}

func (t *transactionMessageProcessor) handleMessageByQueue(ctx context.Context, queueName string, message amqp.Delivery, workerId int) {
	var err error
	var txType models.TransactionType

	switch queueName {
	case t.cfg.TransactionService.Queues.DepositQueue:
		txType = models.TransactionTypeDeposit
	case t.cfg.TransactionService.Queues.WithdrawQueue:
		txType = models.TransactionTypeWithdraw
	case t.cfg.TransactionService.Queues.TransferQueue:
		txType = models.TransactionTypeTransfer
	default:
		t.log.Error("Worker-%v: Unsupported queue name: %s", workerId, queueName)
		message.Nack(false, false) // Kuyruktan tamamen at (requeue: false)
		return
	}

	t.log.Info(fmt.Sprintf("Worker-%v [%s] message received", "workerId", workerId, "transActionType", txType))

	txModel := dto.ByteToTransAction(message.Body)
	if txModel == nil {
		t.log.Error("Worker-%v: JSON unmarshal error: %v", workerId, err)
		message.Nack(false, false)
		return
	}

	txModel.TransactionType = txType

	err = t.ts.Command.CreateTransAction(ctx, *txModel)

	if err != nil {
		t.log.Error("Worker-%v, Transaction execution error: %v", "workerId", workerId, "error", err)

		if strings.Contains(err.Error(), "insufficient balance") {
			message.Nack(false, false)
		} else {
			message.Nack(false, true) // Hata sistemselse sıraya geri bırak
		}
		return
	}

	if err := message.Ack(false); err != nil {
		t.log.Error("Worker-%v, ACK error: %v", workerId, err)
	}
}
