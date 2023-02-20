package utils

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Alveona/go-events-enricher/app/constants"
	"github.com/Alveona/go-events-enricher/app/generated/models"
	"github.com/sirupsen/logrus"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = strings.Split(r.RemoteAddr, ":")[0]
	}
	return IPAddress
}

func LoggedError(ctx context.Context, hintCode constants.HintCode, err error) *models.Error {
	logrus.Errorf("code: %d, error: %+v", hintCode, err)

	errorMsg, ok := constants.ErrorMessages[hintCode]
	if !ok {
		errorMsg = constants.ErrorMessages[constants.InternalServerError]
	}

	code := int64(hintCode)
	return &models.Error{
		Code:    &code,
		Message: &errorMsg,
	}
}

type CustomTime struct {
	// We have non-default time format in unserialized json, so we need custom time.
	// Hence, this time is not scannable by sql driver, so we should use both custom and default time in DTO struct.
	time.Time
}

const layout = "2006-01-02 15:04:05"

func (t *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	v, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	*t = CustomTime{v}
	return nil
}
