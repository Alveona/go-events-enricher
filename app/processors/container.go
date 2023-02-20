package processors

import (
	"github.com/Alveona/go-events-enricher/app/metrics"
	"github.com/Alveona/go-events-enricher/app/processors/events"
	"github.com/Alveona/go-events-enricher/app/storages"
)

type Container struct {
	EventsProcessor *events.Processor
}

func New(storageContainer *storages.Container, processorMetrics *metrics.ProcessorContainer) *Container {
	eventsProcessor := events.New(storageContainer.ClickhouseStorage, processorMetrics)
	return &Container{
		EventsProcessor: eventsProcessor,
	}
}
