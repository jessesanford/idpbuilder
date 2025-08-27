package v2

import (
	"context"
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
}