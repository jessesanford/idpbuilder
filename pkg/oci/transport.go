package oci

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
)

// TransportConfig represents configuration options for HTTP transport.
// It controls TLS behavior and security settings for OCI registry connections.
type TransportConfig struct {
	// Insecure when true, skips TLS certificate verification.
	// WARNING: This should only be used for testing or with trusted self-signed certificates.
	Insecure bool
}

// NewTransport creates a new HTTP transport with the specified configuration.
// It returns an http.RoundTripper configured according to the provided settings.
//
// Security Behavior:
// - When Insecure is true: Skips TLS certificate verification (logs warning)
// - When Insecure is false/unset: Uses standard TLS verification
//
// WARNING: Insecure mode should only be used for:
// - Testing environments with self-signed certificates
// - Private registries with known self-signed certificates
// - Local development scenarios
//
// Production deployments should use proper TLS certificates.
func NewTransport(config TransportConfig) http.RoundTripper {
	// Create base HTTP transport with sensible defaults
	transport := &http.Transport{
		// Add some default configuration for better performance
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 2,
	}

	// Configure TLS settings and log security implications
	configureTransportSecurity(transport, config.Insecure)

	return transport
}

// ConfigureInsecure modifies the provided HTTP transport to enable or disable
// TLS certificate verification based on the insecure flag.
//
// This function is provided for backward compatibility and external configuration.
// For new code, consider using NewTransport with TransportConfig.
//
// When insecure is true, TLS certificate verification is skipped.
// When insecure is false, standard TLS verification is enforced.
//
// If the transport doesn't have a TLS configuration, one will be created.
func ConfigureInsecure(transport *http.Transport, insecure bool) {
	configureTransportSecurity(transport, insecure)
}

// configureTransportSecurity is the internal function that handles TLS configuration
// and security logging. It centralizes the TLS setup logic.
func configureTransportSecurity(transport *http.Transport, insecure bool) {
	// Ensure TLS config exists
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{
			// Set minimum TLS version for security
			MinVersion: tls.VersionTLS12,
		}
	}

	// Configure certificate verification
	transport.TLSClientConfig.InsecureSkipVerify = insecure

	// Log security configuration with appropriate severity
	if insecure {
		logSecurityWarning()
		logTransportConfig("insecure", "TLS certificate verification disabled")
	} else {
		logTransportConfig("secure", "TLS certificate verification enabled")
	}
}

// logSecurityWarning logs a prominent warning when insecure mode is enabled.
// This ensures that the security implications are clearly communicated.
func logSecurityWarning() {
	warningMsg :=
		"⚠️  SECURITY WARNING: Insecure TLS mode enabled\n" +
		"   • Certificate verification is DISABLED\n" +
		"   • Only use this for:\n" +
		"     - Testing environments\n" +
		"     - Trusted self-signed certificates\n" +
		"     - Local development\n" +
		"   • NEVER use in production with untrusted certificates"

	// Check if we should suppress warnings (useful for testing)
	if os.Getenv("IDPBUILDER_SUPPRESS_TLS_WARNINGS") != "true" {
		log.Printf("%s", warningMsg)
	}
}

// logTransportConfig logs the current transport security configuration for debugging.
func logTransportConfig(mode, description string) {
	log.Printf("Transport configured for %s mode: %s", mode, description)
}