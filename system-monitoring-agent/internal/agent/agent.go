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
	// Use tenant-specific sample rate if configured, otherwise use global interval
	interval := cfg.Interval
	if cfg.Tenant.CollectionRules.SampleRate > 0 {
		interval = cfg.Tenant.CollectionRules.SampleRate
	}

	return &Agent{
		config:    cfg,
		metrics:   mc,
		exporters: exporters,
		interval:  time.Duration(interval) * time.Second,
	}
}

// validateTenantContext checks if the tenant context is valid
func (a *Agent) validateTenantContext() error {
	// Temporarily disabled tenant ID requirement
	// API key is now optional
	return nil
}

func (a *Agent) Start(done chan struct{}) {
	// Validate tenant context before starting
	if err := a.validateTenantContext(); err != nil {
		log.Printf("Error validating tenant context: %v", err)
		return
	}

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			data := a.metrics.Collect()

			// Export metrics with retry logic
			for _, exp := range a.exporters {
				retries := 0
				maxRetries := a.config.HTTP.RetryAttempts

				for retries < maxRetries {
					if err := exp.Export(data); err != nil {
						if retries == maxRetries-1 {
							log.Printf("Error exporting metrics after %d retries: %v", maxRetries, err)
							break
						}
						retries++
						time.Sleep(time.Duration(a.config.HTTP.RetryDelay) * time.Second)
						continue
					}
					break
				}
			}
		}
	}
}
