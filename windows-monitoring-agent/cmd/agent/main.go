// cmd/agent/main.go
package main

import (
	"fmt"
	"log"

	"github.com/travism26/windows-monitoring-agent/internal/config"
	"github.com/travism26/windows-monitoring-agent/internal/logger"
)

func main() {
	fmt.Println("Starting Windows Monitoring Agent...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

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
	// Monitoring logic goes here
	log.Println("Agent is running...")
	// Add system monitoring code
}
