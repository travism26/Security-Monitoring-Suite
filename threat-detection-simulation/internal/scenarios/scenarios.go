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
			ThreatIndicators: []types.ThreatIndicator{
				{
					Type:        "scenario",
					Description: "High CPU usage",
					Severity:    "HIGH",
					Score:       100,
					Timestamp:   time.Now(),
					Metadata:    types.Metadata{Tags: []string{"high_cpu"}},
				},
				{
					Type:        "scenario",
					Description: "Malicious process",
					Severity:    "HIGH",
					Score:       100,
					Timestamp:   time.Now(),
					Metadata:    types.Metadata{Tags: []string{"malicious_process"}},
				},
				{
					Type:        "scenario",
					Description: "High memory usage",
					Severity:    "HIGH",
					Score:       100,
					Timestamp:   time.Now(),
					Metadata:    types.Metadata{Tags: []string{"high_memory"}},
				},
				{
					Type:        "scenario",
					Description: "Network activity",
					Severity:    "HIGH",
					Score:       100,
					Timestamp:   time.Now(),
					Metadata:    types.Metadata{Tags: []string{"network_activity"}},
				},
			},
		},
	}
}

// MaliciousProcessScenario simulates a known malicious process
func MaliciousProcessScenario() types.MetricPayload {
	return types.MetricPayload{
		Timestamp: time.Now().Format(time.RFC3339),
		Data: types.MetricData{
			HostInfo: types.HostInfo{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			Metrics: map[string]interface{}{
				"processes": []map[string]interface{}{
					{
						"name":           "malicious.exe",
						"cpu_percent":    80.0,
						"memory_percent": 20.0,
					},
				},
			},
		},
	}
}
