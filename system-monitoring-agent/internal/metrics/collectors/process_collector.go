// internal/metrics/collectors/process_collector.go
package collectors

import (
	"github.com/travism26/shared-monitoring-libs/types"
	"github.com/travism26/system-monitoring-agent/internal/core"
)

type ProcessCollector struct {
	monitor core.ProcessMonitor
}

func NewProcessCollector(monitor core.ProcessMonitor) *ProcessCollector {
	return &ProcessCollector{monitor: monitor}
}

func (c *ProcessCollector) Name() string {
	return "processes"
}

func (c *ProcessCollector) Collect() (map[string]interface{}, error) {
	processes, err := c.monitor.GetProcesses()
	if err != nil {
		return nil, err
	}

	// Create summary metrics
	var totalCPU float64
	var totalMemory uint64
	processCount := len(processes)

	// Convert directly to ProcessInfo
	processList := make([]types.ProcessInfo, 0, processCount)

	for _, proc := range processes {
		totalCPU += proc.CPUPercent
		totalMemory += proc.MemoryUsage

		// Create ProcessInfo directly
		processMetric := types.ProcessInfo{
			Name:        proc.Name,
			PID:         proc.PID,
			CPUPercent:  proc.CPUPercent,
			MemoryUsage: proc.MemoryUsage,
			Status:      proc.Status,
		}
		processList = append(processList, processMetric)
	}

	result := map[string]interface{}{
		"processes": map[string]interface{}{
			"total_count":        processCount,
			"total_cpu_percent":  totalCPU,
			"total_memory_usage": totalMemory,
			"process_list":       processList,
		},
	}

	return result, nil
}
