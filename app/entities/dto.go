package entities

import (
	"net/http"
	"time"

	"github.com/Alveona/go-events-enricher/app/utils"
	"github.com/google/uuid"
)

type EventDTO struct {
	Event         string           `json:"event" db:"event"`
	DeviceID      uuid.UUID        `json:"device_id" db:"device_id"`
	DeviceOS      string           `json:"device_os" db:"device_os"`
	Session       string           `json:"session" db:"session"`
	Sequence      int64            `json:"sequence" db:"sequence"`
	ParamInt      int64            `json:"param_int" db:"param_int"`
	ParamStr      string           `json:"param_str" db:"param_str"`
	ClientTimeRaw utils.CustomTime `json:"client_time" db:"-"` // See description of CustomTime
	ClientTime    time.Time        `json:"-" db:"client_time"`
	ServerTime    time.Time        `json:"-" db:"server_time"`
	IP            string           `json:"-" db:"ip"`
}

func (e *EventDTO) enrichEventDTO(ip string, serverTime time.Time) {
	e.ServerTime = serverTime
	e.IP = ip
	e.ClientTime = e.ClientTimeRaw.Time
}

// EnrichEvents enriches each event struct with request params and current time.
func EnrichEvents(req *http.Request, events []*EventDTO) {
	serverTime := time.Now()
	ip := utils.ReadUserIP(req)
	for _, event := range events {
		event.enrichEventDTO(ip, serverTime)
	}
}
