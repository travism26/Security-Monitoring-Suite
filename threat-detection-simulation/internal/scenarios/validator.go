package scenarios

import (
	"fmt"

	"github.com/travism26/shared-monitoring-libs/types"
)

func ValidateMetrics(metrics types.MetricPayload) error {
	if metrics.CPUUsage < 0 || metrics.CPUUsage > 100 {
		return fmt.Errorf("invalid CPU usage: %f", metrics.CPUUsage)
	}

	// Add more validation rules

	return nil
}
