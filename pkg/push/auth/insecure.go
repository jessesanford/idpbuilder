package auth

import (
	"crypto/tls"
	"net/http"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GetInsecureTransport returns an HTTP transport that skips TLS verification
// This should only be used with the --insecure flag for self-signed certificates
func GetInsecureTransport() *http.Transport {
	// Create a copy of the default transport
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// Configure TLS to skip certificate verification
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return transport
}

// GetInsecureOption returns a remote.Option for insecure registries
// This configures go-containerregistry to skip TLS certificate verification
func GetInsecureOption() remote.Option {
	return remote.WithTransport(GetInsecureTransport())
}

// GetSecureTransport returns the default HTTP transport with TLS verification enabled
// This is the default behavior and is included for completeness
func GetSecureTransport() *http.Transport {
	// Return a copy of the default transport which has TLS verification enabled
	return http.DefaultTransport.(*http.Transport).Clone()
}

// GetSecureOption returns a remote.Option for secure registries
// This explicitly configures go-containerregistry to use TLS certificate verification
func GetSecureOption() remote.Option {
	return remote.WithTransport(GetSecureTransport())
}

// NewTransportOption creates a remote.Option with the appropriate transport based on insecure flag
func NewTransportOption(insecure bool) remote.Option {
	if insecure {
		return GetInsecureOption()
	}
	return GetSecureOption()
}

// IsInsecureRegistry determines if a registry URL suggests insecure communication
// This is a helper function to warn users about potential insecure connections
func IsInsecureRegistry(registryURL string) bool {
	// Check for common indicators of insecure registries
	insecureIndicators := []string{
		"http://",   // Plain HTTP
		"localhost", // Local development
		"127.0.0.1", // Loopback
		"::1",       // IPv6 loopback
		".local",    // mDNS domains
	}

	for _, indicator := range insecureIndicators {
		if len(registryURL) >= len(indicator) &&
			registryURL[:len(indicator)] == indicator ||
			len(registryURL) > len(indicator) && registryURL[len(registryURL)-len(indicator):] == indicator {
			return true
		}
	}

	return false
}
