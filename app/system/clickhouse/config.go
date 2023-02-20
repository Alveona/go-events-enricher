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
}
