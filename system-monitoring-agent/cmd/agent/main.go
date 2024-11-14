// cmd/agent/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/travism26/system-monitoring-agent/internal/config"

	"github.com/travism26/system-monitoring-agent/internal/agent"
	"github.com/travism26/system-monitoring-agent/internal/exporter"
	"github.com/travism26/system-monitoring-agent/internal/logger"
	"github.com/travism26/system-monitoring-agent/internal/metrics"
	"github.com/travism26/system-monitoring-agent/internal/monitor"
)

func main() {
	fmt.Println("Starting Windows Monitoring Agent...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize components
	mon, err := monitor.NewSystemMonitor()
	if err != nil {
		log.Fatalf("Error creating system monitor: %v", err)
	}
	mc := metrics.NewMetricsCollector(mon)
	exp := exporter.NewExporter(cfg.LogFilePath)

	ag := agent.NewAgent(cfg, mc, exp)

	// Handle termination signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Handle termination signal
	go func() {
		<-ctx.Done()
		fmt.Println("Received termination signal, stopping agent...")
		stop()
		os.Exit(0)
	}()

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