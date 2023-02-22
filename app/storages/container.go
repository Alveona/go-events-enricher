package storages

import (
	"github.com/Alveona/go-events-enricher/app/config"
	"github.com/Alveona/go-events-enricher/app/storages/events"
	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
)

type Container struct {
	ClickhouseStorage *events.CHStorage
}

func New(eventsCHConn clickhouse.ChConn, storageMetrics events.StorageMetrics, bufferMetrics events.BufferMetrics, config *config.Config) *Container {
	eventsCHStorage := events.NewCHRepo(eventsCHConn, storageMetrics, bufferMetrics, config.Clickhouse.BufferConfig)
	return &Container{
		ClickhouseStorage: eventsCHStorage,
	}
}
