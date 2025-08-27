<<<<<<< HEAD
=======
// Package v2 provides OCI certificate management interfaces and types.
// This file contains the CertificateService interface and related types
// that should be implemented by certificate service providers.
>>>>>>> origin/idpbuidler-oci-mgmt/phase3/wave1/E3.1.3-certificate-validator
package v2

import (
	"context"
<<<<<<< HEAD
	"crypto/tls"
	"crypto/x509"
	"time"
)

// CertificateService defines the contract for certificate management
type CertificateService interface {
	// LoadCertificateBundle loads certificates from various formats
	LoadCertificateBundle(ctx context.Context, path string, format CertFormat) (*CertBundle, error)

	// SetVerificationMode configures how certificates are verified
	SetVerificationMode(mode VerificationMode) error

	// ValidateCertificate validates a single certificate
	ValidateCertificate(cert *x509.Certificate) error

	// LoadGiteaCertificate auto-discovers and loads Gitea certificates
	LoadGiteaCertificate(ctx context.Context, giteaURL string) (*CertBundle, error)

	// GetTLSConfig returns configured TLS configuration
	GetTLSConfig() (*tls.Config, error)

	// AddCACertificate adds a CA to the certificate pool
	AddCACertificate(cert *x509.Certificate) error

	// RemoveCACertificate removes a CA from the pool
	RemoveCACertificate(cert *x509.Certificate) error

	// ListCertificates returns all managed certificates
	ListCertificates() ([]*x509.Certificate, error)

	// RotateCertificate handles certificate rotation
	RotateCertificate(old, new *x509.Certificate) error

	// GetCertificatePool returns the current certificate pool
	GetCertificatePool() *x509.CertPool
}

// CertFormat represents supported certificate formats
type CertFormat string

const (
	CertFormatPEM    CertFormat = "pem"
	CertFormatDER    CertFormat = "der"
	CertFormatPKCS7  CertFormat = "pkcs7"
	CertFormatPKCS12 CertFormat = "pkcs12"
)

// VerificationMode defines certificate verification strategies
type VerificationMode string

const (
	VerificationModeStrict     VerificationMode = "strict"
	VerificationModePermissive VerificationMode = "permissive"
	VerificationModeSkip       VerificationMode = "skip"
	VerificationModeCustomCA   VerificationMode = "custom-ca"
)

// CertBundle represents a collection of certificates
type CertBundle struct {
	Certificates []*x509.Certificate
	CAs          []*x509.Certificate
	Format       CertFormat
	LoadedAt     time.Time
	Source       string
}

// CertificateConfig holds certificate service configuration
type CertificateConfig struct {
	BundlePath           string
	VerificationMode     VerificationMode
	AutoDiscoverGitea    bool
	SkipVerifyFallback   bool
	RefreshInterval      time.Duration
	CustomCAPath         string
}

// CertificateError represents certificate-specific errors
type CertificateError struct {
	Code    string
	Message string
	Cert    *x509.Certificate
	Err     error
}

func (e *CertificateError) Error() string {
	return e.Message
=======
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
>>>>>>> origin/idpbuidler-oci-mgmt/phase3/wave1/E3.1.3-certificate-validator
}