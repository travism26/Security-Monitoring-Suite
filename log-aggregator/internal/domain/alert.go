package domain

import "time"

type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "LOW"
	SeverityMedium   AlertSeverity = "MEDIUM"
	SeverityHigh     AlertSeverity = "HIGH"
	SeverityCritical AlertSeverity = "CRITICAL"
)

// IsValid checks if the alert severity is a valid value
func (s AlertSeverity) IsValid() bool {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	default:
		return false
	}
}

type AlertStatus string

const (
	StatusOpen     AlertStatus = "OPEN"
	StatusResolved AlertStatus = "RESOLVED"
	StatusIgnored  AlertStatus = "IGNORED"
)

// IsValid checks if the alert status is a valid value
func (s AlertStatus) IsValid() bool {
	switch s {
	case StatusOpen, StatusResolved, StatusIgnored:
		return true
	default:
		return false
	}
}

type Alert struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Severity    AlertSeverity `json:"severity"`
	Status      AlertStatus   `json:"status"`
	Source      string        `json:"source"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ResolvedAt  *time.Time    `json:"resolved_at,omitempty"`

	// Reference to related logs
	RelatedLogs []string `json:"related_logs,omitempty"`

	// Additional context about the alert
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AlertRepository defines the interface for alert storage operations
type AlertRepository interface {
	Store(alert *Alert) error
	FindByID(id string) (*Alert, error)
	List(limit, offset int) ([]*Alert, error)
	Update(alert *Alert) error
	FindByStatus(status AlertStatus, limit, offset int) ([]*Alert, error)
	FindBySeverity(severity AlertSeverity, limit, offset int) ([]*Alert, error)
}
