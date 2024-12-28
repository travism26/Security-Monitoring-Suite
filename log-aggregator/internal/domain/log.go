package domain

import "time"

// Package domain contains the core business logic and entities of the application.
// It defines the central types and interfaces that other packages will implement.

// Log represents a single log entry in the system.
// This is the core domain model that will flow through all layers of the application:
// 1. Received from Kafka
// 2. Processed by services
// 3. Stored in database
// 4. Returned via API
type Log struct {
	// Unique identifier for the log entry
	ID string `json:"id"`

	// When the log event occurred
	Timestamp time.Time `json:"timestamp"`

	// Source of the log (e.g., server name, container ID)
	Host string `json:"host"`

	// The actual log message
	Message string `json:"message"`

	// Log level (e.g., INFO, ERROR, WARN)
	Level string `json:"level"`

	// Additional structured data associated with the log
	Metadata    map[string]interface{} `json:"-"`
	MetadataStr string                 `json:"metadata"`

	// Total number of processes in the log
	ProcessCount int `json:"process_count"`

	// Total CPU usage percentage for all processes
	TotalCPUPercent float64 `json:"total_cpu_percent"`

	// Total memory usage for all processes
	TotalMemoryUsage int64 `json:"total_memory_usage"`
}

// LogRepository defines the interface for storing and retrieving logs.
// This interface allows us to:
// 1. Abstract the database implementation details
// 2. Easily switch between different database types
// 3. Write mock implementations for testing
type LogRepository interface {
	// Store saves a single log entry to the database
	Store(log *Log) error

	// FindByID retrieves a specific log entry by its ID
	FindByID(id string) (*Log, error)

	// List retrieves multiple logs with pagination support
	// limit: maximum number of logs to return
	// offset: number of logs to skip
	List(limit, offset int) ([]*Log, error)
}
