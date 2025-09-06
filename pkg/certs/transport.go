// Package certs provides transport configuration for go-containerregistry
package certs

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// DefaultTransportConfig returns a default transport configuration
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		Timeout:             30 * time.Second,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     90 * time.Second,
	}
}

// ConfigureTransport creates a remote.Option with proper TLS configuration for a registry
func (m *trustStoreManager) ConfigureTransport(registry string) (remote.Option, error) {
	return m.ConfigureTransportWithConfig(registry, DefaultTransportConfig())
}

// ConfigureTransportWithConfig creates a remote.Option with custom transport configuration
func (m *trustStoreManager) ConfigureTransportWithConfig(registry string, config *TransportConfig) (remote.Option, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry name cannot be empty")
	}
	
	if config == nil {
		config = DefaultTransportConfig()
	}

	// Check if registry is marked as insecure
	if m.IsInsecure(registry) {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		}

		return remote.WithTransport(transport), nil
	}

	// Get certificate pool for secure connection
	certPool, err := m.GetCertPool(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate pool for %s: %w\n"+
			"To fix this issue:\n"+
			"1. Ensure the certificate is valid and properly formatted\n"+
			"2. Check if the certificate has expired\n"+
			"3. Verify certificate file permissions\n"+
			"4. Or use --insecure flag for testing", registry, err)
	}

	// Create transport with custom CA pool
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12, // Enforce minimum TLS version
		},
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
	}

	return remote.WithTransport(transport), nil
}

// ConfigureAuthenticatedTransport creates a remote.Option with authentication and TLS
func (m *trustStoreManager) ConfigureAuthenticatedTransport(registry string, auth remote.Option) ([]remote.Option, error) {
	transportOption, err := m.ConfigureTransport(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to configure transport for %s: %w", registry, err)
	}

	return []remote.Option{transportOption, auth}, nil
}

// CreateHTTPClient creates an HTTP client with proper TLS configuration
func (m *trustStoreManager) CreateHTTPClient(registry string) (*http.Client, error) {
	return m.CreateHTTPClientWithConfig(registry, DefaultTransportConfig())
}

// CreateHTTPClientWithConfig creates an HTTP client with custom configuration
func (m *trustStoreManager) CreateHTTPClientWithConfig(registry string, config *TransportConfig) (*http.Client, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry name cannot be empty")
	}
	
	if config == nil {
		config = DefaultTransportConfig()
	}

	var transport *http.Transport

	if m.IsInsecure(registry) {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		}
	} else {
		// Get certificate pool for secure connection
		certPool, err := m.GetCertPool(registry)
		if err != nil {
			return nil, fmt.Errorf("failed to get certificate pool for %s: %w", registry, err)
		}

		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return client, nil
}

// TestConnection tests the TLS connection to a registry
func (m *trustStoreManager) TestConnection(registry string) error {
	if registry == "" {
		return fmt.Errorf("registry name cannot be empty")
	}

	client, err := m.CreateHTTPClient(registry)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Test connection with a simple GET request to /v2/
	// This is the standard registry API endpoint
	url := fmt.Sprintf("https://%s/v2/", registry)
	
	resp, err := client.Get(url)
	if err != nil {
		if m.IsInsecure(registry) {
			return fmt.Errorf("connection test failed for insecure registry %s: %w\n"+
				"This might indicate the registry is not running or accessible", registry, err)
		}
		return fmt.Errorf("TLS connection test failed for %s: %w\n"+
			"This might indicate:\n"+
			"1. Certificate mismatch or invalid certificate\n"+
			"2. Registry is not accessible\n"+
			"3. Network connectivity issues\n"+
			"Try using --insecure flag for testing", registry, err)
	}
	defer resp.Body.Close()

	// Check for expected registry API response
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized {
		// 200 OK or 401 Unauthorized are both valid responses for /v2/
		// 401 just means authentication is required, but the connection works
		return nil
	}

	return fmt.Errorf("unexpected response from registry %s: status %d (%s)",
		registry, resp.StatusCode, resp.Status)
}

// GetTLSConnectionInfo returns TLS connection information for debugging
func (m *trustStoreManager) GetTLSConnectionInfo(registry string) (*ConnectionInfo, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry name cannot be empty")
	}

	client, err := m.CreateHTTPClient(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	url := fmt.Sprintf("https://%s/v2/", registry)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", registry, err)
	}
	defer resp.Body.Close()

	// Extract TLS connection state
	if resp.TLS == nil {
		return &ConnectionInfo{
			Registry:    registry,
			IsSecure:    false,
			IsInsecure:  m.IsInsecure(registry),
			Error:       "No TLS connection state available",
		}, nil
	}

	info := &ConnectionInfo{
		Registry:         registry,
		IsSecure:         true,
		IsInsecure:       m.IsInsecure(registry),
		TLSVersion:       tlsVersionString(resp.TLS.Version),
		CipherSuite:      tls.CipherSuiteName(resp.TLS.CipherSuite),
		ServerCerts:      resp.TLS.PeerCertificates,
		VerifiedChains:   resp.TLS.VerifiedChains,
		HandshakeComplete: resp.TLS.HandshakeComplete,
	}

	return info, nil
}

// ConnectionInfo type is now defined in types.go

// tlsVersionString converts TLS version constant to string
func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (%d)", version)
	}
}