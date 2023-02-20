package app

import (
	swgservice "github.com/Alveona/go-base-service"
	"github.com/Alveona/go-events-enricher/app/config"
	"github.com/Alveona/go-events-enricher/app/generated/restapi/operations"
	"github.com/Alveona/go-events-enricher/app/handlers"
	"github.com/Alveona/go-events-enricher/app/metrics"
	"github.com/Alveona/go-events-enricher/app/processors"
	"github.com/Alveona/go-events-enricher/app/storages"
	"github.com/Alveona/go-events-enricher/app/system"
	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
	"github.com/sirupsen/logrus"
)

// EventsEnricher struct
type EventsEnricher struct {
	swgservice.BaseService
	system     *system.Container
	storages   *storages.Container
	processors *processors.Container
	config     *config.Config
}

// New initialize a new service
func New(baseService swgservice.BaseService) swgservice.ServiceImplementation {
	return &EventsEnricher{
		BaseService: baseService,
	}
}

// ConfigureService service configuration
func (svc *EventsEnricher) ConfigureService() error {
	cfg, err := config.InitConfig(svc.Name())
	if err != nil {
		return err
	}
	svc.config = cfg

	clickhouseDB, err := clickhouse.New(cfg.Clickhouse)
	if err != nil {
		return err
	}

	storageMetrics := metrics.NewStorageContainer(svc.Name())
	svc.system = system.New(clickhouseDB)

	eventsProcessorMetrics := metrics.NewProcessorContainer(svc.Name())
	svc.Metrics().MustRegister(metrics.GetMetrics(
		storageMetrics,
		eventsProcessorMetrics,
	)...)

	svc.storages = storages.New(clickhouseDB, storageMetrics, cfg)
	svc.processors = processors.New(svc.storages, eventsProcessorMetrics)
	return nil
}

// SetupSwaggerHandlers service handlers
func (svc *EventsEnricher) SetupSwaggerHandlers(iapi interface{}) {
	api, ok := iapi.(*operations.EventsEnricherAPI)
	if !ok {
		logrus.Error("iapi is not a operations.EventsEnricherAPI type")
		return
	}
	eventsHandlers := handlers.NewEventsHandler(svc.processors.EventsProcessor)
	api.ProduceEventsHandler = operations.ProduceEventsHandlerFunc(eventsHandlers.ProduceEvents)
}

// HealthCheckers list of health check
func (svc *EventsEnricher) HealthCheckers() []swgservice.Checker {
	checkers := make([]swgservice.Checker, 0)

	// service health check
	checkers = append(checkers, func() swgservice.CheckerResult {
		return swgservice.CheckerResult{
			Service: "System",
			Status:  true,
		}
	})
	// Clickhouse(master) health check
	checkers = append(checkers, func() swgservice.CheckerResult {
		return swgservice.CheckerResult{
			Service: "Clickhouse",
			Status:  svc.system.Clickhouse.PingMaster(),
		}
	})

	return checkers
}

func (svc *EventsEnricher) OnShutdown() {
	logrus.Warn("shutting down clickhouse...")
	if err := svc.system.Clickhouse.Master().Close(); err != nil {
		logrus.Errorf("close clickhouse connection: %+v", err)
	}
}
