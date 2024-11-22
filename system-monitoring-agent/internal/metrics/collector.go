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

func NewMetricsCollector(m core.Monitor, cfg *config.Config) *MetricsCollector {
	collectors := []MetricCollector{
		collectors.NewCPUCollector(m),
		collectors.NewMemoryCollector(m),
		collectors.NewDiskCollector(m),
		collectors.NewNetworkCollector(m),
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

	// Add metadata
	metrics["timestamp"] = now.Unix()
	metrics["timestamp_utc"] = now.UTC().Format(time.RFC3339)
	metrics["host_info"] = map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}

	// Collect from each collector
	for _, collector := range mc.collectors {
		if data, err := collector.Collect(); err == nil {
			for k, v := range data {
				metrics[k] = v
			}
		}
	}

	// Analyze metrics for threats
	if mc.analyzer != nil {
		indicators := mc.analyzer.AnalyzeMetrics(metrics)
		if len(indicators) > 0 {
			metrics["threat_indicators"] = indicators
		}
	}

	mc.lastCheck = now
	return metrics
}
