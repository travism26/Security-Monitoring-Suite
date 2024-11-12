package monitor

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// SystemCPU implements CPU interface
type SystemCPU struct{}

func NewCPUMonitor() CPU {
	return &SystemCPU{}
}

func (c *SystemCPU) Percent(interval float64, percpu bool) ([]float64, error) {
	return cpu.Percent(time.Duration(interval)*time.Second, percpu)
}

// SystemMem implements Mem interface
type SystemMem struct{}

func NewMemMonitor() Mem {
	return &SystemMem{}
}

func (m *SystemMem) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}
