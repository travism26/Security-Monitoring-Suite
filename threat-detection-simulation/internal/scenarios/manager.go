package scenarios

import (
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
	"github.com/travism26/threat-detection-simulation/internal/config"
)

type ScenarioManager struct {
	config *config.Config
}

func NewScenarioManager(cfg *config.Config) *ScenarioManager {
	return &ScenarioManager{
		config: cfg,
	}
}

func (sm *ScenarioManager) NextScenario() types.MetricPayload {
	// Rotate between different scenarios
	now := time.Now().Unix()
	if now%2 == 0 {
		return HighCPUScenario()
	}
	return MaliciousProcessScenario()
}
