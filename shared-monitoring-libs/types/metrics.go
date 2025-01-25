package types

import (
	"time"
)

// MetricPayload represents the common structure for system metrics
// TenantContext holds tenant-specific information
type TenantContext struct {
	ID       string            `json:"id"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type MetricPayload struct {
	Timestamp string        `json:"timestamp"`
	Tenant    TenantContext `json:"tenant"`
	Data      MetricData    `json:"data"`
}

type MetricData struct {
	HostInfo         HostInfo               `json:"host_info"`
	Metrics          map[string]interface{} `json:"metrics"`
	ThreatIndicators []ThreatIndicator      `json:"threat_indicators"`
	Metadata         MetadataInfo           `json:"metadata"`
	Processes        SystemProcessStats     `json:"processes"`
}

type HostInfo struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Hostname  string `json:"hostname"`
	CPUCores  int    `json:"cpu_cores"`
	GoVersion string `json:"go_version"`
}

// Keep these for backwards compatibility or specific parsing needs
type MemoryMetrics struct {
	Used    int64   `json:"used"`
	Total   int64   `json:"total"`
	Percent float64 `json:"percent"`
}

type ProcessInfo struct {
	Name        string  `json:"name"`
	PID         int     `json:"pid"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsage uint64  `json:"memory_usage"`
	Status      string  `json:"status"`
}

type SystemProcessStats struct {
	TotalCount       int           `json:"total_count"`
	TotalCPUPercent  float64       `json:"total_cpu_percent"`
	TotalMemoryUsage uint64        `json:"total_memory_usage"`
	ProcessList      []ProcessInfo `json:"process_list"`
}

type CPUUsage struct {
	Usage float64 `json:"usage"`
	Total float64 `json:"total"`
}

type MetadataInfo struct {
	CollectionDuration string            `json:"collection_duration"`
	CollectorCount     int               `json:"collector_count"`
	Errors             []string          `json:"errors,omitempty"`
	TenantMetadata     map[string]string `json:"tenant_metadata,omitempty"`
}

type ThreatIndicator struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Score       float64                `json:"score"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    Metadata               `json:"metadata"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

type Metadata struct {
	Tags []string `json:"tags"`
}
