package domain

import (
	"time"
)

// Severity represents the severity level of an event
type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

// EventType represents the type of security event
type EventType string

const (
	EventTypeSystemMetric    EventType = "SYSTEM_METRIC"
	EventTypeNetworkTraffic  EventType = "NETWORK_TRAFFIC"
	EventTypeSecurityAlert   EventType = "SECURITY_ALERT"
	EventTypeProcessActivity EventType = "PROCESS_ACTIVITY"
)

// Event represents a normalized security event in the XDR system
type Event struct {
	ID          string            `json:"id"`
	Timestamp   time.Time         `json:"timestamp"`
	Type        EventType         `json:"type"`
	Source      string            `json:"source"`
	Severity    Severity          `json:"severity"`
	Description string            `json:"description"`
	RawData     []byte            `json:"raw_data,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Validate performs basic validation on the event
func (e *Event) Validate() error {
	if e.ID == "" {
		return NewValidationError("event ID is required")
	}
	if e.Source == "" {
		return NewValidationError("event source is required")
	}
	if e.Type == "" {
		return NewValidationError("event type is required")
	}
	if e.Timestamp.IsZero() {
		return NewValidationError("event timestamp is required")
	}
	return nil
}

// ValidationError represents an error that occurs during event validation
type ValidationError struct {
	Message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func (e *ValidationError) Error() string {
	return e.Message
}
