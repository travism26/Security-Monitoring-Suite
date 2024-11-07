// cmd/agent/main.go
package main

import (
	"fmt"
	"log"

	"github.com/travism26/windows-monitoring-agent/internal/agent"
	"github.com/travism26/windows-monitoring-agent/internal/config"
	"github.com/travism26/windows-monitoring-agent/internal/exporter"
	"github.com/travism26/windows-monitoring-agent/internal/logger"
	"github.com/travism26/windows-monitoring-agent/internal/metrics"
	"github.com/travism26/windows-monitoring-agent/internal/monitor"
)

func main() {
	fmt.Println("Starting Windows Monitoring Agent...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize components
	mon := monitor.NewMonitor()
	mc := metrics.NewMetricsCollector(mon)
	exp := exporter.NewExporter(cfg.LogFilePath)

	ag := agent.NewAgent(cfg, mc, exp)
	ag.Start()

	// Initialize logger
	logFile, err := logger.Init(cfg.LogFilePath)
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logFile.Close() // Ensure log file is closed when the application exits

	log.Println("Logging started...")

	// Placeholder for monitoring logic
	runMonitoringAgent(cfg)
}

func runMonitoringAgent(cfg *config.Config) {
	log.Println("Agent is running and logging to", cfg.LogFilePath)
	// Monitoring logic goes here
	log.Println("Agent is running...")
	// Add system monitoring code
}
