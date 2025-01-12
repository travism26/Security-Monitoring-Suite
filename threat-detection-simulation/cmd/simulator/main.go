package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
	"github.com/travism26/threat-detection-simulation/internal/config"
	"github.com/travism26/threat-detection-simulation/internal/logger"
	"github.com/travism26/threat-detection-simulation/internal/scenarios"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "internal/config/config.yaml", "Path to configuration file")
	scenarioType := flag.String("scenario", "high-cpu", "Type of scenario to run (high-cpu, malicious-process)")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version if requested
	if *showVersion {
		fmt.Println("Threat Detection Simulator v0.1.0")
		os.Exit(0)
	}

	// Initialize logger
	l := logger.New("[Simulator]")

	// Load configuration
	l.Printf("Loading configuration from %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	// Start simulation in a goroutine
	go func() {
		l.Printf("Starting %s scenario simulation", *scenarioType)
		for {
			var metrics types.MetricPayload

			// Generate scenario based on type
			switch *scenarioType {
			case "high-cpu":
				metrics = scenarios.HighCPUScenario()
			case "malicious-process":
				metrics = scenarios.MaliciousProcessScenario()
			default:
				l.Printf("Unknown scenario type: %s", *scenarioType)
				continue
			}

			// Validate metrics
			if err := scenarios.ValidateMetrics(metrics); err != nil {
				l.Printf("Invalid metrics: %v", err)
				continue
			}

			// Send metrics to endpoint
			if err := scenarios.SendMetrics(cfg.Endpoint, metrics); err != nil {
				l.Printf("Failed to send metrics: %v", err)
			} else {
				l.LogMetricsSent(*scenarioType, cfg.Endpoint)
			}

			select {
			case <-done:
				l.Printf("Stopping simulation")
				return
			case <-time.After(time.Duration(cfg.Interval) * time.Second):
				// Continue to next iteration
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	l.Printf("Received shutdown signal")
	done <- true
	l.Printf("Simulation stopped")
}
