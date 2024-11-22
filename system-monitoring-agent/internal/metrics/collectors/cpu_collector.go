// internal/metrics/collectors/cpu_collector.go
package collectors

import "github.com/travism26/system-monitoring-agent/internal/core"

type CPUCollector struct {
	monitor core.Monitor
}

func NewCPUCollector(monitor core.Monitor) *CPUCollector {
	return &CPUCollector{monitor: monitor}
}

func (c *CPUCollector) Name() string {
	return "cpu"
}

func (c *CPUCollector) Collect() (map[string]interface{}, error) {
	cpuUsage, err := c.monitor.GetCPUUsage()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"cpu_usage": cpuUsage,
	}, nil
}
