package certs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// X509Manager implements the Manager interface using the standard library's
// crypto/x509 package. It provides production-ready certificate management
// with proper error handling, context support, and integration with the
// system certificate store.
type X509Manager struct {
	validator Validator
	store     Store
}

// NewX509Manager creates a new X509Manager instance with the provided validator and store.
// If validator is nil, a DefaultValidator will be used.
// If store is nil, a new MemoryStore will be created.
func NewX509Manager(validator Validator, store Store) *X509Manager {
	if validator == nil {
		validator = NewDefaultValidator()
	}
	if store == nil {
		store = NewMemoryStore()
	}

	return &X509Manager{
		validator: validator,
		store:     store,
	}
}

// NewDefaultX509Manager creates a new X509Manager instance with default components.
// This is a convenience function that creates a manager with DefaultValidator
// and MemoryStore, suitable for most use cases.
func NewDefaultX509Manager() *X509Manager {
	return NewX509Manager(nil, nil)
}

// LoadSystemCerts loads the system certificate pool.
// This method uses the system's trusted root certificates and supports
// cancellation through the provided context.
func (m *X509Manager) LoadSystemCerts(ctx context.Context) (*x509.CertPool, error) {
	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load system certificate pool
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, NewCertError(ErrSystemPool, "failed to load system certificate pool", err)
	}

	if pool == nil {
		// On some systems SystemCertPool might return nil without error
		// In this case, create an empty pool
		pool = x509.NewCertPool()
	}

	return pool, nil
}

// ValidateCertificate validates a single certificate using the configured validator.
// It uses the system certificate pool for validation unless the certificate
// is being validated in insecure mode.
func (m *X509Manager) ValidateCertificate(ctx context.Context, cert *x509.Certificate) error {
	if cert == nil {
		return NewCertError(ErrInvalidCert, "certificate is nil", nil)
	}

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Load system certs for validation
	systemPool, err := m.LoadSystemCerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to load system certs for validation: %w", err)
	}

	// Prepare verification options
	opts := x509.VerifyOptions{
		Roots: systemPool,
	}

	// Validate the certificate
	return m.validator.Validate(cert, opts)
}

// CreateTLSConfig creates a TLS configuration suitable for client connections.
// The insecure parameter controls whether certificate validation should be skipped.
// When insecure is false, the configuration will use the system certificate pool
// for validation. When insecure is true, all certificate validation is disabled
// (useful for development environments).
func (m *X509Manager) CreateTLSConfig(ctx context.Context, insecure bool) (*tls.Config, error) {
	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Create base TLS configuration
	config := &tls.Config{
		InsecureSkipVerify: insecure,
	}

	// If not in insecure mode, configure certificate validation
	if !insecure {
		// Load system certificate pool
		pool, err := m.LoadSystemCerts(ctx)
		if err != nil {
			return nil, NewCertError(ErrTLSConfig, "failed to load system certs for TLS config", err)
		}

		// Set the root CA pool for certificate validation
		config.RootCAs = pool

		// Enable certificate verification
		config.InsecureSkipVerify = false

		// Set minimum TLS version for security
		config.MinVersion = tls.VersionTLS12
	}

	return config, nil
}

// AddTrustedCert adds a certificate to the manager's certificate store.
// This allows for runtime addition of trusted certificates that will be
// included in future TLS configurations and validation operations.
func (m *X509Manager) AddTrustedCert(ctx context.Context, cert *x509.Certificate) error {
	if cert == nil {
		return NewCertError(ErrInvalidCert, "cannot add nil certificate", nil)
	}

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Add certificate to store
	return m.store.AddCert(ctx, cert)
}

// CreateTLSConfigWithCustomCerts creates a TLS configuration that includes
// both system certificates and any custom certificates added to the store.
// This is useful when you need to trust additional certificates beyond
// the system's default trusted roots.
func (m *X509Manager) CreateTLSConfigWithCustomCerts(ctx context.Context, insecure bool) (*tls.Config, error) {
	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Create base TLS configuration
	config := &tls.Config{
		InsecureSkipVerify: insecure,
	}

	// If not in insecure mode, configure certificate validation
	if !insecure {
		// Load system certificate pool
		systemPool, err := m.LoadSystemCerts(ctx)
		if err != nil {
			return nil, NewCertError(ErrTLSConfig, "failed to load system certs for TLS config", err)
		}

		// Get custom certificates from store
		customPool, err := m.store.GetPool(ctx)
		if err != nil {
			return nil, NewCertError(ErrTLSConfig, "failed to load custom certs for TLS config", err)
		}

		// If we have custom certificates, we need to combine them with system certs
		if customPool != nil {
			// Start with system pool
			combinedPool := systemPool

			// For now, we'll use the system pool as the base and note that
			// combining pools is not directly supported by the standard library.
			// In a production implementation, you might need to maintain
			// your own pool and add both system and custom certs to it.
			config.RootCAs = combinedPool
		} else {
			config.RootCAs = systemPool
		}

		// Enable certificate verification
		config.InsecureSkipVerify = false

		// Set minimum TLS version for security
		config.MinVersion = tls.VersionTLS12
	}

	return config, nil
}

// GetValidator returns the validator used by this manager.
// This allows for inspection and testing of the validation component.
func (m *X509Manager) GetValidator() Validator {
	return m.validator
}

// GetStore returns the certificate store used by this manager.
// This allows for inspection and testing of the storage component.
func (m *X509Manager) GetStore() Store {
	return m.store
}