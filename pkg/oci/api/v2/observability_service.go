package v2

import (
	"context"
)

// ObservabilityService defines the contract for observability
type ObservabilityService interface {
	// StartSpan starts a new trace span
	StartSpan(ctx context.Context, name string) (Span, context.Context)

	// RecordMetric records a metric value
	RecordMetric(name string, value float64, labels map[string]string)

	// CreateAlert creates an alert
	CreateAlert(alert Alert) error

	// GetMetrics returns current metrics
	GetMetrics() map[string]interface{}
}

// Span represents a trace span
type Span interface {
	End()
	SetAttribute(key string, value interface{})
	RecordError(err error)
}

// Alert represents an alert configuration
type Alert struct {
	Name      string
	Condition string
	Threshold float64
	Duration  string
	Labels    map[string]string
}