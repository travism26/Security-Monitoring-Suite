// internal/metrics/collector_interface.go
package metrics

type MetricCollector interface {
	Collect() (map[string]interface{}, error)
	Name() string
}
