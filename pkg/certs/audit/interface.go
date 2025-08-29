package audit

import "time"

// AuditLoggerInterface defines the interface for audit logging implementations
type AuditLoggerInterface interface {
	// LogCertificateValidation logs certificate validation decisions
	LogCertificateValidation(cert, hostname, decision, reason, userID string) error
	
	// LogTrustDecision logs trust anchor decisions
	LogTrustDecision(cert, decision, reason, userID string) error
	
	// LogFallbackActivation logs when fallback mechanisms are activated
	LogFallbackActivation(cert, fallbackType, reason, userID string) error
	
	// LogSecurityOverride logs when security policies are overridden
	LogSecurityOverride(cert, overrideType, justification, userID string) error
	
	// LogError logs certificate validation errors
	LogError(cert, errorType, details, userID string) error
	
	// ReadEntries reads audit entries from the log (limit 0 = all entries)
	ReadEntries(limit int) ([]AuditEntry, error)
	
	// GetRecentEntries returns recent audit entries within the specified duration
	GetRecentEntries(since time.Duration) ([]AuditEntry, error)
	
	// Close closes the audit logger
	Close() error
}

// NoOpAuditLogger provides a no-operation implementation for testing
type NoOpAuditLogger struct{}

// NewNoOpAuditLogger creates a new no-operation audit logger
func NewNoOpAuditLogger() *NoOpAuditLogger {
	return &NoOpAuditLogger{}
}

// LogCertificateValidation does nothing
func (n *NoOpAuditLogger) LogCertificateValidation(cert, hostname, decision, reason, userID string) error {
	return nil
}

// LogTrustDecision does nothing
func (n *NoOpAuditLogger) LogTrustDecision(cert, decision, reason, userID string) error {
	return nil
}

// LogFallbackActivation does nothing
func (n *NoOpAuditLogger) LogFallbackActivation(cert, fallbackType, reason, userID string) error {
	return nil
}

// LogSecurityOverride does nothing
func (n *NoOpAuditLogger) LogSecurityOverride(cert, overrideType, justification, userID string) error {
	return nil
}

// LogError does nothing
func (n *NoOpAuditLogger) LogError(cert, errorType, details, userID string) error {
	return nil
}

// ReadEntries returns empty slice
func (n *NoOpAuditLogger) ReadEntries(limit int) ([]AuditEntry, error) {
	return []AuditEntry{}, nil
}

// GetRecentEntries returns empty slice
func (n *NoOpAuditLogger) GetRecentEntries(since time.Duration) ([]AuditEntry, error) {
	return []AuditEntry{}, nil
}

// Close does nothing
func (n *NoOpAuditLogger) Close() error {
	return nil
}