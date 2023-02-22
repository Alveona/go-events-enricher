package clickhouse

import "time"

type Config struct {
	DSN                string
	DialTimeout        time.Duration `envconfig:"default=1s"`
	MaxOpenConns       int           `envconfig:"default=10"`
	MaxIdleConns       int           `envconfig:"default=5"`
	ConnMaxLifetime    time.Duration `envconfig:"default=60m"`
	MaxExecutionTime   int           `envconfig:"default=60"`
	HealthCheckTimeout time.Duration `envconfig:"default=1000ms"`
	BufferConfig       BufferConfig
}

type BufferConfig struct {
	DequeueDelay      time.Duration `envconfig:"default=1s"`
	BufferWaitRetries time.Duration `envconfig:"default=100ms"`
	BufferSize        uint64        `envconfig:"default=30"`
	WithBuffer        bool          `envconfig:"default=true"`
	Ratelimit         int           `envconfig:"default=1024"`
	DequeueTimeout    time.Duration `envconfig:"default=3m"`
}
