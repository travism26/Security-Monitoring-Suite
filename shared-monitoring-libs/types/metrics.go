package types

import (
	"time"
)

// MetricPayload represents the optimized structure for system metrics
type MetricPayload struct {
	Timestamp string `json:"timestamp"`
	TenantID  string `json:"tenant_id"`

	// Tenant metadata
	TenantMetadata map[string]string `json:"tenant_metadata,omitempty"`

	// Host information
	Host struct {
		OS        string `json:"os"`
		Arch      string `json:"arch"`
		Hostname  string `json:"hostname"`
		CPUCores  int    `json:"cpu_cores"`
		GoVersion string `json:"go_version"`
	} `json:"host"`

	// Core metrics data
	Metrics map[string]interface{} `json:"metrics"`

	// Process information
	Processes struct {
		TotalCount       int           `json:"total_count"`
		TotalCPUPercent  float64       `json:"total_cpu_percent"`
		TotalMemoryUsage uint64        `json:"total_memory_usage"`
		List             []ProcessInfo `json:"list"`
	} `json:"processes"`

	// Threat detection
	ThreatIndicators []ThreatIndicator `json:"threat_indicators,omitempty"`

	// Collection metadata
	Metadata struct {
		CollectionDuration string   `json:"collection_duration"`
		CollectorCount     int      `json:"collector_count"`
		Errors             []string `json:"errors,omitempty"`
	} `json:"metadata"`
}

// ProcessInfo represents information about a single process
type ProcessInfo struct {
	Name        string  `json:"name"`
	PID         int     `json:"pid"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsage uint64  `json:"memory_usage"`
	Status      string  `json:"status"`
}

// ThreatIndicator represents a security threat or anomaly
type ThreatIndicator struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Score       float64                `json:"score"`
	Timestamp   time.Time              `json:"timestamp"`
	Tags        []string               `json:"tags"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// Utility types for specific metric parsing
type MemoryMetrics struct {
	Used    int64   `json:"used"`
	Total   int64   `json:"total"`
	Percent float64 `json:"percent"`
}

type CPUUsage struct {
	Usage float64 `json:"usage"`
	Total float64 `json:"total"`
}
