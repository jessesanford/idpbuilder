// pkg/certs/types.go
package certs

import (
	"context"
	"crypto/x509"
	"net"
	"time"
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