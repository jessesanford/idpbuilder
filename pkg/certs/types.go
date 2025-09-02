// Package certs provides consolidated type definitions for certificate management
package certs

import (
	"context"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// KindCertExtractor defines the interface for extracting certificates from Kind clusters
type KindCertExtractor interface {
	// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
	ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)

	// GetCertInfo retrieves metadata about an extracted certificate
	GetCertInfo(cert *x509.Certificate) *CertificateInfo

	// ValidateCertificate performs basic validation on the extracted certificate
	ValidateCertificate(cert *x509.Certificate) error

	// SaveCertificate saves the certificate to the local trust store
	SaveCertificate(cert *x509.Certificate, path string) error
}

// CertificateInfo contains metadata about an extracted certificate
// This type provides compatibility across all certificate operations in Phase 1
type CertificateInfo struct {
	Subject   string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
	IsCA      bool
	DNSNames  []string
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

// ConnectionInfo holds TLS connection information for debugging and diagnostics
type ConnectionInfo struct {
	Registry          string
	IsSecure          bool
	IsInsecure        bool
	TLSVersion        string
	CipherSuite       string
	ServerCerts       []*x509.Certificate
	VerifiedChains    [][]*x509.Certificate
	HandshakeComplete bool
	Error             string
}

// TrustStoreUtils provides utility functions for trust store operations
type TrustStoreUtils struct{}

// SecurityLevel represents different security validation levels
type SecurityLevel int

const (
	SecurityNone SecurityLevel = iota
	SecurityLow
	SecurityMedium  
	SecurityHigh
)

// ValidationInput contains input parameters for certificate validation
type ValidationInput struct {
	Registry  string
	Operation string
	Error     error
}

// ValidationResult contains the results of a validation attempt
type ValidationResult struct {
	Success       bool
	Strategy      string
	Message       string
	SecurityLevel SecurityLevel
	Actions       []string
	NewConfig     map[string]interface{}
}

// FallbackStrategy defines the interface for fallback validation strategies
type FallbackStrategy interface {
	Name() string
	Priority() int
	CanHandle(err error) bool
	Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
}