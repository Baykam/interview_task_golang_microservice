package rabbitmq

import (
	"context"
	"fmt"
	"interview_task_golang_microservices/pkgs/logger"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Message, RabbitMQ'ya mesaj gönderirken esneklik sağlamak için
// Kafka.Message yapısına benzer bir sarmalayıcı (wrapper) struct.
type Message struct {
	QueueName   string // Mesajın gideceği hedef kuyruk adı
	ContentType string // Örn: "application/json" veya "text/plain"
	Body        []byte // Gönderilecek veri byte array olarak
}

type Publisher interface {
	PublishMessage(ctx context.Context, msgs ...Message) error
	Close() error
}

type publisher struct {
	log  logger.Logger
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewPublisher, RabbitMQ sunucusuna bağlanır, kanalı açar ve Publisher nesnesini döner.
func NewPublisher(log logger.Logger, url string) (*publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq publisher baglanti hatasi: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq publisher kanal acma hatasi: %w", err)
	}

	return &publisher{
		log:  log,
		url:  url,
		conn: conn,
		ch:   ch,
	}, nil
}

// PublishMessage, variadic (...) parametre sayesinde ister tek, ister toplu mesaj göndermeyi sağlar.
func (p *publisher) PublishMessage(ctx context.Context, msgs ...Message) error {
	for _, msg := range msgs {
		// Eğer ContentType boş bırakıldıysa varsayılan olarak JSON ata
		contentType := msg.ContentType
		if contentType == "" {
			contentType = "application/json"
		}

		// 1. Hedef kuyruğun var olduğundan emin ol (Yoksa oluşturur)
		_, err := p.ch.QueueDeclare(
			msg.QueueName, // Kuyruk adı
			true,          // Durable: Sunucu kapansa da kuyruk silinmez
			false,         // Auto-delete
			false,         // Exclusive
			false,         // No-wait
			nil,           // Arguments
		)
		if err != nil {
			p.log.Error("Kuyruk deklare edilirken hata olustu (%s): %v", msg.QueueName, err)
			return errors.Wrapf(err, "kuyruk deklare edilemedi: %s", msg.QueueName)
		}

		// 2. Mesajı default exchange ("") kullanarak direkt kuyruğa gönder
		err = p.ch.PublishWithContext(ctx,
			"",            // Default exchange direkt routingKey'e (kuyruk adına) bakar
			msg.QueueName, // Routing key olarak kuyruk adını veriyoruz
			false,         // Mandatory
			false,         // Immediate
			amqp.Publishing{
				ContentType:  contentType,
				Body:         msg.Body,
				DeliveryMode: amqp.Persistent, // Mesajı diske kaydeder (Kalıcılık/Güvenlik için)
			},
		)
		if err != nil {
			p.log.Error("Mesaj gönderilirken hata olustu (%s): %v", msg.QueueName, err)
			return errors.Wrapf(err, "mesaj gonderilemedi: %s", msg.QueueName)
		}
	}

	return nil
}

// Close, kanalı ve bağlantıyı temiz bir şekilde kapatır.
func (p *publisher) Close() error {
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			p.log.Error("Kanal kapatılırken hata: %v", err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			p.log.Error("Bağlantı kapatılırken hata: %v", err)
			return err
		}
	}
	return nil
}
