package transactionConn

import (
	"context"
	transActionService "interview_task_golang_microservices/cmd/transaction_service/internal/service"
	"interview_task_golang_microservices/pkgs/config"
	"interview_task_golang_microservices/pkgs/logger"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
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

// İş mantığını ayıran yardımcı fonksiyon
func (t *transactionMessageProcessor) handleMessageByQueue(ctx context.Context, queueName string, message amqp.Delivery, workerId int) {
	var err error

	switch queueName {
	case t.cfg.TransactionService.Queues.DepositQueue:
		t.log.Info("Worker-%v [Deposit] mesajı aldı", "workerId", workerId)
		// err = t.ts.HandleDeposit(ctx, message.Body)

	case t.cfg.TransactionService.Queues.WithdrawQueue:
		t.log.Info("Worker-%v [Withdraw] mesajı aldı", "workerId", workerId)
		// err = t.ts.HandleWithdraw(ctx, message.Body)

	case t.cfg.TransactionService.Queues.TransferQueue:
		t.log.Info("Worker-%v [Transfer] mesajı aldı", "workerId", workerId)
		// err = t.ts.HandleTransfer(ctx, message.Body)
	}

	if err != nil {
		t.log.Error("Worker-%v, Mesaj işleme hatası: %v", "workerId", workerId, "error", err)
		// İsteğe bağlı: message.Nack(false, true) ile hata durumunda kuyruğa geri bırakılabilir
		return
	}

	if err := message.Ack(false); err != nil {
		t.log.Error("Worker-%v, ACK hatası: %v", workerId, err)
	}
}
