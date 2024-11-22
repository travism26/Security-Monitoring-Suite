package os_darwin

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/travism26/system-monitoring-agent/internal/core"
)

// DarwinMonitor implements the core.Monitor interface
type DarwinMonitor struct {
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

// NewDarwinMonitor creates a new DarwinMonitor instance
func NewDarwinMonitor() core.SystemMonitor {
	return &DarwinMonitor{
		cpu:  &DarwinCPU{},
		mem:  &DarwinMem{},
		proc: &DarwinProcess{},
	}
}

// GetCPUUsage implements core.Monitor interface
func (m *DarwinMonitor) GetCPUUsage() (float64, error) {
	percentages, err := m.cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return percentages[0], nil
}

// GetMemoryUsage implements core.Monitor interface
func (m *DarwinMonitor) GetMemoryUsage() (uint64, error) {
	vmStat, err := m.mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Used, nil
}

// GetTotalMemory implements core.Monitor interface
func (m *DarwinMonitor) GetTotalMemory() (uint64, error) {
	vmStat, err := m.mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total, nil
}

// Initialize implements core.Monitor interface
func (m *DarwinMonitor) Initialize() error {
	// Add any Darwin-specific initialization if needed
	return nil
}

// Close implements core.Monitor interface
func (m *DarwinMonitor) Close() error {
	// Add any cleanup if needed
	return nil
}

// DarwinCPU implements CPU interface
type DarwinCPU struct{}

func (c *DarwinCPU) Percent(interval float64, percpu bool) ([]float64, error) {
	return cpu.Percent(time.Duration(interval)*time.Second, percpu)
}

// DarwinMem implements Mem interface
type DarwinMem struct{}

func (m *DarwinMem) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

// DarwinMonitor implements core.Monitor interface
func (m *DarwinMonitor) GetDiskUsage() (core.DiskStats, error) {
	_, err := disk.Partitions(false)
	if err != nil {
		return core.DiskStats{}, err
	}

	// Get root partition usage
	usage, err := disk.Usage("/")
	if err != nil {
		return core.DiskStats{}, err
	}

	fmt.Println(usage)

	return core.DiskStats{
		Total: usage.Total,
		Used:  usage.Used,
		Free:  usage.Free,
	}, nil
}

func (m *DarwinMonitor) GetNetworkStats() (core.NetworkStats, error) {
	stats, err := net.IOCounters(false)
	if err != nil {
		return core.NetworkStats{}, err
	}

	if len(stats) == 0 {
		return core.NetworkStats{}, fmt.Errorf("no network stats available")
	}

	// Aggregate all interfaces
	var totalSent, totalReceived uint64
	for _, stat := range stats {
		totalSent += stat.BytesSent
		totalReceived += stat.BytesRecv
	}

	return core.NetworkStats{
		BytesSent:     totalSent,
		BytesReceived: totalReceived,
	}, nil
}

// DarwinProcess implements Process interface
type DarwinProcess struct{}

func (p *DarwinProcess) Processes() ([]*process.Process, error) {
	return process.Processes()
}

// Implement GetProcesses method
func (m *DarwinMonitor) GetProcesses() ([]core.ProcessInfo, error) {
	processes, err := m.proc.Processes()
	if err != nil {
		return nil, err
	}

	var processInfos []core.ProcessInfo
	for _, p := range processes {
		pid := p.Pid
		name, err := p.Name()
		if err != nil {
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		status, err := p.Status()
		if err != nil {
			status = "unknown"
		}

		processInfos = append(processInfos, core.ProcessInfo{
			PID:         int(pid),
			Name:        name,
			CPUPercent:  cpuPercent,
			MemoryUsage: memInfo.RSS,
			Status:      status,
		})
	}

	return processInfos, nil
}
