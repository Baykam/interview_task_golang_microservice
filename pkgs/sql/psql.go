package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL Sürücüsü
)

// Config, veritabanı bağlantı ve havuz ayarlarını tutar.
type Config struct {
	DSN             string        // Veritabanı bağlantı adresi (örn: "postgres://user:pass@localhost:5432/db")
	MaxOpenConns    int           // Aynı anda açık olabilecek maksimum bağlantı sayısı
	MaxIdleConns    int           // Havuzda boşta (ready) bekleyecek maksimum bağlantı sayısı
	ConnMaxLifetime time.Duration // Bir bağlantının maksimum yaşam süresi
	ConnMaxIdleTime time.Duration // Boştaki bir bağlantının havuzda bekleyebileceği maksimum süre
}

// NewPostgresDB, dışarıdan verilen Config nesnesine göre bağlantı havuzunu başlatır.
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	// 1. Bağlantıyı aç (Sadece nesne oluşturur)
	db, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open hatasi: %w", err)
	}

	// 2. Bağlantı Havuzu (Connection Pool) Ayarları (Config'den dinamik alınır)
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// 3. Gerçek bir Ping atarak veritabanının orada olduğunu doğrula
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close() // Hata durumunda açılan yapıyı temizle
		return nil, fmt.Errorf("veritabani ping hatasi: %w", err)
	}

	return db, nil
}
