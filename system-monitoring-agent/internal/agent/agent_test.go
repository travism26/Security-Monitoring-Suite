package agent

import (
	"testing"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
	"github.com/travism26/system-monitoring-agent/internal/config"
	"github.com/travism26/system-monitoring-agent/internal/core"
	"github.com/travism26/system-monitoring-agent/internal/exporter"
	"github.com/travism26/system-monitoring-agent/internal/metrics"
)

type MockCPU struct{}

// Implement necessary methods for MockCPU to satisfy the interface
// For example, if the interface has a method called GetUsage, implement it:
func (m *MockCPU) GetUsage() float64 {
	return 0.0 // Return a mock value
}

// Add the missing Percent method with the correct signature
func (m *MockCPU) Percent(value float64, flag bool) ([]float64, error) {
	return []float64{0.0}, nil // Return a mock value and nil error
}

type MockDisk struct{}

type MockNetwork struct{}

type MockMem struct{}

// Implement necessary methods for MockMem to satisfy the interface
// For example, if the interface has a method called GetMemoryUsage, implement it:
func (m *MockMem) GetMemoryUsage() float64 {
	return 0.0 // Return a mock value
}

func (m *MockMem) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return &mem.VirtualMemoryStat{}, nil
}

type MockMonitor struct{}

func (m *MockMonitor) GetCPUUsage() (float64, error) {
	return 0.0, nil
}

func (m *MockMonitor) Close() error {
	return nil
}

func (m *MockMonitor) GetMemoryUsage() (uint64, error) {
	return 0, nil
}

func (m *MockMonitor) GetTotalMemory() (uint64, error) {
	return 0, nil
}

func (m *MockMonitor) Initialize() error {
	return nil
}

func (m *MockMonitor) GetDiskUsage() (core.DiskStats, error) {
	return core.DiskStats{}, nil
}

func (m *MockMonitor) GetNetworkStats() (core.NetworkStats, error) {
	return core.NetworkStats{}, nil
}

type MockHTTPExporter struct{}

func (m *MockHTTPExporter) Export(data map[string]interface{}) error {
	return nil
}

func TestNewAgent(t *testing.T) {
	cfg := &config.Config{
		LogFilePath: "./agent.log",
		Interval:    60,
	}
	mon := &MockMonitor{}
	mc := metrics.NewMetricsCollector(mon, cfg)
	exporters := []exporter.Exporter{
		exporter.NewFileExporter(cfg.LogFilePath),
		exporter.NewHTTPExporter(cfg.HTTP.Endpoint),
	}
	agent := NewAgent(cfg, mc, exporters...)

	assert.NotNil(t, agent)
	assert.Equal(t, cfg, agent.config)
	assert.Equal(t, time.Duration(cfg.Interval)*time.Second, agent.interval)
}

func TestAgentStart(t *testing.T) {
	cfg := &config.Config{
		LogFilePath: "./agent.log",
		Interval:    1, // Set a short interval for testing
	}
	mon := &MockMonitor{}
	mc := metrics.NewMetricsCollector(mon, cfg)
	exporters := []exporter.Exporter{
		exporter.NewFileExporter(cfg.LogFilePath),
		exporter.NewHTTPExporter(cfg.HTTP.Endpoint),
	}
	agent := NewAgent(cfg, mc, exporters...)

	// Run Start in a separate goroutine
	go agent.Start()

	// Allow some time for the ticker to tick
	time.Sleep(3 * time.Second)

	// Here you would typically check if the metrics were exported correctly.
	// This might involve checking the output file or using a mock exporter.
	// For simplicity, we will just assert that the function runs without panic.
	assert.NotNil(t, agent)
}
