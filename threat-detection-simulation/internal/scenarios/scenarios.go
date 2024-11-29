package scenarios

import (
	"runtime"
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
)

// HighCPUScenario simulates a high CPU usage attack
func HighCPUScenario() types.MetricPayload {
	return types.MetricPayload{
		Timestamp: time.Now().Format(time.RFC3339),
		Data: types.MetricData{
			HostInfo: types.HostInfo{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			Metrics: map[string]interface{}{
				"cpu_usage": types.CPUUsage{Usage: 95.0, Total: 100},
			},
		},
	}
}

// MaliciousProcessScenario simulates a known malicious process
func MaliciousProcessScenario() types.MetricPayload {
	return types.MetricPayload{
		// ... scenario configuration
	}
}
