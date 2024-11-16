package threat

import "time"

type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

type ThreatIndicator struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Severity    Severity  `json:"severity"`
	Score       float64   `json:"score"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    Metadata  `json:"metadata"`
}

type Metadata struct {
	ProcessID   int32             `json:"process_id,omitempty"`
	ProcessName string            `json:"process_name,omitempty"`
	Username    string            `json:"username,omitempty"`
	Command     string            `json:"command,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

type ThreatScore struct {
	Overall float64                   `json:"overall_score"`
	Details map[string]ComponentScore `json:"component_scores"`
}

type ComponentScore struct {
	Score       float64 `json:"score"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}
