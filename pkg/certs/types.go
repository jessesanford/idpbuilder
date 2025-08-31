// pkg/certs/types.go
package certs

import (
	"context"
	"crypto/x509"
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