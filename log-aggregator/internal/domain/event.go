package domain

import "time"

type EventType string

const (
	EventTypeAlert  EventType = "ALERT"
	EventTypeMetric EventType = "METRIC"
	EventTypeAudit  EventType = "AUDIT"
	EventTypeSystem EventType = "SYSTEM"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`

	// Reference to related entities
	RelatedID   string `json:"related_id,omitempty"`
	RelatedType string `json:"related_type,omitempty"`

	// Additional event data
	Data map[string]interface{} `json:"data,omitempty"`
}

// EventRepository defines the interface for event storage operations
type EventRepository interface {
	Store(event *Event) error
	FindByID(id string) (*Event, error)
	List(limit, offset int) ([]*Event, error)
	FindByTimeRange(start, end time.Time) ([]*Event, error)
	FindByType(eventType EventType, limit, offset int) ([]*Event, error)
}
