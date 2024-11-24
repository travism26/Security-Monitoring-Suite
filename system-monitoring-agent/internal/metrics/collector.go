// internal/metrics/collector.go
package metrics

import (
	"runtime"
	"time"

	"github.com/travism26/system-monitoring-agent/internal/config"
	"github.com/travism26/system-monitoring-agent/internal/core"
	"github.com/travism26/system-monitoring-agent/internal/metrics/collectors"
	"github.com/travism26/system-monitoring-agent/internal/threat"
)

type MetricsCollector struct {
	collectors  []MetricCollector
	config      *config.Config
	lastNetwork map[string]core.NetworkStats
	lastCheck   time.Time
	analyzer    *threat.Analyzer
}

func NewMetricsCollector(monitor core.SystemMonitor, cfg *config.Config) *MetricsCollector {
	collectors := []MetricCollector{
		collectors.NewCPUCollector(monitor),
		collectors.NewMemoryCollector(monitor),
		collectors.NewDiskCollector(monitor),
		collectors.NewNetworkCollector(monitor),
		collectors.NewProcessCollector(monitor),
	}

	return &MetricsCollector{
		collectors:  collectors,
		config:      cfg,
		lastNetwork: make(map[string]core.NetworkStats),
		lastCheck:   time.Now(),
		analyzer:    threat.NewAnalyzer(),
	}
}

func (mc *MetricsCollector) Collect() map[string]interface{} {
	now := time.Now()
	metrics := make(map[string]interface{})

	// Collect from each collector
	for _, collector := range mc.collectors {
		if data, err := collector.Collect(); err == nil {
			for k, v := range data {
				metrics[k] = v
			}
		}
	}

	// Structure the response according to the API requirements
	return map[string]interface{}{
		"timestamp": now.UTC().Format(time.RFC3339),
		"data": map[string]interface{}{
			"host_info": map[string]string{
				"os":   runtime.GOOS,
				"arch": runtime.GOARCH,
			},
			"metrics": metrics,
		},
	}
}
