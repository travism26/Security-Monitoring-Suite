package domain

import "time"

type MetricType string

const (
	MetricTypeCPU     MetricType = "CPU"
	MetricTypeMemory  MetricType = "MEMORY"
	MetricTypeDisk    MetricType = "DISK"
	MetricTypeNetwork MetricType = "NETWORK"
)

type Metric struct {
	ID        string     `json:"id"`
	Type      MetricType `json:"type"`
	Host      string     `json:"host"`
	Value     float64    `json:"value"`
	Unit      string     `json:"unit"`
	Timestamp time.Time  `json:"timestamp"`

	// Additional labels/tags for the metric
	Labels map[string]string `json:"labels,omitempty"`

	// Raw metric data
	RawData map[string]interface{} `json:"raw_data,omitempty"`
}

// MetricThreshold defines alerting thresholds for metrics
type MetricThreshold struct {
	Type     MetricType    `json:"type"`
	Warning  float64       `json:"warning"`
	Critical float64       `json:"critical"`
	Duration time.Duration `json:"duration"` // Duration the threshold must be exceeded
}

// MetricRepository defines the interface for metric storage operations
type MetricRepository interface {
	Store(metric *Metric) error
	FindByID(id string) (*Metric, error)
	List(limit, offset int) ([]*Metric, error)
	FindByTimeRange(start, end time.Time) ([]*Metric, error)
	FindByHost(host string, limit, offset int) ([]*Metric, error)
	GetAggregates(metricType MetricType, interval time.Duration) ([]*MetricAggregate, error)
}

// MetricAggregate represents aggregated metric data
type MetricAggregate struct {
	Timestamp time.Time `json:"timestamp"`
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	Avg       float64   `json:"avg"`
	Count     int64     `json:"count"`
}
