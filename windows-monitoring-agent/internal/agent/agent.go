// internal/agent/agent.go
package agent

import (
	"log"
	"time"

	"github.com/travism26/windows-monitoring-agent/internal/config"
	"github.com/travism26/windows-monitoring-agent/internal/exporter"
	"github.com/travism26/windows-monitoring-agent/internal/metrics"
)

type Agent struct {
	config   *config.Config
	metrics  *metrics.MetricsCollector
	exporter *exporter.Exporter
	interval time.Duration
}

func NewAgent(cfg *config.Config, mc *metrics.MetricsCollector, exp *exporter.Exporter) *Agent {
	return &Agent{
		config:   cfg,
		metrics:  mc,
		exporter: exp,
		interval: time.Duration(cfg.Interval) * time.Second,
	}
}

func (a *Agent) Start() {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		data := a.metrics.Collect()
		if err := a.exporter.Export(data); err != nil {
			log.Printf("Error exporting metrics: %v", err)
		}
	}
}
