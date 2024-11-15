// internal/metrics/collector.go
package metrics

import (
	"runtime"
	"time"

	"github.com/travism26/system-monitoring-agent/internal/config"
	"github.com/travism26/system-monitoring-agent/internal/core"
)

type MetricsCollector struct {
	monitor     core.Monitor
	config      *config.Config
	lastNetwork map[string]core.NetworkStats
	lastCheck   time.Time
}

func NewMetricsCollector(m core.Monitor, cfg *config.Config) *MetricsCollector {
	return &MetricsCollector{
		monitor:     m,
		config:      cfg,
		lastNetwork: make(map[string]core.NetworkStats),
		lastCheck:   time.Now(),
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

	// CPU metrics
	cpuUsage, err := mc.monitor.GetCPUUsage()
	if err == nil {
		metrics["cpu_usage"] = cpuUsage
	}

	// Memory metrics
	memUsage, err := mc.monitor.GetMemoryUsage()
	if err == nil {
		metrics["memory_usage"] = memUsage
	}

	totalMem, err := mc.monitor.GetTotalMemory()
	if err == nil {
		metrics["total_memory"] = totalMem
		if memUsage > 0 {
			metrics["memory_usage_percent"] = float64(memUsage) / float64(totalMem) * 100
		}
	}

	// Disk metrics
	diskStats, err := mc.monitor.GetDiskUsage()
	if err == nil {
		metrics["disk"] = map[string]interface{}{
			"total":         diskStats.Total,
			"used":          diskStats.Used,
			"free":          diskStats.Free,
			"usage_percent": float64(diskStats.Used) / float64(diskStats.Total) * 100,
		}
	}

	// Network metrics
	netStats, err := mc.monitor.GetNetworkStats()
	if err == nil {
		timeDiff := now.Sub(mc.lastCheck).Seconds()
		networkMetrics := map[string]interface{}{
			"bytes_sent":     netStats.BytesSent,
			"bytes_received": netStats.BytesReceived,
		}

		// Calculate transfer rates if we have previous measurements
		if lastStats, exists := mc.lastNetwork[""]; exists && timeDiff > 0 {
			bytesSentRate := float64(netStats.BytesSent-lastStats.BytesSent) / timeDiff
			bytesReceivedRate := float64(netStats.BytesReceived-lastStats.BytesReceived) / timeDiff

			networkMetrics["bytes_sent_per_second"] = bytesSentRate
			networkMetrics["bytes_received_per_second"] = bytesReceivedRate
		}

		metrics["network"] = networkMetrics
		mc.lastNetwork[""] = netStats
	}

	mc.lastCheck = now
	return metrics
}
