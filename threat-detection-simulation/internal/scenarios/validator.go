package scenarios

import (
	"fmt"

	"github.com/travism26/shared-monitoring-libs/types"
)

func ValidateMetrics(metrics types.MetricPayload) error {
	if metrics.Data.Metrics["cpu_usage"].(types.CPUUsage).Usage < 0 || metrics.Data.Metrics["cpu_usage"].(types.CPUUsage).Usage > 100 {
		return fmt.Errorf("invalid CPU usage: %f", metrics.Data.Metrics["cpu_usage"].(types.CPUUsage).Usage)
	}

	// Add more validation rules

	return nil
}
