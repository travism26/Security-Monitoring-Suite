package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rickjms/security-monitoring-suite/mini-xdr/internal/domain"
	"go.uber.org/zap"
)

// EventProcessor defines the interface for processing specific event types
type EventProcessor interface {
	ProcessEvent(ctx context.Context, event *domain.Event) error
	SupportsEventType(eventType domain.EventType) bool
}

// EventService handles the processing of security events
type EventService struct {
	logger     *zap.Logger
	processors []EventProcessor
	mu         sync.RWMutex
	metrics    struct {
		processedCount   int64
		lastProcessedAt  time.Time
		processingErrors int64
	}
}

// NewEventService creates a new instance of EventService
func NewEventService(logger *zap.Logger, processors ...EventProcessor) *EventService {
	return &EventService{
		logger:     logger,
		processors: processors,
	}
}

// RegisterProcessor adds a new event processor
func (s *EventService) RegisterProcessor(processor EventProcessor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processors = append(s.processors, processor)
}

// HandleEvent processes incoming events and routes them to appropriate processors
func (s *EventService) HandleEvent(ctx context.Context, event *domain.Event) error {
	start := time.Now()
	s.logger.Info("Processing event",
		zap.String("event_id", event.ID),
		zap.String("event_type", string(event.Type)),
		zap.String("source", event.Source))

	// Find matching processors
	var matchingProcessors []EventProcessor
	s.mu.RLock()
	for _, p := range s.processors {
		if p.SupportsEventType(event.Type) {
			matchingProcessors = append(matchingProcessors, p)
		}
	}
	s.mu.RUnlock()

	if len(matchingProcessors) == 0 {
		s.logger.Warn("No processors found for event type",
			zap.String("event_type", string(event.Type)))
		return fmt.Errorf("no processors available for event type: %s", event.Type)
	}

	// Process the event with all matching processors
	var processingErrors []error
	for _, processor := range matchingProcessors {
		if err := processor.ProcessEvent(ctx, event); err != nil {
			processingErrors = append(processingErrors, err)
			s.logger.Error("Processor failed to handle event",
				zap.Error(err),
				zap.String("event_id", event.ID),
				zap.String("event_type", string(event.Type)))
		}
	}

	// Update metrics
	s.mu.Lock()
	s.metrics.processedCount++
	s.metrics.lastProcessedAt = time.Now()
	if len(processingErrors) > 0 {
		s.metrics.processingErrors += int64(len(processingErrors))
	}
	s.mu.Unlock()

	s.logger.Info("Event processing completed",
		zap.String("event_id", event.ID),
		zap.Duration("processing_time", time.Since(start)),
		zap.Int("processor_count", len(matchingProcessors)),
		zap.Int("error_count", len(processingErrors)))

	if len(processingErrors) > 0 {
		return fmt.Errorf("event processing completed with %d errors", len(processingErrors))
	}

	return nil
}

// GetMetrics returns current processing metrics
func (s *EventService) GetMetrics() (int64, time.Time, int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics.processedCount, s.metrics.lastProcessedAt, s.metrics.processingErrors
}

// SystemMetricsProcessor handles system metric events
type SystemMetricsProcessor struct {
	logger *zap.Logger
}

func NewSystemMetricsProcessor(logger *zap.Logger) *SystemMetricsProcessor {
	return &SystemMetricsProcessor{logger: logger}
}

func (p *SystemMetricsProcessor) ProcessEvent(ctx context.Context, event *domain.Event) error {
	p.logger.Info("Processing system metrics event",
		zap.String("event_id", event.ID),
		zap.Any("metadata", event.Metadata))
	// TODO: Implement system metrics specific processing logic
	return nil
}

func (p *SystemMetricsProcessor) SupportsEventType(eventType domain.EventType) bool {
	return eventType == domain.EventTypeSystemMetric
}

// NetworkTrafficProcessor handles network traffic events
type NetworkTrafficProcessor struct {
	logger *zap.Logger
}

func NewNetworkTrafficProcessor(logger *zap.Logger) *NetworkTrafficProcessor {
	return &NetworkTrafficProcessor{logger: logger}
}

func (p *NetworkTrafficProcessor) ProcessEvent(ctx context.Context, event *domain.Event) error {
	p.logger.Info("Processing network traffic event",
		zap.String("event_id", event.ID),
		zap.Any("metadata", event.Metadata))
	// TODO: Implement network traffic specific processing logic
	return nil
}

func (p *NetworkTrafficProcessor) SupportsEventType(eventType domain.EventType) bool {
	return eventType == domain.EventTypeNetworkTraffic
}
