// internal/core/monitor.go
package core

// Interface for the monitor component to be cross-compatible with different implementations
// Windows and MacOS ATM
type Monitor interface {
	Initialize() error
	Close() error
	GetCPUUsage() (float64, error)
	GetMemoryUsage() (uint64, error)
	GetTotalMemory() (uint64, error)
	GetDiskUsage() (DiskStats, error)
	GetNetworkStats() (NetworkStats, error)
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
