package config

import (
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"

	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
)

// Config service configuration
type Config struct {
	Clickhouse *clickhouse.Config
}

func InitConfig(prefix string) (*Config, error) {
	config := &Config{}
	if err := envconfig.InitWithPrefix(config, prefix); err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}
