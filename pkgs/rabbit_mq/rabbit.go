package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config, RabbitMQ bağlantı ayarlarını tutar.
type Config struct {
	URL         string        // RabbitMQ bağlantı adresi (örn: "amqp://guest:guest@localhost:5672/")
	DialTimeout time.Duration // Bağlantı kurulurken beklenecek maksimum süre
}

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

// Close, hem kanalı hem de bağlantıyı güvenli bir şekilde kapatır.
func (r *RabbitMQ) Close() {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// Publish (Producer), belirtilen kuyruğa bir mesaj gönderir.
func (r *RabbitMQ) Publish(ctx context.Context, queueName string, body []byte) error {
	// Kuyruğun var olduğundan emin ol (Yoksa oluşturur)
	_, err := r.ch.QueueDeclare(
		queueName, // Kuyruk adı
		true,      // Durable: RabbitMQ kapansa da kuyruk silinmez
		false,     // Auto-delete: Tüketici bittiğinde kuyruk silinmez
		false,     // Exclusive: Sadece bu bağlantıya özel değil
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("kuyruk deklarasyon hatasi: %w", err)
	}

	// Mesajı kuyruğa gönder
	err = r.ch.PublishWithContext(ctx,
		"",        // Exchange (Boş bırakılırsa varsayılan exchange kullanılır)
		queueName, // Routing key (Kuyruk adı)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType:  "text/plain", // Veya "application/json"
			Body:         body,
			DeliveryMode: amqp.Persistent, // Mesajı diske kaydet (Kalıcılık için)
		},
	)
	if err != nil {
		return fmt.Errorf("mesaj gonderim hatasi: %w", err)
	}

	return nil
}

// Consume (Consumer), belirtilen kuyruktan gelen mesajları dinler ve bir callback fonksiyonuna aktarır.
func (r *RabbitMQ) Consume(queueName string, handler func(body []byte)) error {
	// Kuyruğun var olduğundan emin ol
	_, err := r.ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("kuyruk deklarasyon hatasi: %w", err)
	}

	// Kuyruktan mesaj akışını (channel) başlat
	msgs, err := r.ch.Consume(
		queueName, // Kuyruk adı
		"",        // Consumer tag (Boş bırakılırsa otomatik üretilir)
		true,      // Auto-Ack: Mesaj alındığı an otomatik onaylanır (Basitlik için true)
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		return fmt.Errorf("kuyruk dinleme hatasi: %w", err)
	}

	// Ayrı bir goroutine içinde mesajları sürekli dinle ve işle
	go func() {
		for d := range msgs {
			handler(d.Body) // Gelen mesajı dışarıdan verilen fonksiyona pasla
		}
	}()

	return nil
}
