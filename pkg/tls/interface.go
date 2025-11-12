// Package tls provides TLS configuration for secure registry connections.
package tls

import (
	"crypto/tls"
)

// TLSProvider generates TLS configurations for HTTPS connections to OCI registries.
type TLSProvider interface {
	// GetTLSConfig returns a TLS configuration for registry HTTPS connections.
	//
	// In secure mode (default):
	//   - Uses the system certificate pool
	//   - Verifies certificate chains
	//   - Checks that hostnames match certificates
	//
	// In insecure mode:
	//   - Sets InsecureSkipVerify = true
	//   - Accepts self-signed certificates
	//   - Should only be used for local development
	//
	// Example:
	//   tlsConfig := tlsProvider.GetTLSConfig()
	//   transport := &http.Transport{TLSClientConfig: tlsConfig}
	GetTLSConfig() *tls.Config

	// IsInsecure returns true if certificate verification is disabled.
	//
	// Example:
	//   if tlsProvider.IsInsecure() {
	//       fmt.Println("WARNING: Certificate verification disabled")
	//   }
	IsInsecure() bool

	// GetWarningMessage returns a user-facing warning message for insecure mode.
	// Returns an empty string if in secure mode.
	//
	// Example:
	//   if warning := tlsProvider.GetWarningMessage(); warning != "" {
	//       fmt.Println(warning)
	//   }
	GetWarningMessage() string
}

// NewTLSProvider creates a new TLS configuration provider.
//
// If insecure is true:
//   - Certificate verification is disabled (InsecureSkipVerify = true)
//   - Self-signed certificates are accepted
//   - A warning message is generated
//   - Use only for local development with Gitea
//
// If insecure is false:
//   - Standard TLS verification is enabled
//   - System certificate pool is used
//   - Production-ready configuration
//
// Example (secure mode):
//   provider, err := tls.NewTLSProvider(false)
//
// Example (insecure mode for local Gitea):
//   provider, err := tls.NewTLSProvider(true)
//   if warning := provider.GetWarningMessage(); warning != "" {
//       log.Println(warning)
//   }
func NewTLSProvider(insecure bool) (TLSProvider, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}
