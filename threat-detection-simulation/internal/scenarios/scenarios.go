package scenarios

import (
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
)

// HighCPUScenario simulates a high CPU usage attack
func HighCPUScenario() types.MetricPayload {
	return types.MetricPayload{
		Timestamp: time.Now(),
		CPUUsage:  95.0,
	}
}

// MaliciousProcessScenario simulates a known malicious process
func MaliciousProcessScenario() types.MetricPayload {
	return types.MetricPayload{
		// ... scenario configuration
	}
}
