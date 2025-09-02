// Package certs provides consolidated type definitions for certificate management and validation
// This file consolidates types from all Phase 1 efforts to ensure compatibility
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
	SecurityHigh SecurityLevel = iota
	SecurityMedium
	SecurityLow
	SecurityNone
)

// ValidationInput contains the input data for certificate validation operations
// This consolidates both basic validation and advanced certificate validation
type ValidationInput struct {
	Certificates []*x509.Certificate   // New field from certificate-validation-pipeline
	Registry     string
	Operation    string
	Error        error
	Options      map[string]interface{} // New field from certificate-validation-pipeline
}

// ValidationResult contains the result of a certificate validation attempt
type ValidationResult struct {
	Success       bool
	Strategy      string
	Message       string
	SecurityLevel SecurityLevel
	Actions       []string
	NewConfig     map[string]interface{}
}

// RecoveryConfig contains configuration for certificate recovery operations
// This consolidates fields from both fallback.go and recovery.go versions
type RecoveryConfig struct {
	// From recovery.go
	EnableCertRefresh   bool
	EnableTrustUpdate   bool
	EnableChainRepair   bool
	MaxAttempts         int
	Timeout             time.Duration
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration
	
	// From fallback.go (additional fields)  
	AllowInsecure      bool
	FallbackRegistries []string
	RecoveryStrategies []string
	TrustStoreUpdate   bool
}

// RecoveryResult contains the result of a recovery operation
// This consolidates fields from both fallback.go and recovery.go versions
type RecoveryResult struct {
	Success          bool
	Method           string    // From recovery.go (was 'Method')
	Strategy         string    // From fallback.go (additional)
	Actions          []string
	NewConfig        interface{}
	Message          string    // From recovery.go
	FailureReason    string    // From fallback.go (additional)
	RecoveredCerts   []*x509.Certificate // From fallback.go (additional)
	SecurityDowngrade bool     // From fallback.go (additional)
	NextRetryAfter   time.Duration // From fallback.go (additional)
}

// FallbackStrategy defines the interface for certificate validation fallback strategies
type FallbackStrategy interface {
	Name() string
	Priority() int
	CanHandle(err error) bool
	Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
}

// CertDiagnostics contains comprehensive certificate diagnostic information
type CertDiagnostics struct {
	Subject          string
	Issuer           string
	SerialNumber     string
	NotBefore        time.Time
	NotAfter         time.Time
	DNSNames         []string
	IPAddresses      []net.IP
	ValidationErrors []ValidationError
	Warnings         []string
}

// ValidationError represents a certificate validation error
type ValidationError struct {
	Type    string // e.g., "chain", "expiry", "hostname"
	Message string
	Detail  string
}