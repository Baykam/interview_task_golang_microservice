package rabbitmq

import "time"

// Config, RabbitMQ'nun genel bağlantı ve başlatma ayarlarını tutar.
type Config struct {
	URL         string        `mapstructure:"url"`          // Örn: "amqp://guest:guest@localhost:5672/"
	DialTimeout time.Duration `mapstructure:"dial_timeout"` // Bağlantı kurulurken beklenecek maksimum süre
	InitQueues  bool          `mapstructure:"initQueues"`   // Uygulama başlarken kuyrukları otomatik oluşturmak için (Kafka'daki InitTopics gibi)
}

// QueueConfig, dinamik olarak oluşturulacak veya dinlenecek kuyrukların ayarlarını tutar.
type QueueConfig struct {
	QueueName  string `mapstructure:"queueName"`  // Kuyruğun benzersiz adı
	Durable    bool   `mapstructure:"durable"`    // Sunucu kapansa da kuyruk silinsin mi? (Genelde true)
	AutoDelete bool   `mapstructure:"autoDelete"` // Tüketici bağlantısı bittiğinde kuyruk silinsin mi? (Genelde false)
	Exclusive  bool   `mapstructure:"exclusive"`  // Sadece bu bağlantıya özel mi olsun? (Genelde false)
}

// ExchangeConfig, eğer gelişmiş pub/sub mimarisi kullanacaksan exchange ayarlarını tutar.
type ExchangeConfig struct {
	ExchangeName string `mapstructure:"exchangeName"` // Exchange adı
	Type         string `mapstructure:"type"`         // direct, topic, fanout, headers
	Durable      bool   `mapstructure:"durable"`
	AutoDelete   bool   `mapstructure:"autoDelete"`
}
