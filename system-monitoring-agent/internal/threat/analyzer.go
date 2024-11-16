package threat

import (
	"math"
	"time"
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

func (a *Analyzer) AnalyzeMetrics(metrics map[string]interface{}) []ThreatIndicator {
	var indicators []ThreatIndicator
	now := time.Now()

	// Analyze CPU usage
	if cpuUsage, ok := metrics["cpu_usage"].(float64); ok {
		if cpuUsage > a.thresholds["cpu_usage"] {
			indicators = append(indicators, ThreatIndicator{
				Type:        "high_cpu_usage",
				Description: "CPU usage exceeds threshold",
				Severity:    a.getSeverity(cpuUsage, a.thresholds["cpu_usage"]),
				Score:       a.calculateScore(cpuUsage, a.thresholds["cpu_usage"]),
				Timestamp:   now,
				Metadata: Metadata{
					Tags: []string{"performance", "resource_usage"},
				},
			})
		}
	}

	// Add memory analysis
	if memUsage, ok := metrics["memory_usage_percent"].(float64); ok {
		if memUsage > a.thresholds["memory_usage"] {
			indicators = append(indicators, ThreatIndicator{
				Type:        "high_memory_usage",
				Description: "Memory usage exceeds threshold",
				Severity:    a.getSeverity(memUsage, a.thresholds["memory_usage"]),
				Score:       a.calculateScore(memUsage, a.thresholds["memory_usage"]),
				Timestamp:   now,
				Metadata: Metadata{
					Tags: []string{"performance", "resource_usage"},
				},
			})
		}
	}

	return indicators
}

func (a *Analyzer) getSeverity(value, threshold float64) Severity {
	excess := value - threshold
	switch {
	case excess > 20:
		return SeverityHigh
	case excess > 10:
		return SeverityMedium
	default:
		return SeverityLow
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
