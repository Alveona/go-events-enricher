package entities

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var enrichEventsTestCases = []struct {
	name  string
	input []*EventDTO
}{
	{
		name:  "default",
		input: []*EventDTO{{}},
	},
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func Test_EnrichEvents(t *testing.T) {
	for _, tt := range enrichEventsTestCases {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			start := time.Now()
			EnrichEvents(&http.Request{RemoteAddr: "8.8.8.8"}, testCase.input)
			end := time.Now()
			for _, event := range testCase.input {
				assert.Equal(t, event.IP, "8.8.8.8")
				assert.True(t, inTimeSpan(start, end, event.ServerTime))
			}
		})
	}
}
