// Package certs provides a clean abstraction layer for TLS certificate handling and validation,
// wrapping the standard library's crypto/x509 package behind testable interfaces.
package certs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
)

// Manager handles TLS certificate operations and provides a clean abstraction
// over the standard library's crypto/x509 package for certificate management.
// This interface enables mockability for testing and provides a foundation
// for secure registry connections.
type Manager interface {
	// LoadSystemCerts loads the system certificate pool.
	// Returns the system's root certificate pool or an error if loading fails.
	// The context can be used for cancellation during long-running operations.
	LoadSystemCerts(ctx context.Context) (*x509.CertPool, error)

	// ValidateCertificate validates a single certificate using standard x509 validation.
	// This method performs basic certificate validation including expiry checks
	// and signature verification against the system certificate pool.
	ValidateCertificate(ctx context.Context, cert *x509.Certificate) error

	// CreateTLSConfig creates a TLS configuration suitable for client connections.
	// The insecure parameter controls whether certificate validation should be skipped,
	// which is useful for development environments. When insecure is false,
	// the returned config will use the system certificate pool for validation.
	CreateTLSConfig(ctx context.Context, insecure bool) (*tls.Config, error)
}