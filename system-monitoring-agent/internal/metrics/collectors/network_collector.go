// internal/metrics/collectors/network_collector.go
package collectors

import "github.com/travism26/system-monitoring-agent/internal/core"

type NetworkCollector struct {
	monitor core.NetworkMonitor
}

func NewNetworkCollector(monitor core.NetworkMonitor) *NetworkCollector {
	return &NetworkCollector{monitor: monitor}
}

func (c *NetworkCollector) Name() string {
	return "network"
}

func (c *NetworkCollector) Collect() (map[string]interface{}, error) {
	networkStats, err := c.monitor.GetNetworkStats()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"network": networkStats,
	}, nil
}
