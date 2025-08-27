// Package v2 provides OCI certificate management interfaces and types.
// Imported from E3.1.1 contracts to maintain API compatibility.
package v2

import (
	"context"
	"crypto/x509"
)

// VerificationMode defines the certificate verification behavior
type VerificationMode string

const (
	// VerificationModeStrict uses only system certificate pool
	VerificationModeStrict VerificationMode = "strict"
	// VerificationModeCustomCA uses system and custom certificate pools
	VerificationModeCustomCA VerificationMode = "custom-ca"
	// VerificationModeSkip skips certificate verification (development only)
	VerificationModeSkip VerificationMode = "skip"
)

// CertificateInfo represents metadata about a certificate
type CertificateInfo struct {
	Subject     string `json:"subject"`
	Issuer      string `json:"issuer"`
	NotBefore   string `json:"not_before"`
	NotAfter    string `json:"not_after"`
	KeyUsage    string `json:"key_usage"`
	Fingerprint string `json:"fingerprint"`
}

// ValidationResult represents the result of certificate validation
type ValidationResult struct {
	Valid       bool     `json:"valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	Certificate *CertificateInfo `json:"certificate"`
}

// CertificateService defines the interface for certificate management operations
type CertificateService interface {
	// LoadCertificateBundle loads certificates from the specified bundle path
	LoadCertificateBundle(ctx context.Context, bundlePath string) error

	// SetVerificationMode sets the certificate verification mode
	SetVerificationMode(ctx context.Context, mode VerificationMode) error

	// ValidateCertificate validates a certificate against the current configuration
	ValidateCertificate(ctx context.Context, cert *x509.Certificate) (*ValidationResult, error)

	// LoadGiteaCertificate loads Gitea-specific certificates
	LoadGiteaCertificate(ctx context.Context, giteaURL string) error

	// GetCertPool returns the current certificate pool (thread-safe)
	GetCertPool() *x509.CertPool

	// AddCertificate adds a certificate to the pool
	AddCertificate(ctx context.Context, cert *x509.Certificate) error

	// RemoveCertificate removes a certificate from the pool
	RemoveCertificate(ctx context.Context, fingerprint string) error
}