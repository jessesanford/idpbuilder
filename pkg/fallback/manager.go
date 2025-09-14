package fallback

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"sync"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// FallbackManager handles fallback scenarios for registry connections
// when certificate validation fails or insecure connections are required
type FallbackManager struct {
	mu           sync.RWMutex
	trustStore   *certs.DefaultTrustStore
	insecureMode bool
	enabled      bool
}

// NewFallbackManager creates a new fallback manager with the given trust store
func NewFallbackManager(trustStore *certs.DefaultTrustStore) *FallbackManager {
	return &FallbackManager{
		trustStore:   trustStore,
		insecureMode: false,
		enabled:      true,
	}
}

// SetInsecureMode enables or disables insecure mode for connections
func (fm *FallbackManager) SetInsecureMode(insecure bool) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.insecureMode = insecure
}

// IsInsecureMode returns whether insecure mode is enabled
func (fm *FallbackManager) IsInsecureMode() bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.insecureMode
}

// CreateHTTPClient creates an HTTP client with appropriate TLS configuration
func (fm *FallbackManager) CreateHTTPClient() (*http.Client, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if !fm.enabled {
		return nil, fmt.Errorf("fallback manager disabled")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: fm.insecureMode,
		},
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

// ConfigureTLS configures TLS settings based on fallback policies
func (fm *FallbackManager) ConfigureTLS(registry string) (*tls.Config, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if !fm.enabled {
		return nil, fmt.Errorf("fallback manager disabled")
	}

	config := &tls.Config{
		InsecureSkipVerify: fm.insecureMode,
	}

	// If not in insecure mode, try to get trusted certificates from trust store
	if !fm.insecureMode && fm.trustStore != nil {
		if trustedCerts, err := fm.trustStore.GetTrustedCerts(registry); err == nil && len(trustedCerts) > 0 {
			// Create certificate pool and add trusted certificates
			pool, err := x509.SystemCertPool()
			if err != nil {
				pool = x509.NewCertPool()
			}

			for _, cert := range trustedCerts {
				pool.AddCert(cert)
			}

			config.RootCAs = pool
		}
	}

	return config, nil
}

// Enable enables the fallback manager
func (fm *FallbackManager) Enable() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.enabled = true
}

// Disable disables the fallback manager
func (fm *FallbackManager) Disable() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.enabled = false
}

// IsEnabled returns whether the fallback manager is enabled
func (fm *FallbackManager) IsEnabled() bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.enabled
}