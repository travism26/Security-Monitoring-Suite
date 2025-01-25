// cmd/agent/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/travism26/system-monitoring-agent/internal/config"

	"github.com/travism26/system-monitoring-agent/internal/agent"
	"github.com/travism26/system-monitoring-agent/internal/exporter"
	"github.com/travism26/system-monitoring-agent/internal/metrics"
	"github.com/travism26/system-monitoring-agent/internal/monitor"
)

func main() {
	fmt.Println("Starting System Monitoring Agent...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	storage, err := exporter.NewMetricStorage(cfg.StorageDir)
	if err != nil {
		log.Fatalf("Error creating metric storage: %v", err)
	}

	// Initialize components
	mon, err := monitor.NewSystemMonitor()
	if err != nil {
		log.Fatalf("Error creating system monitor: %v", err)
	}
	mc := metrics.NewMetricsCollector(mon, cfg)

	// Initialize exporters
	exporters := []exporter.MetricsExporter{
		exporter.NewFileExporter(cfg.LogFilePath),
	}

	// Initialize HTTP exporter with tenant context
	httpExporter, err := exporter.NewHTTPExporter(cfg, storage)
	if err != nil {
		log.Fatalf("Error creating HTTP exporter: %v", err)
	}
	exporters = append(exporters, httpExporter)

	// Initialize agent
	ag := agent.NewAgent(cfg, mc, exporters...)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start agent in goroutine
	done := make(chan struct{})
	go func() {
		ag.Start(done)
	}()

	// Wait for termination signal
	<-sigChan
	fmt.Println("Received termination signal, stopping agent...")

	// Signal agent to stop
	close(done)

	// Give the agent time to clean up
	time.Sleep(time.Second)
	fmt.Println("Agent shutdown complete")
}
