package os_windows

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/travism26/system-monitoring-agent/internal/core"
)

// WindowsMonitor implements the core.Monitor interface
type WindowsMonitor struct {
	cpu  CPU
	mem  Mem
	proc Process
}

// CPU interface (same as in monitor.go)
type CPU interface {
	Percent(interval float64, percpu bool) ([]float64, error)
}

// Mem interface (same as in monitor.go)
type Mem interface {
	VirtualMemory() (*mem.VirtualMemoryStat, error)
}

// Process interface
type Process interface {
	Processes() ([]*process.Process, error)
}

// NewWindowsMonitor creates a new WindowsMonitor instance
func NewWindowsMonitor() core.SystemMonitor {
	return &WindowsMonitor{
		cpu:  &WindowsCPU{},
		mem:  &WindowsMem{},
		proc: &WindowsProcess{},
	}
}

func (m *WindowsMonitor) GetDiskUsage() (core.DiskStats, error) {
	return core.DiskStats{}, nil
}

func (m *WindowsMonitor) GetNetworkStats() (core.NetworkStats, error) {
	return core.NetworkStats{}, nil
}

// GetCPUUsage implements core.Monitor interface
func (m *WindowsMonitor) GetCPUUsage() (float64, error) {
	percentages, err := m.cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return percentages[0], nil
}

// GetMemoryUsage implements core.Monitor interface
func (m *WindowsMonitor) GetMemoryUsage() (uint64, error) {
	vmStat, err := m.mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Used, nil
}

// GetTotalMemory implements core.Monitor interface
func (m *WindowsMonitor) GetTotalMemory() (uint64, error) {
	vmStat, err := m.mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total, nil
}

// Initialize implements core.Monitor interface
func (m *WindowsMonitor) Initialize() error {
	// Add any Windows-specific initialization if needed
	return nil
}

// Close implements core.Monitor interface
func (m *WindowsMonitor) Close() error {
	// Add any cleanup if needed
	return nil
}

// WindowsCPU implements CPU interface
type WindowsCPU struct{}

func (c *WindowsCPU) Percent(interval float64, percpu bool) ([]float64, error) {
	return cpu.Percent(time.Duration(interval)*time.Second, percpu)
}

// WindowsMem implements Mem interface
type WindowsMem struct{}

func (m *WindowsMem) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

// WindowsProcess implements Process interface
type WindowsProcess struct{}

func (p *WindowsProcess) Processes() ([]*process.Process, error) {
	return process.Processes()
}

// Add GetProcesses method
func (m *WindowsMonitor) GetProcesses() ([]core.ProcessInfo, error) {
	processes, err := m.proc.Processes()
	if err != nil {
		return nil, err
	}

	result := make([]core.ProcessInfo, 0, len(processes))
	for _, p := range processes {
		pid := int(p.Pid)
		name, err := p.Name()
		if err != nil {
			continue
		}
		result = append(result, core.ProcessInfo{
			PID:  pid,
			Name: name,
		})
	}
	return result, nil
}
