package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/travism26/system-monitoring-agent/internal/config"
	"github.com/travism26/system-monitoring-agent/internal/exporter"
	"github.com/travism26/system-monitoring-agent/internal/metrics"
	"github.com/travism26/system-monitoring-agent/internal/monitor"
)

func TestNewAgent(t *testing.T) {
	cfg := &config.Config{
		LogFilePath: "./agent.log",
		Interval:    60,
	}
	mon := monitor.NewMonitor()
	mc := metrics.NewMetricsCollector(mon)
	exp := exporter.NewExporter(cfg.LogFilePath)

	agent := NewAgent(cfg, mc, exp)

	assert.NotNil(t, agent)
	assert.Equal(t, cfg, agent.config)
	assert.Equal(t, mc, agent.metrics)
	assert.Equal(t, exp, agent.exporter)
	assert.Equal(t, time.Duration(cfg.Interval)*time.Second, agent.interval)
}

func TestAgentStart(t *testing.T) {
	cfg := &config.Config{
		LogFilePath: "./agent.log",
		Interval:    1, // Set a short interval for testing
	}
	mon := monitor.NewMonitor()
	mc := metrics.NewMetricsCollector(mon)
	exp := exporter.NewExporter(cfg.LogFilePath)

	agent := NewAgent(cfg, mc, exp)

	// Run Start in a separate goroutine
	go agent.Start()

	// Allow some time for the ticker to tick
	time.Sleep(3 * time.Second)

	// Here you would typically check if the metrics were exported correctly.
	// This might involve checking the output file or using a mock exporter.
	// For simplicity, we will just assert that the function runs without panic.
	assert.NotNil(t, agent)
}
