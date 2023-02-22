package events

import (
	"context"

	"github.com/Alveona/go-events-enricher/app/entities"
)

type clickhouseStorage interface {
	ProcessInsertEvents(ctx context.Context, events []*entities.EventDTO) error
}

type processorMetrics interface {
	TypesInc(eventType string)
	OSInc(os string)
}

type Processor struct {
	chStorage clickhouseStorage
	metrics   processorMetrics
}

func New(chStorage clickhouseStorage, metrics processorMetrics) *Processor {
	return &Processor{
		chStorage: chStorage,
		metrics:   metrics,
	}
}

func (p *Processor) Process(ctx context.Context, events []*entities.EventDTO) error {
	for _, event := range events {
		p.metrics.TypesInc(event.Event)
		p.metrics.OSInc(event.DeviceOS)
	}
	err := p.chStorage.ProcessInsertEvents(ctx, events)
	if err != nil {
		return err
	}
	return nil
}
