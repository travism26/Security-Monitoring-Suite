package domain

import (
	"time"
)

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

	// Organization that owns this log entry
	OrganizationID string `json:"organization_id"`

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

	// Enrichment fields
	Environment    string    `json:"environment"`     // e.g., production, staging, development
	Application    string    `json:"application"`     // Application name
	Component      string    `json:"component"`       // Component within the application
	CorrelationID  string    `json:"correlation_id"`  // For tracking related logs
	Tags           []string  `json:"tags"`            // Custom tags for categorization
	EnrichedAt     time.Time `json:"enriched_at"`     // When the log was enriched
	ProcessedCount int       `json:"processed_count"` // Number of processing attempts
}

// EnrichLog adds additional context and metadata to the log entry
func (l *Log) EnrichLog(env, app, component string) {
	l.Environment = env
	l.Application = app
	l.Component = component
	l.EnrichedAt = time.Now()

	// Add standard tags based on existing data
	l.Tags = append(l.Tags,
		"host:"+l.Host,
		"level:"+l.Level,
		"env:"+env,
	)

	// Add process-related metadata
	if l.Metadata == nil {
		l.Metadata = make(map[string]interface{})
	}
	l.Metadata["process_metrics"] = map[string]interface{}{
		"process_count":     l.ProcessCount,
		"total_cpu_percent": l.TotalCPUPercent,
		"total_memory":      l.TotalMemoryUsage,
	}
}

// LogRepository defines the interface for storing and retrieving logs.
// This interface allows us to:
// 1. Abstract the database implementation details
// 2. Easily switch between different database types
// 3. Write mock implementations for testing
type LogRepository interface {
	// Store saves a single log entry to the database
	Store(log *Log) error

	// StoreBatch saves multiple log entries to the database in a single transaction
	// This is more efficient than storing logs one by one when processing multiple logs
	StoreBatch(logs []*Log) error

	// FindByID retrieves a specific log entry by its ID and organization
	FindByID(orgID, id string) (*Log, error)

	// List retrieves multiple logs for an organization with pagination support
	// limit: maximum number of logs to return
	// offset: number of logs to skip
	List(orgID string, limit, offset int) ([]*Log, error)

	// ListByTimeRange retrieves logs within a specific time range with pagination
	// This is useful for querying logs within a specific window
	ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*Log, error)

	// CountByTimeRange returns the total number of logs within a time range for an organization
	CountByTimeRange(orgID string, start, end time.Time) (int64, error)

	// ListByHost retrieves logs for a specific host within an organization
	ListByHost(orgID string, host string, limit, offset int) ([]*Log, error)

	// ListByLevel retrieves logs of a specific level within an organization
	ListByLevel(orgID string, level string, limit, offset int) ([]*Log, error)
}
