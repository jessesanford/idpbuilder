package registry

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// TransportConfig holds configuration for HTTP transport
type TransportConfig struct {
	// Timeout for HTTP requests
	Timeout time.Duration
	
	// MaxIdleConns controls the maximum number of idle connections
	MaxIdleConns int
	
	// MaxIdleConnsPerHost controls the maximum idle connections per host
	MaxIdleConnsPerHost int
	
	// IdleConnTimeout is the maximum amount of time an idle connection will remain idle
	IdleConnTimeout time.Duration
	
	// TLSHandshakeTimeout is the timeout for TLS handshakes
	TLSHandshakeTimeout time.Duration
	
	// ExpectContinueTimeout is the timeout for Expect: 100-continue headers
	ExpectContinueTimeout time.Duration
	
	// DisableCompression disables automatic gzip compression
	DisableCompression bool
}

// DefaultTransportConfig returns a default transport configuration
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		Timeout:               30 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
	}
}

// ConfigureTransportWithPhase1 configures HTTP transport using Phase 1 certificate infrastructure
func ConfigureTransportWithPhase1(trustStore certs.TrustStoreManager, registry string, insecure bool) (*http.Transport, error) {
	if trustStore == nil {
		return nil, fmt.Errorf("trust store manager cannot be nil")
	}
	
	// Configure insecure mode if requested
	if insecure {
		if err := trustStore.SetInsecureRegistry(registry, true); err != nil {
			return nil, fmt.Errorf("failed to configure insecure registry: %w", err)
		}
	}
	
	// Create HTTP client using Phase 1 infrastructure
	httpClient, err := trustStore.CreateHTTPClient(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}
	
	// Extract the transport
	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("unexpected transport type")
	}
	
	return transport, nil
}

// ConfigureTransportWithConfig creates HTTP transport with custom configuration
func ConfigureTransportWithConfig(trustStore certs.TrustStoreManager, registry string, insecure bool, config *TransportConfig) (*http.Transport, error) {
	if config == nil {
		config = DefaultTransportConfig()
	}
	
	// Get base transport from Phase 1
	transport, err := ConfigureTransportWithPhase1(trustStore, registry, insecure)
	if err != nil {
		return nil, err
	}
	
	// Apply custom configuration
	transport.MaxIdleConns = config.MaxIdleConns
	transport.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
	transport.IdleConnTimeout = config.IdleConnTimeout
	transport.TLSHandshakeTimeout = config.TLSHandshakeTimeout
	transport.ExpectContinueTimeout = config.ExpectContinueTimeout
	transport.DisableCompression = config.DisableCompression
	
	return transport, nil
}

// CreateInsecureTransport creates a transport that skips TLS verification
func CreateInsecureTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// CreateSecureTransport creates a transport with proper certificate validation
func CreateSecureTransport(certPool *tls.Config) *http.Transport {
	return &http.Transport{
		TLSClientConfig:       certPool,
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// ValidateTransportConfig validates transport configuration parameters
func ValidateTransportConfig(config *TransportConfig) error {
	if config == nil {
		return fmt.Errorf("transport config cannot be nil")
	}
	
	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	if config.MaxIdleConns < 0 {
		return fmt.Errorf("max idle connections cannot be negative")
	}
	
	if config.MaxIdleConnsPerHost < 0 {
		return fmt.Errorf("max idle connections per host cannot be negative")
	}
	
	if config.IdleConnTimeout <= 0 {
		return fmt.Errorf("idle connection timeout must be positive")
	}
	
	if config.TLSHandshakeTimeout <= 0 {
		return fmt.Errorf("TLS handshake timeout must be positive")
	}
	
	if config.ExpectContinueTimeout < 0 {
		return fmt.Errorf("expect continue timeout cannot be negative")
	}
	
	return nil
}

// ApplyTransportConfig applies configuration to an existing transport
func ApplyTransportConfig(transport *http.Transport, config *TransportConfig) error {
	if transport == nil {
		return fmt.Errorf("transport cannot be nil")
	}
	
	if err := ValidateTransportConfig(config); err != nil {
		return fmt.Errorf("invalid transport config: %w", err)
	}
	
	transport.MaxIdleConns = config.MaxIdleConns
	transport.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
	transport.IdleConnTimeout = config.IdleConnTimeout
	transport.TLSHandshakeTimeout = config.TLSHandshakeTimeout
	transport.ExpectContinueTimeout = config.ExpectContinueTimeout
	transport.DisableCompression = config.DisableCompression
	
	return nil
}

// CloneTransportConfig creates a copy of transport configuration
func CloneTransportConfig(config *TransportConfig) *TransportConfig {
	if config == nil {
		return DefaultTransportConfig()
	}
	
	return &TransportConfig{
		Timeout:               config.Timeout,
		MaxIdleConns:          config.MaxIdleConns,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
		DisableCompression:    config.DisableCompression,
	}
}