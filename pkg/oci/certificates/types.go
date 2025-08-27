// Package certificates provides certificate management functionality for OCI trust stores.
package certificates

import (
	"crypto/x509"
	"time"
)

// CertificateStatus represents the current status of a certificate.
type CertificateStatus string

const (
	CertificateStatusActive    CertificateStatus = "active"
	CertificateStatusExpiring  CertificateStatus = "expiring"
	CertificateStatusExpired   CertificateStatus = "expired"
	CertificateStatusRevoked   CertificateStatus = "revoked"
	CertificateStatusPending   CertificateStatus = "pending"
)

// Certificate represents a certificate with its metadata and validation status.
type Certificate struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Certificate      *x509.Certificate     `json:"-"`
	PEM              []byte                `json:"pem"`
	Status           CertificateStatus     `json:"status"`
	ValidFrom        time.Time             `json:"valid_from"`
	ValidTo          time.Time             `json:"valid_to"`
	Issuer           string                `json:"issuer"`
	Subject          string                `json:"subject"`
	KeyUsage         []string              `json:"key_usage"`
	ExtendedKeyUsage []string              `json:"extended_key_usage"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	Tags             map[string]string     `json:"tags"`
}

// EventType represents the type of certificate-related event.
type EventType string

const (
	EventCertificateAdded    EventType = "certificate_added"
	EventCertificateUpdated  EventType = "certificate_updated"
	EventCertificateRemoved  EventType = "certificate_removed"
	EventCertificateExpiring EventType = "certificate_expiring"
	EventCertificateExpired  EventType = "certificate_expired"
	EventValidationFailed    EventType = "validation_failed"
	EventPoolUpdated         EventType = "pool_updated"
)

// Event represents a certificate-related event for change notifications.
type Event struct {
	Type          EventType              `json:"type"`
	CertificateID string                 `json:"certificate_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Message       string                 `json:"message"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ValidationRule represents a validation rule that can be applied to certificates.
type ValidationRule struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Critical    bool                   `json:"critical"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ValidationResult represents the result of validating a certificate.
type ValidationResult struct {
	Valid         bool                 `json:"valid"`
	CertificateID string               `json:"certificate_id"`
	ValidatedAt   time.Time            `json:"validated_at"`
	Errors        []ValidationError    `json:"errors,omitempty"`
	Warnings      []ValidationWarning  `json:"warnings,omitempty"`
}

// ValidationWarning represents a non-critical validation issue.
type ValidationWarning struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Common validation error variables for consistency.
var (
	ErrCertificateNotFound = NewValidationError("CERT_NOT_FOUND", "certificate not found")
	ErrCertificateExpired  = NewValidationError("CERT_EXPIRED", "certificate has expired")
	ErrInvalidPEM          = NewValidationError("INVALID_PEM", "invalid PEM data")
	ErrUnsupportedKeyUsage = NewValidationError("UNSUPPORTED_KEY_USAGE", "unsupported key usage")
)