// internal/monitor/monitor.go
package monitor

import (
	"log"

	"github.com/shirou/gopsutil/mem"
)

// CPU interface for mocking
type CPU interface {
	Percent(interval float64, percpu bool) ([]float64, error)
}

// Mem interface for mocking
type Mem interface {
	VirtualMemory() (*mem.VirtualMemoryStat, error)
}

// Monitor struct
type Monitor struct {
	cpu CPU
	mem Mem
}

// NewMonitor creates a new Monitor
func NewMonitor(cpu CPU, mem Mem) *Monitor {
	return &Monitor{cpu: cpu, mem: mem}
}

// GetCPUUsage gets the CPU usage
func (m *Monitor) GetCPUUsage() (float64, error) {
	percentages, err := m.cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return percentages[0], nil
}

// GetMemoryUsage gets the memory usage
func (m *Monitor) GetMemoryUsage() (uint64, error) {
	vmStat, err := m.mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Used, nil
}

// LogSystemMetrics logs the system metrics
func (m *Monitor) LogSystemMetrics() {
	cpuUsage, err := m.GetCPUUsage()
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
	}
	log.Printf("CPU Usage: %v%%", cpuUsage)

	memUsage, err := m.GetMemoryUsage()
	if err != nil {
		log.Printf("Error getting Memory usage: %v", err)
	}
	log.Printf("Memory Usage: %v bytes", memUsage)
}
