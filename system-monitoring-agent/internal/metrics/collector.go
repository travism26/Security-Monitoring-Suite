// internal/metrics/collector.go
package metrics

import "github.com/travism26/system-monitoring-agent/internal/core"

type MetricsCollector struct {
	monitor core.Monitor
}

func NewMetricsCollector(m core.Monitor) *MetricsCollector {
	return &MetricsCollector{
		monitor: m,
	}
}

func (mc *MetricsCollector) Collect() map[string]interface{} {
	metrics := make(map[string]interface{})

	cpuUsage, err := mc.monitor.GetCPUUsage()
	if err == nil {
		metrics["cpu_usage"] = cpuUsage
	}

	memUsage, err := mc.monitor.GetMemoryUsage()
	if err == nil {
		metrics["memory_usage"] = memUsage
	}

	totalMem, err := mc.monitor.GetTotalMemory()
	if err == nil {
		metrics["total_memory"] = totalMem
	}

	return metrics
}
