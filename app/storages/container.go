package storages

import (
	"github.com/Alveona/go-events-enricher/app/config"
	"github.com/Alveona/go-events-enricher/app/storages/events"
	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
)

type Container struct {
	ClickhouseStorage *events.CHStorage
}

func New(eventsCHConn clickhouse.ChConn, metrics events.Metrics, config *config.Config) *Container {
	eventsCHStorage := events.NewCHRepo(eventsCHConn, metrics)
	return &Container{
		ClickhouseStorage: eventsCHStorage,
	}
}
