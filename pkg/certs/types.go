// pkg/certs/types.go
package certs

import (
	"context"
	"crypto/x509"
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// KindCertExtractor defines the interface for extracting certificates from Kind clusters
type KindCertExtractor interface {
	// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
	ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)

	// GetClusterName returns the name of the Kind cluster
	GetClusterName() (string, error)

	// ValidateCertificate performs basic validation on the extracted certificate
	ValidateCertificate(cert *x509.Certificate) error

	// SaveCertificate saves the certificate to the local trust store
	SaveCertificate(cert *x509.Certificate, path string) error
}

// TransportConfig holds configuration options for registry transport
type TransportConfig struct {
	// Timeout for HTTP requests (default: 30 seconds)
	Timeout time.Duration
	
	// MaxIdleConns controls the maximum number of idle connections
	MaxIdleConns int
	
	// MaxIdleConnsPerHost controls the maximum idle connections per host
	MaxIdleConnsPerHost int
	
	// IdleConnTimeout is the maximum amount of time an idle connection will remain idle
	IdleConnTimeout time.Duration
}

// TrustStoreManager manages trusted certificates for registry operations
type TrustStoreManager interface {
	// AddCertificate adds a certificate for a specific registry
	AddCertificate(registry string, cert *x509.Certificate) error

	// RemoveCertificate removes the certificate for a registry
	RemoveCertificate(registry string) error

	// SetInsecureRegistry marks a registry as insecure (skip TLS verification)
	SetInsecureRegistry(registry string, insecure bool) error

	// GetTrustedCerts returns all trusted certificates for a registry
	GetTrustedCerts(registry string) ([]*x509.Certificate, error)

	// GetCertPool returns a configured cert pool for a registry
	GetCertPool(registry string) (*x509.CertPool, error)

	// IsInsecure checks if a registry is marked as insecure
	IsInsecure(registry string) bool

	// LoadFromDisk loads all certificates from persistent storage
	LoadFromDisk() error

	// SaveToDisk saves a certificate to persistent storage
	SaveToDisk(registry string, cert *x509.Certificate) error

	// Transport configuration methods
	// ConfigureTransport creates a remote.Option with proper TLS configuration
	ConfigureTransport(registry string) (remote.Option, error)

	// ConfigureTransportWithConfig creates a remote.Option with custom transport configuration
	ConfigureTransportWithConfig(registry string, config *TransportConfig) (remote.Option, error)

	// CreateHTTPClient creates an HTTP client with proper TLS configuration
	CreateHTTPClient(registry string) (*http.Client, error)

	// CreateHTTPClientWithConfig creates an HTTP client with custom configuration
	CreateHTTPClientWithConfig(registry string, config *TransportConfig) (*http.Client, error)
}

// CertificateInfo contains metadata about an extracted certificate
type CertificateInfo struct {
	Subject   string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
	IsCA      bool
	DNSNames  []string
}

// CertValidator provides comprehensive X.509 certificate validation
type CertValidator interface {
	// ValidateChain verifies the certificate chain against trusted roots
	ValidateChain(cert *x509.Certificate) error
	
	// CheckExpiry checks if certificate is expired or expiring soon
	CheckExpiry(cert *x509.Certificate) (*time.Duration, error)
	
	// VerifyHostname checks if the certificate is valid for the given hostname
	VerifyHostname(cert *x509.Certificate, hostname string) error
	
	// GenerateDiagnostics creates a detailed diagnostic report for the certificate
	GenerateDiagnostics(cert *x509.Certificate) (*CertDiagnostics, error)
}

// CertDiagnostics contains detailed information about certificate validation
type CertDiagnostics struct {
	Subject         string
	Issuer          string
	SerialNumber    string
	NotBefore       time.Time
	NotAfter        time.Time
	DNSNames        []string
	IPAddresses     []net.IP
	ValidationErrors []ValidationError
	Warnings        []string
}

// ValidationError represents a specific validation failure
type ValidationError struct {
	Type    string // "chain", "expiry", "hostname", etc.
	Message string
	Detail  string
}