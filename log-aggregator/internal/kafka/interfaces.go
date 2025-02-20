package kafka

import "github.com/travism26/log-aggregator/internal/domain"

// LogService defines the interface for log operations
type LogService interface {
	StoreLog(log *domain.Log) error
	GetLog(userID, id string) (*domain.Log, error)
	ListLogs(userID string, limit, offset int) ([]*domain.Log, error)
}

// AlertService defines the interface for alert operations
type AlertService interface {
	ProcessMetrics(log *domain.Log) error
	GetAlert(id string) (*domain.Alert, error)
	ListAlerts(limit, offset int) ([]*domain.Alert, error)
	ListAlertsByStatus(status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error)
	ListAlertsBySeverity(severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error)
	UpdateAlertStatus(id string, status domain.AlertStatus) error
}

// ProcessRepository defines the interface for process storage operations
// TODO: Rename this to ProcessService to follow other naming conventions
type ProcessRepository interface {
	StoreBatch(processes []domain.Process) error
	FindByLogID(logID string) ([]domain.Process, error)
}
