package exporter

import "github.com/travism26/shared-monitoring-libs/types"

// MetricsExporter defines the interface for all metric exporters
type MetricsExporter interface {
	Export(data types.MetricPayload) error
	Close() error
}
