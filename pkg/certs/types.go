// Package certs provides consolidated type definitions for certificate management
package certs

import (
	"crypto/x509"
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// CertificateInfo contains metadata about an extracted certificate
// This type provides compatibility with E1.1.1 (kind-certificate-extraction)
// and consolidates certificate information across all certificate operations
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
	
	// MaxConnsPerHost limits connections per host
	MaxConnsPerHost int
	
	// TLSHandshakeTimeout specifies TLS handshake timeout
	TLSHandshakeTimeout time.Duration
	
	// InsecureSkipVerify allows skipping TLS verification (development only)
	InsecureSkipVerify bool
	
	// CustomCA allows specifying custom CA certificates
	CustomCA *x509.CertPool
	
	// DebugTLS enables detailed TLS debugging
	DebugTLS bool
}

// TrustStoreManager provides an interface for managing certificate trust stores
// This interface is implemented by E1.1.2 (registry-tls-trust-integration)
type TrustStoreManager interface {
	// AddCertificate adds a certificate to the trust store
	AddCertificate(cert *x509.Certificate) error
	
	// RemoveCertificate removes a certificate from the trust store
	RemoveCertificate(cert *x509.Certificate) error
	
	// GetCertPool returns the current certificate pool
	GetCertPool() *x509.CertPool
	
	// SaveTrustStore persists the trust store to disk
	SaveTrustStore() error
	
	// LoadTrustStore loads the trust store from disk
	LoadTrustStore() error
	
	// CreateTransport creates an HTTP transport with the trust store
	CreateTransport(config *TransportConfig) (http.RoundTripper, error)
	
	// CreateRemoteOptions creates GGCR remote options with the trust store
	CreateRemoteOptions(config *TransportConfig) []remote.Option
	
	// ValidateCertificate validates a certificate against the trust store
	ValidateCertificate(cert *x509.Certificate) error
	
	// ListCertificates returns all certificates in the trust store
	ListCertificates() ([]*x509.Certificate, error)
}

// TLSDebugInfo contains detailed TLS handshake information for debugging
type TLSDebugInfo struct {
	Version           uint16
	CipherSuite       string
	ServerCerts       []*x509.Certificate
	VerifiedChains    [][]*x509.Certificate
	HandshakeComplete bool
	Error             string
}

// TrustStoreUtils provides utility functions for trust store operations
type TrustStoreUtils struct{}

// CertValidator provides comprehensive X.509 certificate validation
// Added from fallback-strategies effort
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
// Added from fallback-strategies effort
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
// Added from fallback-strategies effort
type ValidationError struct {
	Type    string // "chain", "expiry", "hostname", etc.
	Message string
	Detail  string
}
