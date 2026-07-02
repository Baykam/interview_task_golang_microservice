package rabbitmq

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ, bağlantıyı ve kanalı sarmalayan ana yapıdır.
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQ, verilen yapılandırma ile bir RabbitMQ bağlantısı ve kanalı başlatır.
func NewRabbitMQ(cfg Config) (*RabbitMQ, error) {
	// 1. Bağlantı ayarlarını yapılandır
	amqpCfg := amqp.Config{
		Dial: amqp.DefaultDial(cfg.DialTimeout),
	}
	if cfg.DialTimeout == 0 {
		amqpCfg.Dial = amqp.DefaultDial(3 * time.Second) // Varsayılan timeout
	}

	// 2. RabbitMQ sunucusuna bağlan
	conn, err := amqp.DialConfig(cfg.URL, amqpCfg)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq baglanti hatasi: %w", err)
	}

	// 3. İletişim için bir kanal aç
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq kanal acma hatasi: %w", err)
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}

func (r *RabbitMQ) Channel() *amqp.Channel {
	return r.ch
}

func (r *RabbitMQ) NewChannel() (*amqp.Channel, error) {
	if r.conn == nil || r.conn.IsClosed() {
		return nil, fmt.Errorf("rabbitmq bağlantısı kapalı, yeni kanal açılamaz")
	}
	return r.conn.Channel()
}
