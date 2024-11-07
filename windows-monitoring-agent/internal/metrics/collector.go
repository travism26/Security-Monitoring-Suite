// internal/metrics/collector.go
package metrics

import "github.com/travism26/windows-monitoring-agent/internal/monitor"

type MetricsCollector struct {
	monitor *monitor.Monitor
}

func NewMetricsCollector(m *monitor.Monitor) *MetricsCollector {
	return &MetricsCollector{
		monitor: m,
	}
}

func (mc *MetricsCollector) Collect() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Collect CPU usage
	cpuUsage, err := mc.monitor.GetCPUUsage()
	if err == nil {
		metrics["cpu_usage"] = cpuUsage
	}

	// Collect Memory usage
	memUsage, err := mc.monitor.GetMemoryUsage()
	if err == nil {
		metrics["memory_usage"] = memUsage
	}

	return metrics
}
