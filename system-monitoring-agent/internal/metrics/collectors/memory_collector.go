// internal/metrics/collectors/memory_collector.go
package collectors

import "github.com/travism26/system-monitoring-agent/internal/core"

type MemoryCollector struct {
	monitor core.MemoryMonitor
}

func NewMemoryCollector(monitor core.MemoryMonitor) *MemoryCollector {
	return &MemoryCollector{monitor: monitor}
}

func (c *MemoryCollector) Name() string {
	return "memory"
}

func (c *MemoryCollector) Collect() (map[string]interface{}, error) {
	memUsage, err := c.monitor.GetMemoryUsage()
	if err != nil {
		return nil, err
	}

	totalMem, err := c.monitor.GetTotalMemory()
	if err != nil {
		return nil, err
	}

	metrics := map[string]interface{}{
		"memory_usage": memUsage,
		"total_memory": totalMem,
	}

	if memUsage > 0 {
		metrics["memory_usage_percent"] = float64(memUsage) / float64(totalMem) * 100
	}

	return metrics, nil
}
