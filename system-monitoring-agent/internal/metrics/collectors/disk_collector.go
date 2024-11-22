// internal/metrics/collectors/disk_collector.go
package collectors

import "github.com/travism26/system-monitoring-agent/internal/core"

type DiskCollector struct {
	monitor core.Monitor
}

func NewDiskCollector(monitor core.Monitor) *DiskCollector {
	return &DiskCollector{monitor: monitor}
}

func (c *DiskCollector) Name() string {
	return "disk"
}

func (c *DiskCollector) Collect() (map[string]interface{}, error) {
	diskStats, err := c.monitor.GetDiskUsage()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"disk": map[string]interface{}{
			"total":         diskStats.Total,
			"used":          diskStats.Used,
			"free":          diskStats.Free,
			"usage_percent": float64(diskStats.Used) / float64(diskStats.Total) * 100,
		},
	}, nil
}
