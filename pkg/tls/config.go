// Package tls provides TLS configuration implementations.
package tls

import (
	"crypto/tls"
	"crypto/x509"
)

// tlsConfigProvider implements the ConfigProvider interface.
// It provides TLS configuration for HTTP transport with support for
// both secure and insecure (certificate verification skipped) modes.
type tlsConfigProvider struct {
	config Config // Configuration from Wave 1 interface
}

// NewConfigProvider creates a TLS configuration provider.
//
// The provider can operate in two modes:
//   - Secure mode (insecure=false, default): Full certificate verification using
//     system certificate pool. This is the recommended mode for production use.
//   - Insecure mode (insecure=true, --insecure flag): Certificate verification is
//     disabled. This mode should ONLY be used for development/testing with
//     self-signed certificates.
//
// WARNING: Insecure mode disables TLS certificate verification, making connections
// vulnerable to man-in-the-middle attacks. Never use in production without
// understanding the security implications.
//
// Parameters:
//   - insecure: Whether to enable insecure mode (skip cert verification).
//     Typically set from the --insecure / -k CLI flag.
//
// Returns:
//   - ConfigProvider: TLS configuration provider interface implementation
//
// Example (secure mode - recommended):
//
//	provider := tls.NewConfigProvider(false)
//	tlsConfig := provider.GetTLSConfig()
//	// tlsConfig has full certificate verification enabled
//
// Example (insecure mode - development only):
//
//	provider := tls.NewConfigProvider(true)
//	if provider.IsInsecure() {
//	    log.Warn("TLS certificate verification disabled (insecure mode)")
//	}
//	tlsConfig := provider.GetTLSConfig()
//	// tlsConfig.InsecureSkipVerify == true
func NewConfigProvider(insecure bool) ConfigProvider {
	return &tlsConfigProvider{
		config: Config{
			InsecureSkipVerify: insecure,
		},
	}
}

// GetTLSConfig returns a tls.Config for HTTP transport.
//
// The returned configuration depends on the mode set during provider creation:
//
// Secure mode (default, InsecureSkipVerify=false):
//   - Loads system certificate pool for verification
//   - Falls back to empty pool if system certs unavailable
//   - Full certificate chain verification enabled
//   - Suitable for production use
//
// Insecure mode (InsecureSkipVerify=true):
//   - Certificate verification disabled
//   - No certificate pool needed
//   - Accepts any certificate (including self-signed, expired, invalid)
//   - Should only be used for development/testing
//
// The returned tls.Config is ready to be used with http.Transport for
// establishing HTTPS connections to container registries.
//
// Returns:
//   - *tls.Config: TLS configuration for HTTP transport. Never returns nil.
//
// Example usage with HTTP client:
//
//	provider := tls.NewConfigProvider(false)
//	tlsConfig := provider.GetTLSConfig()
//	transport := &http.Transport{
//	    TLSClientConfig: tlsConfig,
//	}
//	client := &http.Client{Transport: transport}
//	// client now uses proper TLS configuration
func (p *tlsConfigProvider) GetTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: p.config.InsecureSkipVerify,
	}

	// If secure mode, load system certificate pool for verification
	if !p.config.InsecureSkipVerify {
		// Attempt to load system certificates
		certPool, err := x509.SystemCertPool()
		if err != nil {
			// Fallback to empty pool if system certs unavailable
			// This allows operation even if system cert pool can't be loaded,
			// but certificate verification will use the empty pool
			certPool = x509.NewCertPool()
		}
		tlsConfig.RootCAs = certPool
	}

	return tlsConfig
}

// IsInsecure returns whether insecure mode is enabled.
//
// This method allows callers to check if the TLS provider is operating in
// insecure mode (certificate verification disabled). This is useful for:
//   - Logging warnings when insecure mode is active
//   - Conditional behavior based on security mode
//   - Validation that security requirements are met
//
// Returns:
//   - bool: true if --insecure flag was set (certificate verification disabled),
//     false if in secure mode (certificate verification enabled)
//
// Example usage for warning:
//
//	provider := tls.NewConfigProvider(insecureFlag)
//	if provider.IsInsecure() {
//	    log.Warn("WARNING: TLS certificate verification disabled (insecure mode)")
//	    log.Warn("This makes connections vulnerable to man-in-the-middle attacks")
//	}
//
// Example usage for validation:
//
//	provider := tls.NewConfigProvider(insecureFlag)
//	if provider.IsInsecure() && productionMode {
//	    return fmt.Errorf("insecure mode not allowed in production")
//	}
func (p *tlsConfigProvider) IsInsecure() bool {
	return p.config.InsecureSkipVerify
}
