package oci

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransport_InsecureModeSkipsTLSVerify verifies that when insecure mode is enabled,
// the transport configuration properly skips TLS certificate verification.
func TestTransport_InsecureModeSkipsTLSVerify(t *testing.T) {
	// Arrange
	config := TransportConfig{
		Insecure: true,
	}

	// Act
	transport := NewTransport(config)

	// Assert
	require.NotNil(t, transport, "Transport should not be nil")

	// Verify that the transport is an *http.Transport and has InsecureSkipVerify set to true
	httpTransport, ok := transport.(*http.Transport)
	require.True(t, ok, "Transport should be an *http.Transport")
	require.NotNil(t, httpTransport.TLSClientConfig, "TLS config should not be nil")
	assert.True(t, httpTransport.TLSClientConfig.InsecureSkipVerify, "InsecureSkipVerify should be true when insecure mode is enabled")
}

// TestTransport_SecureModeEnforcesTLS verifies that when secure mode is used (default),
// the transport configuration enforces strict TLS certificate validation.
func TestTransport_SecureModeEnforcesTLS(t *testing.T) {
	// Arrange
	config := TransportConfig{
		Insecure: false, // Explicit false, but this should be default
	}

	// Act
	transport := NewTransport(config)

	// Assert
	require.NotNil(t, transport, "Transport should not be nil")

	// Verify that the transport enforces TLS verification
	httpTransport, ok := transport.(*http.Transport)
	require.True(t, ok, "Transport should be an *http.Transport")

	// In secure mode, either TLSClientConfig is nil (uses defaults) or InsecureSkipVerify is false
	if httpTransport.TLSClientConfig != nil {
		assert.False(t, httpTransport.TLSClientConfig.InsecureSkipVerify, "InsecureSkipVerify should be false in secure mode")
	}
}

// TestTransport_ConfigurationOptions verifies that the transport configuration
// accepts and properly handles various TLS configuration options.
func TestTransport_ConfigurationOptions(t *testing.T) {
	testCases := []struct {
		name     string
		config   TransportConfig
		expected bool // Expected InsecureSkipVerify value
	}{
		{
			name:     "Default configuration (secure)",
			config:   TransportConfig{},
			expected: false,
		},
		{
			name:     "Explicit secure mode",
			config:   TransportConfig{Insecure: false},
			expected: false,
		},
		{
			name:     "Explicit insecure mode",
			config:   TransportConfig{Insecure: true},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			transport := NewTransport(tc.config)

			// Assert
			require.NotNil(t, transport, "Transport should not be nil")

			httpTransport, ok := transport.(*http.Transport)
			require.True(t, ok, "Transport should be an *http.Transport")

			if tc.expected {
				// Insecure mode - should have TLS config with skip verify
				require.NotNil(t, httpTransport.TLSClientConfig, "TLS config should not be nil in insecure mode")
				assert.True(t, httpTransport.TLSClientConfig.InsecureSkipVerify, "InsecureSkipVerify should match expected value")
			} else {
				// Secure mode - may have nil TLS config (defaults) or explicit false
				if httpTransport.TLSClientConfig != nil {
					assert.False(t, httpTransport.TLSClientConfig.InsecureSkipVerify, "InsecureSkipVerify should be false in secure mode")
				}
			}
		})
	}
}

// TestTransport_ConnectionWithSelfSignedCert verifies that the transport handles
// self-signed certificates correctly based on the insecure mode setting.
// This is a more integration-focused test.
func TestTransport_ConnectionWithSelfSignedCert(t *testing.T) {
	// Arrange - Create a self-signed certificate scenario
	// Note: This test focuses on transport configuration rather than actual connections

	// Test insecure mode with self-signed cert scenario
	t.Run("insecure mode accepts self-signed", func(t *testing.T) {
		config := TransportConfig{
			Insecure: true,
		}

		transport := NewTransport(config)
		require.NotNil(t, transport)

		httpTransport := transport.(*http.Transport)
		require.NotNil(t, httpTransport.TLSClientConfig)
		assert.True(t, httpTransport.TLSClientConfig.InsecureSkipVerify,
			"Insecure mode should skip certificate verification for self-signed certs")
	})

	// Test secure mode with self-signed cert scenario
	t.Run("secure mode rejects self-signed", func(t *testing.T) {
		config := TransportConfig{
			Insecure: false,
		}

		transport := NewTransport(config)
		require.NotNil(t, transport)

		httpTransport := transport.(*http.Transport)
		// Secure mode should either have nil TLS config or explicit verification enabled
		if httpTransport.TLSClientConfig != nil {
			assert.False(t, httpTransport.TLSClientConfig.InsecureSkipVerify,
				"Secure mode should enforce certificate verification")
		}
	})
}

// TestTransport_ErrorOnInvalidCertInSecureMode verifies that in secure mode,
// the transport properly validates certificates and would reject invalid ones.
func TestTransport_ErrorOnInvalidCertInSecureMode(t *testing.T) {
	// Arrange
	config := TransportConfig{
		Insecure: false,
	}

	// Act
	transport := NewTransport(config)

	// Assert
	require.NotNil(t, transport, "Transport should not be nil")

	httpTransport, ok := transport.(*http.Transport)
	require.True(t, ok, "Transport should be an *http.Transport")

	// Verify that certificate validation is enabled (InsecureSkipVerify is false or nil config)
	if httpTransport.TLSClientConfig != nil {
		assert.False(t, httpTransport.TLSClientConfig.InsecureSkipVerify,
			"Secure mode should validate certificates and reject invalid ones")

		// Verify that we can set custom certificate validation if needed
		assert.Nil(t, httpTransport.TLSClientConfig.VerifyPeerCertificate,
			"Should use default certificate verification unless custom verification is needed")
	}
}

// TestConfigureInsecure verifies that the ConfigureInsecure function properly
// modifies an existing transport's TLS configuration.
func TestConfigureInsecure(t *testing.T) {
	testCases := []struct {
		name     string
		insecure bool
		expected bool
	}{
		{
			name:     "Enable insecure mode",
			insecure: true,
			expected: true,
		},
		{
			name:     "Disable insecure mode",
			insecure: false,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{},
			}

			// Act
			ConfigureInsecure(transport, tc.insecure)

			// Assert
			require.NotNil(t, transport.TLSClientConfig, "TLS config should not be nil after configuration")
			assert.Equal(t, tc.expected, transport.TLSClientConfig.InsecureSkipVerify,
				"InsecureSkipVerify should match the configured value")
		})
	}
}

// TestConfigureInsecure_NilTLSConfig verifies that ConfigureInsecure handles
// the case where the transport has no existing TLS configuration.
func TestConfigureInsecure_NilTLSConfig(t *testing.T) {
	// Arrange
	transport := &http.Transport{
		TLSClientConfig: nil, // No existing TLS config
	}

	// Act
	ConfigureInsecure(transport, true)

	// Assert
	require.NotNil(t, transport.TLSClientConfig, "TLS config should be created if it was nil")
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify,
		"InsecureSkipVerify should be set even when starting with nil TLS config")
}

// TestTransport_SecurtiyWarnings verifies that insecure mode usage generates
// appropriate warnings or logging (this will be implemented in the GREEN phase).
func TestTransport_SecurityWarnings(t *testing.T) {
	// This test will verify that using insecure mode generates appropriate warnings
	// Implementation will be added during the GREEN phase

	// Arrange
	config := TransportConfig{
		Insecure: true,
	}

	// Act
	transport := NewTransport(config)

	// Assert
	require.NotNil(t, transport, "Transport should be created even in insecure mode")

	// TODO: In GREEN phase, add verification that warnings are logged
	// when insecure mode is used
}