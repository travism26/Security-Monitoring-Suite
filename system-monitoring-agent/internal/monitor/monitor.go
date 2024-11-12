// internal/monitor/monitor.go
package monitor

import (
	"log"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (m *Monitor) GetCPUUsage() (float64, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return percent[0], nil
}

func (m *Monitor) GetMemoryUsage() (uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return v.Used, nil
}

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
