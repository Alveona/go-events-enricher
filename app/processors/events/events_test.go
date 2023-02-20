package events

import (
	"context"
	"testing"

	"github.com/Alveona/go-events-enricher/app/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	syncerSuite struct {
		suite.Suite
		ctx    context.Context
		cancel context.CancelFunc
	}
)

func Test_SyncerSuite(t *testing.T) {
	suite.Run(t, new(syncerSuite))
}

func (s *syncerSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
}

func (s *syncerSuite) TearDownTest() {
	s.cancel()
}

func (s *syncerSuite) TestIncVendors() {
	s.T().Parallel()
	type testCase struct { //nolint:maligned
		name        string
		mockStorage func(mock *MockclickhouseStorage)
		mockMetric  func(mock *MockprocessorMetrics)
		input       []*entities.EventDTO
		wantErr     error
	}

	tests := []testCase{
		{
			name: "success",
			mockStorage: func(storage *MockclickhouseStorage) {
				storage.EXPECT().InsertEvents(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			mockMetric: func(metric *MockprocessorMetrics) {
				metric.EXPECT().TypesInc("app_start").Times(1)
				metric.EXPECT().OSInc("iOS").Times(1)
			},
			input: []*entities.EventDTO{{
				Event:    "app_start",
				DeviceOS: "iOS",
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tc := tt
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockclickhouseStorage(ctrl)
			metric := NewMockprocessorMetrics(ctrl)
			syncer := New(storage, metric)
			tc.mockStorage(storage)
			tc.mockMetric(metric)

			err := syncer.Process(s.ctx, tc.input)
			assert.Equal(t, tc.wantErr, err)
		})

	}
}
