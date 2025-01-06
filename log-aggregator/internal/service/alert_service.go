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
	repo         domain.AlertRepository
	thresholds   AlertThresholds
	systemMemory int64
	timeNowFn    func() time.Time
}

// AlertServiceConfig allows customizing service behavior
type AlertServiceConfig struct {
	SystemMemory int64
	TimeNowFn    func() time.Time
}

func NewAlertService(repo domain.AlertRepository, config *AlertServiceConfig) *AlertService {
	if config == nil {
		config = &AlertServiceConfig{
			SystemMemory: 16 * 1024 * 1024 * 1024, // 16GB default
			TimeNowFn:    time.Now,
		}
	}
	return &AlertService{
		repo:         repo,
		thresholds:   DefaultThresholds,
		systemMemory: config.SystemMemory,
		timeNowFn:    config.TimeNowFn,
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
	now := s.timeNowFn()
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

// getSystemTotalMemory returns the configured system memory in bytes
func (s *AlertService) getSystemTotalMemory() int64 {
	return s.systemMemory
}

// UpdateAlertStatus updates the status of an alert
func (s *AlertService) UpdateAlertStatus(id string, status domain.AlertStatus) error {
	existingAlert, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find alert: %w", err)
	}

	// Create a copy of the alert to avoid modifying the original
	updatedAlert := *existingAlert
	updatedAlert.Status = status
	updatedAlert.UpdatedAt = s.timeNowFn()

	if status == domain.StatusResolved {
		now := s.timeNowFn()
		updatedAlert.ResolvedAt = &now
	}

	return s.repo.Update(&updatedAlert)
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

// AlertTrends represents trend analysis data for alerts
type AlertTrends struct {
	TotalAlerts      int            `json:"total_alerts"`
	AlertsBySeverity map[string]int `json:"alerts_by_severity"`
	AlertsByStatus   map[string]int `json:"alerts_by_status"`
	TimeDistribution map[string]int `json:"time_distribution"` // Hourly distribution
	TopSources       map[string]int `json:"top_sources"`
}

// GetAlertTrends analyzes alerts within a time range to provide trend information
func (s *AlertService) GetAlertTrends(start, end time.Time) (*AlertTrends, error) {
	// Get all alerts within the time range
	// Note: This is a simplified implementation. In production, you'd want to:
	// 1. Add a repository method to get alerts by time range
	// 2. Add proper pagination to handle large datasets
	// 3. Potentially use database aggregations for better performance
	alerts, err := s.repo.List(1000, 0) // Get last 1000 alerts as a sample
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve alerts: %w", err)
	}

	trends := &AlertTrends{
		AlertsBySeverity: make(map[string]int),
		AlertsByStatus:   make(map[string]int),
		TimeDistribution: make(map[string]int),
		TopSources:       make(map[string]int),
	}

	for _, alert := range alerts {
		// Only consider alerts within the time range
		if alert.CreatedAt.Before(start) || alert.CreatedAt.After(end) {
			continue
		}

		trends.TotalAlerts++

		// Count by severity
		trends.AlertsBySeverity[string(alert.Severity)]++

		// Count by status
		trends.AlertsByStatus[string(alert.Status)]++

		// Count by hour of day (for time distribution analysis)
		hour := alert.CreatedAt.Format("15:00")
		trends.TimeDistribution[hour]++

		// Count by source
		trends.TopSources[alert.Source]++
	}

	return trends, nil
}
