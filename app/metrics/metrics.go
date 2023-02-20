package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Container
type Container interface {
	Collectors() []prometheus.Collector
}

// ProcessorContainer processor metrics container
type ProcessorContainer struct {
	eventsByTypes *prometheus.CounterVec
	eventsByOS    *prometheus.CounterVec
}

// NewProcessorContainer creates new service metrics container
func NewProcessorContainer(appName string) *ProcessorContainer {
	return &ProcessorContainer{
		eventsByTypes: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: appName,
			Name:      "added_by_types",
			Help:      "added events by types",
		}, []string{"type"}),
		eventsByOS: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: appName,
			Name:      "added_by_os",
			Help:      "added events by types",
		}, []string{"os"}),
	}
}

// Build merges all metrics container to a common one
func GetMetrics(containers ...Container) []prometheus.Collector {
	collectors := make([]prometheus.Collector, 0)
	for _, container := range containers {
		collectors = append(collectors, container.Collectors()...)
	}

	return collectors
}

func (a *ProcessorContainer) TypesInc(eventType string) {
	a.eventsByTypes.With(
		prometheus.Labels{"type": eventType},
	).Inc()
}

func (a *ProcessorContainer) OSInc(os string) {
	a.eventsByOS.With(
		prometheus.Labels{"os": os},
	).Inc()
}

func (a *ProcessorContainer) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		a.eventsByTypes,
		a.eventsByOS,
	}
}

// StorageContainer service metrics container
type StorageContainer struct {
	storageExec *prometheus.SummaryVec
}

// NewStorageContainer creates new service metrics container
func NewStorageContainer(appName string) *StorageContainer {
	return &StorageContainer{
		storageExec: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: appName,
			Name:      "storage_query_duration",
			Help:      "storage's execute query duration",
		}, []string{"query"}),
	}
}

func (c *StorageContainer) QueryDuration(duration time.Duration, labels ...string) {
	c.storageExec.WithLabelValues(labels...).Observe(
		float64(duration.Milliseconds()),
	)
}

// Collectors возвращает все коллекторы метрик.
func (c *StorageContainer) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.storageExec,
	}
}
