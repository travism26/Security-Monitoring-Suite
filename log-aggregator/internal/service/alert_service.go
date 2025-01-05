package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/travism26/log-aggregator/internal/domain"
)

// AlertThresholds defines the threshold values for different metrics
type AlertThresholds struct {
	CPUUsagePercent    float64
	MemoryUsagePercent float64
	ProcessCount       int
}

// DefaultThresholds provides default values for alert thresholds
var DefaultThresholds = AlertThresholds{
	CPUUsagePercent:    80.0, // Alert when CPU usage is above 80%
	MemoryUsagePercent: 85.0, // Alert when memory usage is above 85%
	ProcessCount:       1000, // Alert when process count exceeds 1000
}

type AlertService struct {
	repo       domain.AlertRepository
	thresholds AlertThresholds
}

func NewAlertService(repo domain.AlertRepository) *AlertService {
	return &AlertService{
		repo:       repo,
		thresholds: DefaultThresholds,
	}
}

// SetThresholds allows customizing alert thresholds
func (s *AlertService) SetThresholds(thresholds AlertThresholds) {
	s.thresholds = thresholds
}

// ProcessMetrics evaluates system metrics and generates alerts if thresholds are exceeded
func (s *AlertService) ProcessMetrics(log *domain.Log) error {
	alerts := make([]*domain.Alert, 0)

	// Check CPU usage
	if log.TotalCPUPercent > s.thresholds.CPUUsagePercent {
		alert := s.createAlert(
			fmt.Sprintf("High CPU Usage on %s", log.Host),
			fmt.Sprintf("CPU usage is %.2f%%, which exceeds the threshold of %.2f%%",
				log.TotalCPUPercent, s.thresholds.CPUUsagePercent),
			domain.SeverityHigh,
			log.Host,
			[]string{log.ID},
			map[string]interface{}{
				"cpu_usage": log.TotalCPUPercent,
				"threshold": s.thresholds.CPUUsagePercent,
			},
		)
		alerts = append(alerts, alert)
	}

	// Check memory usage
	memoryUsagePercent := float64(log.TotalMemoryUsage) / float64(s.getSystemTotalMemory()) * 100
	if memoryUsagePercent > s.thresholds.MemoryUsagePercent {
		alert := s.createAlert(
			fmt.Sprintf("High Memory Usage on %s", log.Host),
			fmt.Sprintf("Memory usage is %.2f%%, which exceeds the threshold of %.2f%%",
				memoryUsagePercent, s.thresholds.MemoryUsagePercent),
			domain.SeverityHigh,
			log.Host,
			[]string{log.ID},
			map[string]interface{}{
				"memory_usage": memoryUsagePercent,
				"threshold":    s.thresholds.MemoryUsagePercent,
			},
		)
		alerts = append(alerts, alert)
	}

	// Check process count
	if log.ProcessCount > s.thresholds.ProcessCount {
		alert := s.createAlert(
			fmt.Sprintf("High Process Count on %s", log.Host),
			fmt.Sprintf("Process count is %d, which exceeds the threshold of %d",
				log.ProcessCount, s.thresholds.ProcessCount),
			domain.SeverityMedium,
			log.Host,
			[]string{log.ID},
			map[string]interface{}{
				"process_count": log.ProcessCount,
				"threshold":     s.thresholds.ProcessCount,
			},
		)
		alerts = append(alerts, alert)
	}

	// Store generated alerts
	for _, alert := range alerts {
		if err := s.repo.Store(alert); err != nil {
			return fmt.Errorf("failed to store alert: %w", err)
		}
	}

	return nil
}

// createAlert is a helper function to create a new alert
func (s *AlertService) createAlert(
	title string,
	description string,
	severity domain.AlertSeverity,
	source string,
	relatedLogs []string,
	metadata map[string]interface{},
) *domain.Alert {
	now := time.Now()
	return &domain.Alert{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Severity:    severity,
		Status:      domain.StatusOpen,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
		RelatedLogs: relatedLogs,
		Metadata:    metadata,
	}
}

// getSystemTotalMemory returns the total system memory in bytes
// This is a placeholder - in a real implementation, this would get the actual system memory
func (s *AlertService) getSystemTotalMemory() int64 {
	return 16 * 1024 * 1024 * 1024 // 16GB as example
}

// UpdateAlertStatus updates the status of an alert
func (s *AlertService) UpdateAlertStatus(id string, status domain.AlertStatus) error {
	alert, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find alert: %w", err)
	}

	alert.Status = status
	alert.UpdatedAt = time.Now()

	if status == domain.StatusResolved {
		now := time.Now()
		alert.ResolvedAt = &now
	}

	return s.repo.Update(alert)
}

// GetAlert retrieves an alert by ID
func (s *AlertService) GetAlert(id string) (*domain.Alert, error) {
	return s.repo.FindByID(id)
}

// ListAlerts retrieves alerts with pagination
func (s *AlertService) ListAlerts(limit, offset int) ([]*domain.Alert, error) {
	return s.repo.List(limit, offset)
}

// ListAlertsByStatus retrieves alerts filtered by status
func (s *AlertService) ListAlertsByStatus(status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error) {
	return s.repo.FindByStatus(status, limit, offset)
}

// ListAlertsBySeverity retrieves alerts filtered by severity
func (s *AlertService) ListAlertsBySeverity(severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error) {
	return s.repo.FindBySeverity(severity, limit, offset)
}
