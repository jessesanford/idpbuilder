package errors

import (
	"fmt"
	"time"
)

// ErrorType represents the category of certificate error
type ErrorType string

// Core certificate error types
const (
	ErrCertNotFound    ErrorType = "CERT_NOT_FOUND"
	ErrCertExpired     ErrorType = "CERT_EXPIRED"
	ErrCertUntrusted   ErrorType = "CERT_UNTRUSTED"
	ErrCertMismatch    ErrorType = "CERT_MISMATCH"
	ErrCertPermission  ErrorType = "CERT_PERMISSION"
)

// Severity levels
type Severity string

const (
	SeverityWarning Severity = "WARNING"
	SeverityError   Severity = "ERROR"
)

// CertificateError represents a certificate-related error
type CertificateError struct {
	Type        ErrorType
	Message     string
	Details     map[string]string
	Severity    Severity
	Timestamp   time.Time
	OriginalErr error
}

// Error implements the error interface
func (e *CertificateError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// WithDetail adds detail information
func (e *CertificateError) WithDetail(key, value string) *CertificateError {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// NewCertificateError creates a new certificate error
func NewCertificateError(errorType ErrorType, message string) *CertificateError {
	return &CertificateError{
		Type:      errorType,
		Message:   message,
		Details:   make(map[string]string),
		Severity:  SeverityError,
		Timestamp: time.Now(),
	}
}

// Constructor functions for each error type

func NewCertNotFound(path string) *CertificateError {
	return NewCertificateError(ErrCertNotFound, fmt.Sprintf("Certificate not found at: %s", path)).
		WithDetail("path", path)
}

func NewCertExpired(expiryDate time.Time) *CertificateError {
	daysAgo := int(time.Since(expiryDate).Hours() / 24)
	return NewCertificateError(ErrCertExpired, fmt.Sprintf("Certificate expired %d days ago", daysAgo)).
		WithDetail("expiry_date", expiryDate.Format(time.RFC3339)).
		WithDetail("days_ago", fmt.Sprintf("%d", daysAgo))
}

func NewCertUntrusted(issuer string) *CertificateError {
	return NewCertificateError(ErrCertUntrusted, fmt.Sprintf("Certificate issuer '%s' not in trust store", issuer)).
		WithDetail("issuer", issuer).
		WithSeverity(SeverityWarning)
}

func NewCertMismatch(cn, registry string) *CertificateError {
	return NewCertificateError(ErrCertMismatch, fmt.Sprintf("Certificate CN '%s' doesn't match registry '%s'", cn, registry)).
		WithDetail("cn", cn).
		WithDetail("registry", registry)
}

func NewCertPermission(path string) *CertificateError {
	return NewCertificateError(ErrCertPermission, fmt.Sprintf("Permission denied accessing certificate at %s", path)).
		WithDetail("path", path)
}

// WithSeverity sets the error severity
func (e *CertificateError) WithSeverity(severity Severity) *CertificateError {
	e.Severity = severity
	return e
}