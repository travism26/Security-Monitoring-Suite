package types

import "time"

// MetricPayload represents the common structure for system metrics
type MetricPayload struct {
	Timestamp   time.Time        `json:"timestamp"`
	CPUUsage    float64          `json:"cpu_usage"`
	MemoryUsage MemoryMetrics    `json:"memory_usage"`
	Processes   []ProcessMetrics `json:"processes"`
	// Add other common fields here
}

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
