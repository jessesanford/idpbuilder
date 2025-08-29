// Package audit provides comprehensive audit logging functionality for certificate validation operations.
// This package implements structured audit logging with persistence capabilities and integration
// with the certificate chain validation system.
package audit

import (
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// AuditLogger defines the interface for certificate validation audit logging.
// Implementations of this interface provide structured audit trails for
// certificate validation operations with persistence and querying capabilities.
type AuditLogger interface {
	// LogValidation logs a successful or failed certificate validation operation.
	// The validation result contains all validation details, and metadata provides
	// additional context such as request ID, client information, etc.
	LogValidation(result *certs.ValidationResult, metadata map[string]interface{}) error

	// LogFailure logs a system failure that occurred during certificate validation.
	// This is used for errors outside of normal validation failures (e.g., I/O errors,
	// configuration issues, etc.)
	LogFailure(err error, context map[string]interface{}) error

	// LogEvent logs a general audit event related to certificate operations.
	// This is used for operational events like trust anchor updates, configuration changes, etc.
	LogEvent(eventType string, message string, metadata map[string]interface{}) error

	// Flush ensures all pending audit records are written to persistent storage.
	// This method should be called before application shutdown or at regular intervals.
	Flush() error

	// Close gracefully shuts down the audit logger, ensuring all pending records
	// are flushed and resources are released.
	Close() error
}

// AuditRecord represents a single audit log entry with structured metadata.
// All audit events are normalized into this format for consistent storage and querying.
type AuditRecord struct {
	// ID is a unique identifier for this audit record
	ID string `json:"id"`

	// Timestamp is the exact time when this audit event occurred
	Timestamp time.Time `json:"timestamp"`

	// EventType categorizes the type of audit event (validation, failure, system_event)
	EventType AuditEventType `json:"event_type"`

	// Level indicates the severity/importance of the audit event
	Level AuditLevel `json:"level"`

	// Message provides a human-readable description of the audit event
	Message string `json:"message"`

	// ValidationResult contains the validation details if this is a validation audit
	ValidationResult *certs.ValidationResult `json:"validation_result,omitempty"`

	// Error contains error details if this is a failure audit
	Error *AuditError `json:"error,omitempty"`

	// Metadata contains additional contextual information as key-value pairs
	Metadata map[string]interface{} `json:"metadata"`

	// CorrelationID links related audit events together (e.g., validation attempts for same certificate)
	CorrelationID string `json:"correlation_id,omitempty"`
}

// AuditEventType defines the different categories of audit events.
type AuditEventType string

const (
	// EventTypeValidation represents certificate validation events
	EventTypeValidation AuditEventType = "validation"

	// EventTypeFailure represents system failures during certificate operations
	EventTypeFailure AuditEventType = "failure"

	// EventTypeSystemEvent represents operational events (config changes, trust anchor updates)
	EventTypeSystemEvent AuditEventType = "system_event"
)

// AuditLevel defines the severity levels for audit events.
type AuditLevel string

const (
	// LevelInfo for informational audit events (successful validations)
	LevelInfo AuditLevel = "info"

	// LevelWarn for warning events (validation failures, deprecated usage)
	LevelWarn AuditLevel = "warn"

	// LevelError for error events (system failures, configuration errors)
	LevelError AuditLevel = "error"

	// LevelCritical for critical events (security violations, system compromises)
	LevelCritical AuditLevel = "critical"
)

// AuditError contains structured error information for audit logging.
type AuditError struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Stack   string                 `json:"stack,omitempty"`
}