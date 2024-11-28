package logger

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func New(prefix string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, prefix+" ", log.LstdFlags|log.Lmsgprefix),
	}
}

func (l *Logger) LogMetricsSent(scenario string, endpoint string) {
	l.Printf("Metrics sent - Scenario: %s, Endpoint: %s, Time: %s",
		scenario,
		endpoint,
		time.Now().Format(time.RFC3339))
}
