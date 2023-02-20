package entities

import (
	"testing"
	"time"

	"github.com/Alveona/go-events-enricher/app/generated/models"
	"github.com/Alveona/go-events-enricher/app/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var mapEventsTestCases = []struct {
	name   string
	input  *models.ProduceEventsPayload
	output []*EventDTO
}{
	{
		name: "single",
		input: &models.ProduceEventsPayload{
			Payload: ptrOfString("{\"client_time\":\"2020-12-01 23:59:00\",\"device_id\":\"0287D9AA-4ADF-4B37-A60F-3E9E645C821E\",\"device_os\":\"iOS 13.5.1\",\"session\":\"test1\",\"sequence\":1,\"event\":\"app_start\",\"param_int\":123,\"param_str\":\"some text\"}"),
		},
		output: []*EventDTO{
			{
				Event:         "app_start",
				DeviceID:      uuid.MustParse("0287D9AA-4ADF-4B37-A60F-3E9E645C821E"),
				DeviceOS:      "iOS 13.5.1",
				Session:       "test1",
				Sequence:      int64(1),
				ParamInt:      int64(123),
				ParamStr:      "some text",
				ClientTimeRaw: timeMustParse("2020-12-01T23:59:00Z"),
			},
		},
	},
	{
		name: "multiple",
		input: &models.ProduceEventsPayload{
			Payload: ptrOfString("{\"client_time\":\"2020-12-01 23:59:00\",\"device_id\":\"0287D9AA-4ADF-4B37-A60F-3E9E645C821E\",\"device_os\":\"iOS 13.5.1\",\"session\":\"test1\",\"sequence\":1,\"event\":\"app_start\",\"param_int\":123,\"param_str\":\"some text\"}\n{\"client_time\":\"2015-12-01 23:59:00\",\"device_id\":\"5fc39be0-087d-40c0-a769-7bf0d642b790\",\"device_os\":\"Android 8\",\"session\":\"test2\",\"sequence\":1,\"event\":\"app_close\",\"param_int\":1234,\"param_str\":\"some text2\"}"),
		},
		output: []*EventDTO{
			{
				Event:         "app_start",
				DeviceID:      uuid.MustParse("0287D9AA-4ADF-4B37-A60F-3E9E645C821E"),
				DeviceOS:      "iOS 13.5.1",
				Session:       "test1",
				Sequence:      int64(1),
				ParamInt:      int64(123),
				ParamStr:      "some text",
				ClientTimeRaw: timeMustParse("2020-12-01T23:59:00Z"),
			},
			{
				Event:         "app_close",
				DeviceID:      uuid.MustParse("5fc39be0-087d-40c0-a769-7bf0d642b790"),
				DeviceOS:      "Android 8",
				Session:       "test2",
				Sequence:      int64(1),
				ParamInt:      int64(1234),
				ParamStr:      "some text2",
				ClientTimeRaw: timeMustParse("2015-12-01T23:59:00Z"),
			},
		},
	},
}

func ptrOfString(str string) *string {
	return &str
}

func timeMustParse(str string) utils.CustomTime {
	t, _ := time.Parse(time.RFC3339, str)
	return utils.CustomTime{t} //nolint
}

func Test_MapEventsListToDTO(t *testing.T) {
	for _, tt := range mapEventsTestCases {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			res, err := MapEventsListToDTO(testCase.input)
			assert.Nil(t, err)
			assert.Equal(t, testCase.output, res)
		})
	}
}
