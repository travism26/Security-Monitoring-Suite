// internal/core/monitor.go
package core

// Interface for the monitor component to be cross-compatible with different implementations
// Windows and MacOS ATM
type Monitor interface {
	GetCPUUsage() (float64, error)
	GetMemoryUsage() (uint64, error)
	GetTotalMemory() (uint64, error)
	Initialize() error
	Close() error
}
