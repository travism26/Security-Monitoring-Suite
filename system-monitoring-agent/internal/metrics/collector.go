// internal/metrics/collector.go
package metrics

import (
	"os"
	"runtime"
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
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

func (mc *MetricsCollector) Collect() types.MetricPayload {
	now := time.Now()
	metrics := make(map[string]interface{})
	var collectionErrors []string

	// Collect from each collector
	for _, collector := range mc.collectors {
		if data, err := collector.Collect(); err == nil {
			for k, v := range data {
				metrics[k] = v
			}
		} else {
			collectionErrors = append(collectionErrors, err.Error())
		}
	}

	// Analyze metrics for threats
	threatIndicators := mc.analyzer.AnalyzeMetrics(metrics)

	// Structure the response according to the API requirements
	return types.MetricPayload{
		Timestamp: now.UTC().Format(time.RFC3339),
		Data: types.MetricData{
			HostInfo: types.HostInfo{
				OS:        runtime.GOOS,
				Arch:      runtime.GOARCH,
				Hostname:  hostname(),
				CPUCores:  runtime.NumCPU(),
				GoVersion: runtime.Version(),
			},
			Metrics:          metrics,
			ThreatIndicators: threatIndicators,
			Metadata: types.MetadataInfo{
				CollectionDuration: time.Since(now).String(),
				CollectorCount:     len(mc.collectors),
				Errors:             collectionErrors,
			},
		},
	}
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}
