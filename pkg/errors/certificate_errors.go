package errors

import (
	"fmt"
	"time"
)

// ErrorType represents the category of certificate error
type ErrorType string

// Core error types for certificate issues
const (
	ErrCertNotFound    ErrorType = "CERT_NOT_FOUND"    // Certificate file missing
	ErrCertInvalid     ErrorType = "CERT_INVALID"      // Malformed certificate
	ErrCertExpired     ErrorType = "CERT_EXPIRED"      // Past expiration date
	ErrCertUntrusted   ErrorType = "CERT_UNTRUSTED"    // Not in trust store
	ErrCertMismatch    ErrorType = "CERT_MISMATCH"     // Domain/registry mismatch
	ErrCertPermission  ErrorType = "CERT_PERMISSION"   // Access denied
	ErrCertChainBroken ErrorType = "CERT_CHAIN_BROKEN" // Incomplete chain
	ErrCertFormat      ErrorType = "CERT_FORMAT"       // Unsupported format
)

// Severity levels for certificate errors
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityCritical Severity = "CRITICAL"
)

// CertificateError represents a certificate-related error with rich context
type CertificateError struct {
	Type        ErrorType         // Error type classification
	Message     string            // Human-readable error message
	Details     map[string]string // Additional error details
	Resolution  string            // Resolution guidance
	Severity    Severity          // Error severity level
	Timestamp   time.Time         // When error occurred
	Component   string            // Component that generated error
	Operation   string            // Operation being attempted
	OriginalErr error             // Wrapped original error
}

// Error implements the error interface
func (e *CertificateError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s: %s (original: %v)", e.Type, e.Message, e.OriginalErr)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap implements error unwrapping for compatibility with errors.Is/As
func (e *CertificateError) Unwrap() error {
	return e.OriginalErr
}

// WithDetail adds a detail key-value pair to the error
func (e *CertificateError) WithDetail(key, value string) *CertificateError {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// WithResolution sets the resolution guidance for the error
func (e *CertificateError) WithResolution(resolution string) *CertificateError {
	e.Resolution = resolution
	return e
}

// WithComponent sets the component that generated the error
func (e *CertificateError) WithComponent(component string) *CertificateError {
	e.Component = component
	return e
}

// WithOperation sets the operation being attempted when error occurred
func (e *CertificateError) WithOperation(operation string) *CertificateError {
	e.Operation = operation
	return e
}

// WithSeverity sets the severity level of the error
func (e *CertificateError) WithSeverity(severity Severity) *CertificateError {
	e.Severity = severity
	return e
}

// Wrap wraps an existing error with certificate error context
func (e *CertificateError) Wrap(err error) *CertificateError {
	e.OriginalErr = err
	return e
}

// NewCertificateError creates a new certificate error with basic information
func NewCertificateError(errorType ErrorType, message string) *CertificateError {
	return &CertificateError{
		Type:      errorType,
		Message:   message,
		Details:   make(map[string]string),
		Severity:  SeverityError, // Default severity
		Timestamp: time.Now(),
	}
}

// Constructor functions for each error type

// NewCertNotFound creates a certificate not found error
func NewCertNotFound(path string) *CertificateError {
	return NewCertificateError(ErrCertNotFound, fmt.Sprintf("Certificate not found at: %s", path)).
		WithDetail("path", path).
		WithSeverity(SeverityError)
}

// NewCertInvalid creates a certificate invalid error
func NewCertInvalid(reason string) *CertificateError {
	return NewCertificateError(ErrCertInvalid, fmt.Sprintf("Certificate validation failed: %s", reason)).
		WithDetail("reason", reason).
		WithSeverity(SeverityError)
}

// NewCertExpired creates a certificate expired error
func NewCertExpired(expiryDate time.Time) *CertificateError {
	daysAgo := int(time.Since(expiryDate).Hours() / 24)
	return NewCertificateError(ErrCertExpired, fmt.Sprintf("Certificate expired %d days ago", daysAgo)).
		WithDetail("expiry_date", expiryDate.Format(time.RFC3339)).
		WithDetail("days_ago", fmt.Sprintf("%d", daysAgo)).
		WithSeverity(SeverityError)
}

// NewCertUntrusted creates a certificate untrusted error
func NewCertUntrusted(issuer string) *CertificateError {
	return NewCertificateError(ErrCertUntrusted, fmt.Sprintf("Certificate issuer '%s' not in trust store", issuer)).
		WithDetail("issuer", issuer).
		WithSeverity(SeverityWarning)
}

// NewCertMismatch creates a certificate mismatch error
func NewCertMismatch(cn, registry string) *CertificateError {
	return NewCertificateError(ErrCertMismatch, fmt.Sprintf("Certificate CN '%s' doesn't match registry '%s'", cn, registry)).
		WithDetail("cn", cn).
		WithDetail("registry", registry).
		WithSeverity(SeverityError)
}

// NewCertPermission creates a certificate permission error
func NewCertPermission(path string, permErr error) *CertificateError {
	return NewCertificateError(ErrCertPermission, fmt.Sprintf("Permission denied accessing certificate at %s", path)).
		WithDetail("path", path).
		WithSeverity(SeverityError).
		Wrap(permErr)
}

// NewCertChainBroken creates a certificate chain broken error
func NewCertChainBroken(missingLink string) *CertificateError {
	return NewCertificateError(ErrCertChainBroken, fmt.Sprintf("Certificate chain incomplete: missing %s", missingLink)).
		WithDetail("missing_link", missingLink).
		WithSeverity(SeverityError)
}

// NewCertFormat creates a certificate format error
func NewCertFormat(expected, actual string) *CertificateError {
	return NewCertificateError(ErrCertFormat, fmt.Sprintf("Unsupported certificate format: expected %s, got %s", expected, actual)).
		WithDetail("expected", expected).
		WithDetail("actual", actual).
		WithSeverity(SeverityError)
}