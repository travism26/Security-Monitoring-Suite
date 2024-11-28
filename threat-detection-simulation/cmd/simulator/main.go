package main

import (
	"flag"
	"log"
	"time"

	"github.com/travism26/threat-detection-simulation/internal/scenarios"
)

func main() {
	// Command line flags
	endpoint := flag.String("endpoint", "http://localhost:3000/api/v1/system-metrics/ingest", "API endpoint URL")
	interval := flag.Int("interval", 60, "Interval between metrics in seconds")
	scenario := flag.String("scenario", "high-cpu", "Scenario to run (high-cpu, malicious-process)")
	flag.Parse()

	log.Printf("Starting threat simulation with scenario: %s\n", *scenario)
	log.Printf("Sending metrics to: %s\n", *endpoint)
	log.Printf("Interval: %d seconds\n", *interval)

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var metrics interface{}
		switch *scenario {
		case "high-cpu":
			metrics = scenarios.HighCPUScenario()
		case "malicious-process":
			metrics = scenarios.MaliciousProcessScenario()
		default:
			log.Fatalf("Unknown scenario: %s", *scenario)
		}

		if err := scenarios.SendMetrics(*endpoint, metrics); err != nil {
			log.Printf("Error sending metrics: %v", err)
		}
	}
}
