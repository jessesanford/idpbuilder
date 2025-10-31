package tls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewConfigProvider_SecureMode verifies that NewConfigProvider creates
// a provider in secure mode when insecure=false.
// Test Case: TC-TLS-IMPL-001
func TestNewConfigProvider_SecureMode(t *testing.T) {
	// Given: Secure mode (false = verify certificates)
	insecure := false

	// When: Creating provider
	provider := NewConfigProvider(insecure)

	// Then: Provider created in secure mode
	require.NotNil(t, provider, "Provider should not be nil")

	// Verify it implements the interface
	var _ ConfigProvider = provider

	// Verify it's in secure mode
	assert.False(t, provider.IsInsecure(), "Provider should be in secure mode")
}

// TestNewConfigProvider_InsecureMode verifies that NewConfigProvider creates
// a provider in insecure mode when insecure=true.
// Test Case: TC-TLS-IMPL-002
func TestNewConfigProvider_InsecureMode(t *testing.T) {
	// Given: Insecure mode (true = skip verification)
	insecure := true

	// When: Creating provider
	provider := NewConfigProvider(insecure)

	// Then: Provider created in insecure mode
	require.NotNil(t, provider, "Provider should not be nil")

	// Verify it implements the interface
	var _ ConfigProvider = provider

	// Verify it's in insecure mode
	assert.True(t, provider.IsInsecure(), "Provider should be in insecure mode")
}

// TestGetTLSConfig_SecureMode verifies that GetTLSConfig returns proper
// configuration for secure mode with system certificate pool.
// Test Case: TC-TLS-IMPL-003
func TestGetTLSConfig_SecureMode(t *testing.T) {
	// Given: Provider in secure mode
	provider := NewConfigProvider(false)

	// When: Getting TLS config
	tlsConfig := provider.GetTLSConfig()

	// Then: Config has proper settings
	require.NotNil(t, tlsConfig, "TLS config should not be nil")

	// Verify InsecureSkipVerify is false (certificate verification enabled)
	assert.False(t, tlsConfig.InsecureSkipVerify,
		"InsecureSkipVerify should be false in secure mode")

	// Verify system cert pool is loaded
	assert.NotNil(t, tlsConfig.RootCAs,
		"RootCAs should be set with system certificate pool")
}

// TestGetTLSConfig_InsecureMode verifies that GetTLSConfig returns proper
// configuration for insecure mode with certificate verification disabled.
// Test Case: TC-TLS-IMPL-004
func TestGetTLSConfig_InsecureMode(t *testing.T) {
	// Given: Provider in insecure mode
	provider := NewConfigProvider(true)

	// When: Getting TLS config
	tlsConfig := provider.GetTLSConfig()

	// Then: InsecureSkipVerify is true
	require.NotNil(t, tlsConfig, "TLS config should not be nil")

	// Verify InsecureSkipVerify is true (certificate verification disabled)
	assert.True(t, tlsConfig.InsecureSkipVerify,
		"InsecureSkipVerify should be true in insecure mode")

	// Note: RootCAs is not required in insecure mode
	// (verification is disabled, so cert pool is not used)
}

// TestIsInsecure_Secure verifies that IsInsecure returns false for
// secure mode providers.
// Test Case: TC-TLS-IMPL-005
func TestIsInsecure_Secure(t *testing.T) {
	// Given: Provider in secure mode
	provider := NewConfigProvider(false)

	// When/Then: IsInsecure returns false
	assert.False(t, provider.IsInsecure(),
		"IsInsecure should return false for secure mode")
}

// TestIsInsecure_Insecure verifies that IsInsecure returns true for
// insecure mode providers.
// Test Case: TC-TLS-IMPL-006
func TestIsInsecure_Insecure(t *testing.T) {
	// Given: Provider in insecure mode
	provider := NewConfigProvider(true)

	// When/Then: IsInsecure returns true
	assert.True(t, provider.IsInsecure(),
		"IsInsecure should return true for insecure mode")
}

// TestTLSConfig_UsableWithHTTPClient verifies that the TLS configuration
// returned by GetTLSConfig can be used with http.Client and http.Transport.
// Test Case: TC-TLS-IMPL-007
func TestTLSConfig_UsableWithHTTPClient(t *testing.T) {
	// Given: TLS provider (using insecure mode for simplicity)
	provider := NewConfigProvider(true)
	tlsConfig := provider.GetTLSConfig()

	// When: Creating HTTP transport with TLS config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}

	// Then: Client created successfully
	require.NotNil(t, client, "HTTP client should be created")
	require.NotNil(t, client.Transport, "HTTP client should have transport")

	// Verify the TLS config is properly set
	httpTransport, ok := client.Transport.(*http.Transport)
	require.True(t, ok, "Transport should be *http.Transport")
	assert.Equal(t, tlsConfig, httpTransport.TLSClientConfig,
		"TLS config should be set on transport")
}

// TestGetTLSConfig_MultipleCallsReturnNewInstances verifies that calling
// GetTLSConfig multiple times returns different instances (not cached).
// Test Case: TC-TLS-IMPL-008
func TestGetTLSConfig_MultipleCallsReturnNewInstances(t *testing.T) {
	// Given: TLS provider
	provider := NewConfigProvider(false)

	// When: Calling GetTLSConfig multiple times
	config1 := provider.GetTLSConfig()
	config2 := provider.GetTLSConfig()

	// Then: Both configs are valid but different instances
	require.NotNil(t, config1, "First config should not be nil")
	require.NotNil(t, config2, "Second config should not be nil")

	// Note: We expect different instances (not cached)
	// This is fine - configs are lightweight and stateless
	assert.NotSame(t, config1, config2,
		"GetTLSConfig should return new instances")

	// But they should have same configuration values
	assert.Equal(t, config1.InsecureSkipVerify, config2.InsecureSkipVerify,
		"Configs should have same InsecureSkipVerify setting")
}

// TestTLSConfig_SecureMode_CompatibleWithHTTPSConnections verifies that
// the secure mode TLS config is properly configured for HTTPS connections.
// Test Case: TC-TLS-IMPL-009
func TestTLSConfig_SecureMode_CompatibleWithHTTPSConnections(t *testing.T) {
	// Given: Provider in secure mode
	provider := NewConfigProvider(false)
	tlsConfig := provider.GetTLSConfig()

	// Then: TLS config should be suitable for HTTPS connections
	require.NotNil(t, tlsConfig, "TLS config should not be nil")

	// Verify certificate verification is enabled
	assert.False(t, tlsConfig.InsecureSkipVerify,
		"Certificate verification should be enabled")

	// Verify root CA pool is set (for certificate chain validation)
	assert.NotNil(t, tlsConfig.RootCAs,
		"Root CA pool should be set for certificate validation")

	// Create HTTP client to verify config is usable
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}

	require.NotNil(t, client, "HTTP client should be created successfully")
}

// TestTLSConfig_InsecureMode_AcceptsAnyCertificate verifies that
// insecure mode disables certificate verification as expected.
// Test Case: TC-TLS-IMPL-010
func TestTLSConfig_InsecureMode_AcceptsAnyCertificate(t *testing.T) {
	// Given: Provider in insecure mode
	provider := NewConfigProvider(true)
	tlsConfig := provider.GetTLSConfig()

	// Then: Certificate verification should be disabled
	require.NotNil(t, tlsConfig, "TLS config should not be nil")

	// Verify InsecureSkipVerify is true (will accept any certificate)
	assert.True(t, tlsConfig.InsecureSkipVerify,
		"InsecureSkipVerify should be true - will accept any certificate")

	// This means the config will accept:
	// - Self-signed certificates
	// - Expired certificates
	// - Invalid certificates
	// - Certificates with wrong hostname
	// (This is intentional behavior for --insecure flag)
}
