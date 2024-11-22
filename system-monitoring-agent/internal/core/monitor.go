// internal/core/monitor.go
package core

// SystemMonitor composes all monitoring capabilities
type SystemMonitor interface {
	CPUMonitor
	MemoryMonitor
	DiskMonitor
	NetworkMonitor
	ProcessMonitor
}

// CPUMonitor handles CPU-specific metrics
type CPUMonitor interface {
	GetCPUUsage() (float64, error)
}

// MemoryMonitor handles memory-specific metrics
type MemoryMonitor interface {
	GetMemoryUsage() (uint64, error)
	GetTotalMemory() (uint64, error)
}

// DiskMonitor handles disk-specific metrics
type DiskMonitor interface {
	GetDiskUsage() (DiskStats, error)
}

// NetworkMonitor handles network-specific metrics
type NetworkMonitor interface {
	GetNetworkStats() (NetworkStats, error)
}

// ProcessMonitor handles process-specific metrics
type ProcessMonitor interface {
	GetProcesses() ([]ProcessInfo, error)
}

type DiskStats struct {
	Total uint64
	Used  uint64
	Free  uint64
}

type NetworkStats struct {
	BytesSent     uint64
	BytesReceived uint64
}

type ProcessInfo struct {
	PID         int
	Name        string
	CPUPercent  float64
	MemoryUsage uint64
	Status      string
}
