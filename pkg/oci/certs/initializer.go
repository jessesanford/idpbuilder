package certs

import (
	"context"
	"fmt"
	"os"
)

// TrustManager manages certificate trust for OCI operations
type TrustManager struct {
	config     *CertConfig
	configured bool
}

// InitializeTrustEnvironment initializes the trust environment for OCI operations.
// This function:
// 1. Auto-configures certificates if not already done
// 2. Loads existing trust configuration
// 3. Sets up environment variables for OCI tools
// 4. Returns configured trust manager
func InitializeTrustEnvironment(ctx context.Context) (*TrustManager, error) {
	// Step 1: Auto-configure if not already done
	config, err := AutoConfigureCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-configure certificates: %w", err)
	}

	// Step 2: Create trust manager
	trustManager := &TrustManager{
		config:     config,
		configured: false,
	}

	// Step 3: Set up environment variables for OCI tools
	if err := trustManager.setupEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to setup trust environment: %w", err)
	}

	trustManager.configured = true
	return trustManager, nil
}

// setupEnvironment configures environment variables for certificate trust
func (tm *TrustManager) setupEnvironment() error {
	if !tm.config.IsValid() {
		return fmt.Errorf("certificate configuration is not valid")
	}

	// Set SSL_CERT_FILE for buildah and other OCI tools
	if err := os.Setenv("SSL_CERT_FILE", tm.config.TrustStorePath); err != nil {
		return fmt.Errorf("failed to set SSL_CERT_FILE: %w", err)
	}

	// Set CURL_CA_BUNDLE for curl-based operations
	if err := os.Setenv("CURL_CA_BUNDLE", tm.config.TrustStorePath); err != nil {
		return fmt.Errorf("failed to set CURL_CA_BUNDLE: %w", err)
	}

	// Set REQUESTS_CA_BUNDLE for Python-based tools
	if err := os.Setenv("REQUESTS_CA_BUNDLE", tm.config.TrustStorePath); err != nil {
		return fmt.Errorf("failed to set REQUESTS_CA_BUNDLE: %w", err)
	}

	return nil
}

// GetCertificatePath returns the path to the certificate file
func (tm *TrustManager) GetCertificatePath() string {
	if tm.config == nil {
		return ""
	}
	return tm.config.TrustStorePath
}

// GetGiteaURL returns the URL of the Gitea registry
func (tm *TrustManager) GetGiteaURL() string {
	if tm.config == nil {
		return ""
	}
	return tm.config.GiteaURL
}

// IsConfigured returns true if the trust manager is properly configured
func (tm *TrustManager) IsConfigured() bool {
	return tm.configured && tm.config != nil && tm.config.IsValid()
}

// Reset clears the cached certificate configuration and forces re-extraction
func (tm *TrustManager) Reset() error {
	if tm.config == nil {
		return nil
	}

	// Remove cached configuration
	cacheFile := fmt.Sprintf("%s/cert-config.json", tm.config.CacheDir)
	if err := os.Remove(cacheFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cached config: %w", err)
	}

	// Clear environment variables
	os.Unsetenv("SSL_CERT_FILE")
	os.Unsetenv("CURL_CA_BUNDLE")
	os.Unsetenv("REQUESTS_CA_BUNDLE")

	tm.configured = false
	tm.config = nil

	return nil
}

// InitializeFallbackEnvironment initializes trust environment for insecure mode
func InitializeFallbackEnvironment() (*TrustManager, error) {
	trustManager := &TrustManager{
		config: &CertConfig{
			TrustStorePath: "",
			GiteaURL:       "https://gitea.idpbuilder.localtest.me",
			CertsExtracted: false,
			CacheDir:       getCacheDir(),
		},
		configured: true,
	}

	// Set environment variables to skip certificate verification
	os.Setenv("BUILDAH_TLS_VERIFY", "false")
	os.Setenv("PODMAN_TLS_VERIFY", "false")

	return trustManager, nil
}