package certificate

import (
	"crypto/x509"
	"time"
)

// CertificateManager defines the interface for certificate management operations.
// It provides methods to generate, store, retrieve, and validate certificates
// for various use cases within the IDPBuilder system.
type CertificateManager interface {
	// GenerateSelfSigned creates a new self-signed certificate with the specified options.
	// Returns the certificate and private key in PEM format.
	GenerateSelfSigned(opts *GenerationOptions) (*Certificate, error)

	// Store saves a certificate and its private key for later retrieval.
	// The key parameter is used to identify the certificate for future operations.
	Store(key string, cert *Certificate) error

	// Retrieve gets a previously stored certificate by its key.
	// Returns an error if the certificate is not found.
	Retrieve(key string) (*Certificate, error)

	// IsValid checks if a certificate is currently valid (not expired, proper signature).
	IsValid(cert *Certificate) bool

	// GetExpiration returns the expiration time of a certificate.
	GetExpiration(cert *Certificate) (time.Time, error)

	// List returns all stored certificate keys.
	List() ([]string, error)

	// Delete removes a stored certificate.
	Delete(key string) error
}

// Certificate represents a certificate and its associated private key.
type Certificate struct {
	// CertPEM contains the certificate in PEM format
	CertPEM []byte
	// KeyPEM contains the private key in PEM format
	KeyPEM []byte
	// Metadata contains additional information about the certificate
	Metadata *CertificateMetadata
}

// CertificateMetadata contains additional information about a certificate.
type CertificateMetadata struct {
	// Subject contains the certificate subject information
	Subject string
	// DNSNames contains the DNS names (SANs) included in the certificate
	DNSNames []string
	// NotBefore is the certificate validity start time
	NotBefore time.Time
	// NotAfter is the certificate validity end time
	NotAfter time.Time
	// IsCA indicates if this is a Certificate Authority certificate
	IsCA bool
	// KeyUsage contains the certificate key usage flags
	KeyUsage x509.KeyUsage
	// ExtKeyUsage contains the extended key usage flags
	ExtKeyUsage []x509.ExtKeyUsage
}

// GenerationOptions contains configuration for certificate generation.
type GenerationOptions struct {
	// Subject contains the certificate subject information
	Subject string
	// Organization contains the organization name
	Organization string
	// DNSNames contains the DNS names to include as Subject Alternative Names
	DNSNames []string
	// ValidFor specifies how long the certificate should be valid
	ValidFor time.Duration
	// IsCA indicates whether this should be a CA certificate
	IsCA bool
	// KeyUsage specifies the key usage for the certificate
	KeyUsage x509.KeyUsage
	// ExtKeyUsage specifies the extended key usage for the certificate
	ExtKeyUsage []x509.ExtKeyUsage
}

// DefaultGenerationOptions returns a GenerationOptions struct with sensible defaults
// for creating self-signed certificates suitable for TLS server authentication.
func DefaultGenerationOptions() *GenerationOptions {
	return &GenerationOptions{
		Organization: "IDPBuilder",
		ValidFor:     time.Hour * 24 * 365, // 1 year
		IsCA:         false,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
}