package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config, Redis bağlantı ve havuz ayarlarını tutar.
type Config struct {
	Addr         string        // Redis adresi (örn: "localhost:6379")
	Password     string        // Redis şifresi (boş bırakılabilir)
	DB           int           // Kullanılacak veritabanı indeksi (varsayılan: 0)
	PoolSize     int           // Havuzda tutulacak maksimum bağlantı sayısı
	MinIdleConns int           // Havuzda hazır bekleyecek minimum boşta bağlantı sayısı
	DialTimeout  time.Duration // Bağlantı kurulurken beklenecek maksimum süre
}

// NewRedisClient, verilen Config nesnesine göre optimize edilmiş bir Redis istemcisi başlatır.
func NewRedisClient(cfg Config) (*redis.Client, error) {
	// 1. İstemci seçeneklerini Config'den alarak oluştur
	opts := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	// İsteğe bağlı havuz ayarlarını kontrol ederek ekle (0 ise go-redis varsayılanlarını kullanır)
	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	}
	if cfg.MinIdleConns > 0 {
		opts.MinIdleConns = cfg.MinIdleConns
	}
	if cfg.DialTimeout > 0 {
		opts.DialTimeout = cfg.DialTimeout
	}

	rdb := redis.NewClient(opts)

	// 2. Gerçek bir Ping atarak bağlantıyı doğrula
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close() // Bağlantı başarısızsa kaynakları temizle
		return nil, fmt.Errorf("redis baglanti hatasi: %w", err)
	}

	return rdb, nil
}
