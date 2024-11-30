package main

import (
	"log"
	"time"

	"github.com/travism26/threat-detection-simulation/internal/config"
	"github.com/travism26/threat-detection-simulation/internal/scenarios"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("internal/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run simulation loop
	for {
		// Generate high CPU scenario
		metrics := scenarios.HighCPUScenario()

		// Validate metrics
		if err := scenarios.ValidateMetrics(metrics); err != nil {
			log.Printf("Invalid metrics: %v", err)
			continue
		}

		// Send metrics to endpoint
		if err := scenarios.SendMetrics(cfg.Endpoint, metrics); err != nil {
			log.Printf("Failed to send metrics: %v", err)
		} else {
			log.Printf("Successfully sent metrics to %s", cfg.Endpoint)
		}

		// Wait for configured interval
		time.Sleep(time.Duration(cfg.Interval) * time.Second)
	}
}
