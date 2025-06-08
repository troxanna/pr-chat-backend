package config

import (
	"fmt"
	"time"

	env "github.com/caarlos0/env/v6"
)

type Config struct {

	// Prometheus struct {
	// 	ListenAddress string `env:"PROMETHEUS_LISTEN_ADDRESS" envDefault:":2112"`
	// }
	// Redis struct {
	// 	Username           string `env:"REDIS_USERNAME"`
	// 	Password           string `env:"REDIS_PASSWORD" json:"-"`
	// 	Address            string `env:"REDIS_ADDRESS,notEmpty"`
	// 	DatabaseNumber     int    `env:"REDIS_DATABASE_NUMBER"`
	// 	PoolSize           int    `env:"REDIS_POOL_SIZE" envDefault:"5"`
	// 	MinIdleConnections int    `env:"REDIS_MIN_IDLE_CONNECTIONS" envDefault:"5"`
	// 	MaxIdleConnections int    `env:"REDIS_MAX_IDLE_CONNECTIONS" envDefault:"10"`
	// }

	HTTP struct {
		ListenAddress        string        `env:"HTTP_LISTEN_ADDRESS,notEmpty"`
		ListenAddressPrivate string        `env:"HTTP_LISTEN_ADDRESS_PRIVATE" envDefault:"0.0.0.0:8080"`
		WriteTimeout         time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"90s"`
		ReadTimeout          time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"90s"`
		IdleTimeout          time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		ShutdownTimeout      time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"30s"`
	}

	// JWT struct {
	// 	PrivateKey      string        `env:"JWT_PRIVATE_KEY,notEmpty" json:"-"` // Hide in zap logs
	// 	PublicKey       string        `env:"JWT_PUBLIC_KEY,notEmpty"`
	// 	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"5m"`
	// 	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"168h"`
	// }

	// Postgres struct {
	// 	URL             string        `env:"PG_URL, notEmpty"`
	// 	DSN             string        `env:"PG_DSN,notEmpty" json:"-"` // Hide in zap logs
	// 	MaxIdleConns    int           `env:"PG_MAX_IDLE_CONNS" envDefault:"15"`
	// 	MaxOpenConns    int           `env:"PG_MAX_OPEN_CONNS" envDefault:"15"`
	// 	ConnMaxLifetime time.Duration `env:"PG_CONN_MAX_LIFETIME" envDefault:"5m"`
	// }

	Log struct {
		FieldMaxLen int `env:"LOG_FIELD_MAX_LEN" envDefault:"2000"`

		SensitiveDataMasker struct {
			Enabled bool `env:"LOG_SENSITIVE_DATA_MASKER_ENABLED" envDefault:"true"`
		}
	}

	S3Storage struct {
		Endpoint string `env:"S3_STORAGE_ENDPOINT" envDefault:"https://storage.yandexcloud.net"`
		// SigningRegion          string `env:"S3_STORAGE_SIGNING_REGION" envDefault:"ru-central1"`
		// AccessKeyID            string `env:"S3_STORAGE_ACCESS_KEY_ID,notEmpty"`
		// SecretAccessKey        string `env:"S3_STORAGE_SECRET_ACCESS_KEY,notEmpty"`
		// Session                string `env:"S3_STORAGE_SESSION"`
		// BucketNameFile         string `env:"S3_STORAGE_BUCKET_NAME_FILE" envDefault:"file"`
	}

	Telegram struct {
		BotToken   string `env:"TELEGRAM_BOT_TOKEN,notEmpty"`
		WebhookUrl string `env:"TG_WEBHOOK_URL" envDefault:"-"`
		Port       string `env:"TG_PORT" envDefault:"8697"`
		Cert       string `env:"TG_CERTIFICATE,notEmpty"`
	}
}

func Load() (Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}

	return config, nil
}
