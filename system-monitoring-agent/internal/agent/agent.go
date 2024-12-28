// internal/agent/agent.go
package agent

import (
	"log"
	"time"

	"github.com/travism26/system-monitoring-agent/internal/config"
	"github.com/travism26/system-monitoring-agent/internal/exporter"
	"github.com/travism26/system-monitoring-agent/internal/metrics"
)

type Agent struct {
	config    *config.Config
	metrics   *metrics.MetricsCollector
	exporters []exporter.MetricsExporter
	interval  time.Duration
}

func NewAgent(cfg *config.Config, mc *metrics.MetricsCollector, exporters ...exporter.MetricsExporter) *Agent {
	return &Agent{
		config:    cfg,
		metrics:   mc,
		exporters: exporters,
		interval:  time.Duration(cfg.Interval) * time.Second,
	}
}

func (a *Agent) Start(done chan struct{}) {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			data := a.metrics.Collect()
			for _, exp := range a.exporters {
				if err := exp.Export(data); err != nil {
					log.Printf("Error exporting metrics: %v", err)
				}
			}
		}
	}
}
