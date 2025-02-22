package threat

import (
	"math"
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
)

type Analyzer struct {
	thresholds map[string]float64
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		thresholds: map[string]float64{
			"cpu_usage":    8.0,
			"memory_usage": 85.0,
			"disk_usage":   90.0,
		},
	}
}

func (a *Analyzer) AnalyzeMetrics(metrics map[string]interface{}) []types.ThreatIndicator {
	var indicators []types.ThreatIndicator
	now := time.Now()

	// Analyze CPU usage
	if cpuUsage, ok := metrics["cpu_usage"].(float64); ok {
		if cpuUsage > a.thresholds["cpu_usage"] {
			indicators = append(indicators, types.ThreatIndicator{
				Type:        "high_cpu_usage",
				Description: "CPU usage exceeds threshold",
				Severity:    a.getSeverity(cpuUsage, a.thresholds["cpu_usage"]),
				Score:       a.calculateScore(cpuUsage, a.thresholds["cpu_usage"]),
				Timestamp:   now,
				Tags:        []string{"performance", "resource_usage"},
			})
		}
	}

	// Add memory analysis
	if memUsage, ok := metrics["memory_usage_percent"].(float64); ok {
		if memUsage > a.thresholds["memory_usage"] {
			indicators = append(indicators, types.ThreatIndicator{
				Type:        "high_memory_usage",
				Description: "Memory usage exceeds threshold",
				Severity:    a.getSeverity(memUsage, a.thresholds["memory_usage"]),
				Score:       a.calculateScore(memUsage, a.thresholds["memory_usage"]),
				Timestamp:   now,
				Tags:        []string{"performance", "resource_usage"},
			})
		}
	}

	return indicators
}

func (a *Analyzer) getSeverity(value, threshold float64) string {
	excess := value - threshold
	switch {
	case excess > 20:
		return "high"
	case excess > 10:
		return "medium"
	default:
		return "low"
	}
}

func (a *Analyzer) calculateScore(value, threshold float64) float64 {
	// Normalize to 0-100 scale and apply logarithmic scaling
	excess := value - threshold
	if excess <= 0 {
		return 0
	}
	return math.Min(100, math.Log10(1+excess)*30)
}
