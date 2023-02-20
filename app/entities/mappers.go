package entities

import (
	"encoding/json"
	"strings"

	"github.com/Alveona/go-events-enricher/app/generated/models"
)

func MapEventsListToDTO(payload *models.ProduceEventsPayload) ([]*EventDTO, error) {
	list := strings.Split(*payload.Payload, "\n")
	events := make([]*EventDTO, 0, len(list))
	for _, obj := range list {
		event := &EventDTO{}
		err := json.Unmarshal([]byte(obj), event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
