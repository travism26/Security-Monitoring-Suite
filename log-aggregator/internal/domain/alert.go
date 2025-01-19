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
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organization_id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Severity       AlertSeverity          `json:"severity"`
	Status         AlertStatus            `json:"status"`
	Source         string                 `json:"source"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	RelatedLogs    []string               `json:"related_logs,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// AlertRepository defines the interface for alert storage operations
type AlertRepository interface {
	// Core operations
	Store(alert *Alert) error
	Update(alert *Alert) error
	Delete(orgID, id string) error

	// Retrieval operations
	FindByID(orgID, id string) (*Alert, error)
	List(orgID string, limit, offset int) ([]*Alert, error)

	// Status-based queries
	FindByStatus(orgID string, status AlertStatus, limit, offset int) ([]*Alert, error)
	CountByStatus(orgID string, status AlertStatus) (int64, error)

	// Severity-based queries
	FindBySeverity(orgID string, severity AlertSeverity, limit, offset int) ([]*Alert, error)
	CountBySeverity(orgID string, severity AlertSeverity) (int64, error)

	// Time-based queries
	ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*Alert, error)
	CountByTimeRange(orgID string, start, end time.Time) (int64, error)

	// Source-based queries
	FindBySource(orgID string, source string, limit, offset int) ([]*Alert, error)
	CountBySource(orgID string, source string) (int64, error)
}
