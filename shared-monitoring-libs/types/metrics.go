package types

import (
	"time"
)

// MetricPayload represents the common structure for system metrics
type MetricPayload struct {
	Timestamp string     `json:"timestamp"`
	Data      MetricData `json:"data"`
}

type MetricData struct {
	HostInfo         HostInfo               `json:"host_info"`
	Metrics          map[string]interface{} `json:"metrics"`
	ThreatIndicators []ThreatIndicator      `json:"threat_indicators"`
	Metadata         MetadataInfo           `json:"metadata"`
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

type ProcessMetrics struct {
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
}

type CPUUsage struct {
	Usage float64 `json:"usage"`
	Total float64 `json:"total"`
}

type MetadataInfo struct {
	CollectionDuration string   `json:"collection_duration"`
	CollectorCount     int      `json:"collector_count"`
	Errors             []string `json:"errors,omitempty"`
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
