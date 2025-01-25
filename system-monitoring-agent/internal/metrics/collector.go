// internal/metrics/collector.go
package metrics

import (
	"fmt"
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
	tenantID    string
	tenantMeta  map[string]string
}

func NewMetricsCollector(monitor core.SystemMonitor, cfg *config.Config) *MetricsCollector {
	// Initialize tenant metadata
	tenantMeta := map[string]string{
		"agent_version": "1.0.0",
		"environment":   cfg.Tenant.Environment,
	}
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
		tenantID:    cfg.Tenant.ID,
		tenantMeta:  tenantMeta,
	}
}

func (mc *MetricsCollector) Collect() types.MetricPayload {
	now := time.Now()
	metrics := make(map[string]interface{})
	var collectionErrors []string
	var processData interface{}

	// Get enabled metrics from tenant configuration
	enabledMetrics := make(map[string]bool)
	for _, metric := range mc.config.Tenant.CollectionRules.EnabledMetrics {
		enabledMetrics[metric] = true
	}

	// Collect from each collector
	for _, collector := range mc.collectors {
		// Skip if metric type is not enabled for this tenant
		if !enabledMetrics[collector.Name()] {
			continue
		}

		if data, err := collector.Collect(); err == nil {
			// Special handling for process collector
			if collector.Name() == "processes" {
				processData = data["processes"]
				continue
			}
			// Add other metrics to the metrics map
			for k, v := range data {
				metrics[k] = v
			}
		} else {
			collectionErrors = append(collectionErrors, err.Error())
		}
	}

	// Analyze metrics for threats
	threatIndicators := mc.analyzer.AnalyzeMetrics(metrics)

	// Build tenant metadata
	tenantMeta := map[string]string{
		"agent_version":  "1.0.0",
		"environment":    mc.config.Tenant.Environment,
		"tenant_type":    mc.config.Tenant.Type,
		"tenant_name":    mc.config.Tenant.Name,
		"sample_rate":    fmt.Sprintf("%d", mc.config.Tenant.CollectionRules.SampleRate),
		"retention_days": fmt.Sprintf("%d", mc.config.Tenant.CollectionRules.RetentionDays),
	}

	return types.MetricPayload{
		Timestamp: now.UTC().Format(time.RFC3339),
		Tenant: types.TenantContext{
			ID:       mc.config.Tenant.ID,
			Metadata: tenantMeta,
		},
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
			Processes:        processData.(types.SystemProcessStats),
			Metadata: types.MetadataInfo{
				CollectionDuration: time.Since(now).String(),
				CollectorCount:     len(mc.collectors),
				Errors:             collectionErrors,
				TenantMetadata:     mc.tenantMeta,
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
