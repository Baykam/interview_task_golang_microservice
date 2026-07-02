package rabbitmq

import (
	"context"
	"interview_task_golang_microservices/pkgs/logger"
	"sync"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Worker rabbitmq consumer worker fetches and processes messages from queue
type Worker func(ctx context.Context, ch *amqp.Channel, queueName string, wg *sync.WaitGroup, workerID int)

type ConsumerManager interface {
	GetChannel() (*amqp.Channel, error)
	Publish(ctx context.Context, queueName string, body []byte) error
	ConsumeQueue(ctx context.Context, queueName string, poolSize int, worker Worker)
	Close()
}

type consumerManager struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	log  logger.Logger
}

// NewConsumerManager, verilen yapılandırma ile bir RabbitMQ bağlantısı ve kanalı başlatır.
func NewConsumerManager(cfg Config, log logger.Logger) (ConsumerManager, error) {
	amqpCfg := amqp.Config{
		Dial: amqp.DefaultDial(cfg.DialTimeout),
	}
	if cfg.DialTimeout == 0 {
		amqpCfg.Dial = amqp.DefaultDial(3 * time.Second)
	}

	conn, err := amqp.DialConfig(cfg.URL, amqpCfg)
	if err != nil {
		return nil, errors.Wrap(err, "rabbitmq baglanti hatasi")
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "rabbitmq kanal acma hatasi")
	}

	return &consumerManager{
		conn: conn,
		ch:   ch,
		log:  log,
	}, nil
}

// GetChannel mevcut veya yeni bir kanal dönmek için kullanılabilir
func (c *consumerManager) GetChannel() (*amqp.Channel, error) {
	if c.ch == nil || c.ch.IsClosed() {
		var err error
		c.ch, err = c.conn.Channel()
		if err != nil {
			return nil, errors.Wrap(err, "yeni kanal acilamadi")
		}
	}
	return c.ch, nil
}

// Close, hem kanalı hem de bağlantıyı güvenli bir şekilde kapatır.
func (c *consumerManager) Close() {
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

// Publish (Producer), belirtilen kuyruğa bir mesaj gönderir.
func (c *consumerManager) Publish(ctx context.Context, queueName string, body []byte) error {
	ch, err := c.GetChannel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return errors.Wrapf(err, "kuyruk deklarasyon hatasi: %s", queueName)
	}

	return ch.PublishWithContext(ctx,
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

// ConsumeQueue tıpkı Kafka'daki ConsumeTopic gibi çalışır.
// Verilen worker'ı belirtilen poolSize kadar goroutine üzerinde ayağa kaldırır.
func (c *consumerManager) ConsumeQueue(ctx context.Context, queueName string, poolSize int, worker Worker) {
	ch, err := c.GetChannel()
	if err != nil {
		c.log.Error("ConsumeQueue baslatilamadi, kanal hatasi: %v", err)
		return
	}

	// 1. Kuyruğun varlığından emin ol
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		c.log.Error("kuyruk declare edilemedi: %s, err: %v", queueName, err)
		return
	}

	// 2. Worker Pool'un adil yük dağıtımı (Fair Dispatch) yapabilmesi için QoS ayarı
	err = ch.Qos(1, 0, false)
	if err != nil {
		c.log.Error("QoS ayari yapilamadi, err: %v", err)
		return
	}

	c.log.Info("Starting RabbitMQ consumer, queue: %s, pool size: %v", queueName, poolSize)

	wg := &sync.WaitGroup{}

	// Kafka'daki döngü mantığının aynısıyla worker'ları başlatıyoruz
	for i := 1; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, ch, queueName, wg, i)
	}

	// Tüm worker'ların durmasını bekle (Context iptal edildiğinde tetiklenecekler)
	wg.Wait()
}
