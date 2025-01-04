package metrics

import (
	"testing"

	"github.com/travism26/system-monitoring-agent/internal/config"
	"github.com/travism26/system-monitoring-agent/internal/core"
)

// MockSystemMonitor implements core.SystemMonitor for testing
type MockSystemMonitor struct{}

func (m *MockSystemMonitor) GetCPUUsage() (float64, error) {
	return 25.0, nil // Return a sample CPU usage
}

func (m *MockSystemMonitor) GetMemoryUsage() (uint64, error) {
	return 1024 * 1024 * 512, nil // Return 512MB memory usage
}

func (m *MockSystemMonitor) GetProcesses() ([]core.ProcessInfo, error) {
	return []core.ProcessInfo{
		{
			PID:         1234,
			Name:        "test-process",
			CPUPercent:  1.5,
			MemoryUsage: 1024 * 1024 * 10, // 10MB
			Status:      "running",
		},
	}, nil
}

func (m *MockSystemMonitor) GetNetworkStats() (core.NetworkStats, error) {
	return core.NetworkStats{}, nil
}

func (m *MockSystemMonitor) GetDiskUsage() (core.DiskStats, error) {
	return core.DiskStats{
		Total: 1024 * 1024 * 1024 * 100, // 100GB
		Used:  1024 * 1024 * 1024 * 50,  // 50GB
		Free:  1024 * 1024 * 1024 * 50,  // 50GB
	}, nil
}

func (m *MockSystemMonitor) GetTotalMemory() (uint64, error) {
	return 0, nil
}

func TestNewMetricsCollector(t *testing.T) {
	// Create mock dependencies
	mockMonitor := &MockSystemMonitor{}
	cfg := &config.Config{
		Interval: 5,
	}

	// Create collector
	collector := NewMetricsCollector(mockMonitor, cfg)

	// Verify collector initialization
	if collector == nil {
		t.Error("Expected non-nil collector, got nil")
	}
	if len(collector.collectors) != 5 {
		t.Errorf("Expected 5 collectors, got %d", len(collector.collectors))
	}
	if collector.config != cfg {
		t.Error("Config not set correctly")
	}
	if collector.analyzer == nil {
		t.Error("Expected non-nil analyzer, got nil")
	}
}

func TestCollect(t *testing.T) {
	mockMonitor := &MockSystemMonitor{}
	cfg := &config.Config{Interval: 5}
	collector := NewMetricsCollector(mockMonitor, cfg)

	// Test successful collection
	metrics := collector.Collect()

	// Verify expected metrics are present
	expectedMetrics := []string{"cpu", "memory", "disk", "network"}
	for _, metric := range expectedMetrics {
		if _, exists := metrics.Data.Metrics[metric]; !exists {
			t.Errorf("Expected %s metrics", metric)
		}
	}

	// Verify processes data
	if len(metrics.Data.Processes.ProcessList) == 0 {
		t.Error("Expected process metrics")
	}

	// Test error handling
	// TODO: Add error scenarios once we have error injection in the mock
}

func TestHostname(t *testing.T) {
	host := hostname()
	if host == "" {
		t.Error("Expected non-empty hostname")
	}
}

func TestErrorHandling(t *testing.T) {
	// TODO: Implement error handling tests
	// This will require modifying the mock to support error injection
}
