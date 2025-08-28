package certs

import (
	"crypto/x509"
	"context"
)

// TrustManager defines the interface for managing trust store certificates
type TrustManager interface {
	// AddCertificate adds a certificate to the trust store for a specific registry
	AddCertificate(ctx context.Context, registry string, cert *Certificate) error

	// RemoveCertificate removes a certificate from the trust store for a specific registry
	RemoveCertificate(ctx context.Context, registry string, fingerprint string) error

	// ListCertificates lists all certificates for a specific registry
	ListCertificates(ctx context.Context, registry string) ([]Certificate, error)

	// GetRegistryConfig gets the complete configuration for a registry
	GetRegistryConfig(ctx context.Context, registry string) (*RegistryConfig, error)

	// SetInsecureRegistry configures a registry to skip TLS verification
	SetInsecureRegistry(ctx context.Context, registry string, insecure bool) error

	// ValidateCertificate validates a certificate against the trust store
	ValidateCertificate(ctx context.Context, registry string, cert *x509.Certificate) error
}

// CertificateStore defines the interface for certificate storage operations
type CertificateStore interface {
	// Store writes a certificate to the filesystem
	Store(registry string, cert *Certificate) error

	// Load reads a certificate from the filesystem
	Load(registry string, fingerprint string) (*Certificate, error)

	// Delete removes a certificate from the filesystem
	Delete(registry string, fingerprint string) error

	// List returns all certificates for a registry
	List(registry string) ([]Certificate, error)

	// Exists checks if a certificate exists in the store
	Exists(registry string, fingerprint string) (bool, error)
}

// RegistryConfigManager defines the interface for managing registry configurations
type RegistryConfigManager interface {
	// UpdateInsecureRegistry updates the insecure registry configuration
	UpdateInsecureRegistry(registry string, insecure bool) error

	// GetInsecureRegistries returns a list of registries configured as insecure
	GetInsecureRegistries() ([]string, error)

	// LoadConfig loads the registry configuration from disk
	LoadConfig() error

	// SaveConfig saves the registry configuration to disk
	SaveConfig() error
}