package handlers

import (
	"context"
	"runtime/debug"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/Alveona/go-events-enricher/app/constants"
	"github.com/Alveona/go-events-enricher/app/entities"
	"github.com/Alveona/go-events-enricher/app/generated/models"
	swg "github.com/Alveona/go-events-enricher/app/generated/restapi/operations"
	"github.com/Alveona/go-events-enricher/app/utils"
)

type eventsProcessor interface {
	Process(ctx context.Context, events []*entities.EventDTO) error
}

type EventsHandler struct {
	processor eventsProcessor
}

func NewEventsHandler(processor eventsProcessor) *EventsHandler {
	return &EventsHandler{
		processor: processor,
	}
}

// ProduceEvents POST /v1/events/produce
func (h *EventsHandler) ProduceEvents(params swg.ProduceEventsParams) middleware.Responder {
	events, err := entities.MapEventsListToDTO(params.Payload)
	if err != nil {
		return swg.NewProduceEventsUnprocessableEntity().WithPayload(utils.LoggedError(
			params.HTTPRequest.Context(),
			constants.BadRequestError,
			err,
		))
	}
	go func() {
		entities.EnrichEvents(params.HTTPRequest, events)
		processCtx := context.Background()
		err := h.processor.Process(processCtx, events)
		if err != nil {
			logrus.Errorf("Error inserting events: %+v, %s", events, err.Error())
		}
		if r := recover(); r != nil {
			logrus.Errorf("Panic recovered: %+v, %s", r, debug.Stack())
		}
	}()

	return swg.NewProduceEventsOK().WithPayload(&models.ProduceEventsResponse{
		Status: models.ProduceEventsResponseStatusOK,
	})
}
