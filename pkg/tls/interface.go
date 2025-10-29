// Package tls provides interfaces and types for TLS configuration.
package tls

import (
	"crypto/tls"
)

// ConfigProvider defines operations for providing TLS configuration.
type ConfigProvider interface {
	// GetTLSConfig returns a tls.Config for HTTP transport.
	//
	// Returns:
	//   - *tls.Config: TLS configuration for HTTP transport
	GetTLSConfig() *tls.Config

	// IsInsecure returns whether insecure mode is enabled.
	//
	// Returns:
	//   - bool: true if --insecure flag was set, false otherwise
	IsInsecure() bool
}

// Config holds TLS configuration options.
type Config struct {
	// InsecureSkipVerify controls whether to skip TLS certificate verification.
	//
	// When true: Certificate validity NOT checked (development only)
	// When false: Full certificate verification (production)
	InsecureSkipVerify bool
}

// NewConfigProvider creates a TLS configuration provider.
//
// Parameters:
//   - insecure: Whether to enable insecure mode (skip cert verification)
//
// Returns:
//   - ConfigProvider: TLS configuration provider interface implementation
func NewConfigProvider(insecure bool) ConfigProvider {
	// Implementation will be provided in Wave 2 (pkg/tls/config.go)
	panic("not implemented - interface definition only")
}
