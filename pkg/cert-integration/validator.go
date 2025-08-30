package cert_integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// IntegrationValidator defines the interface for validating integration components
type IntegrationValidator interface {
	// ValidateIntegration validates the overall integration setup
	ValidateIntegration() error
	
	// ValidateRegistryConnection validates connectivity to a registry with certificates
	ValidateRegistryConnection(config *RegistryConfig) error
	
	// ValidateBuildConfiguration validates build configuration and dependencies
	ValidateBuildConfiguration(config *BuildConfig) error
}

// integrationValidator implements the IntegrationValidator interface
type integrationValidator struct {
	config *ManagerConfig
}

// NewIntegrationValidator creates a new integration validator
func NewIntegrationValidator(config *ManagerConfig) IntegrationValidator {
	return &integrationValidator{
		config: config,
	}
}

// ValidateIntegration validates the overall integration setup
func (iv *integrationValidator) ValidateIntegration() error {
	validationErrors := make([]string, 0)
	
	// Validate Phase 1 components accessibility
	if err := iv.validatePhase1Components(); err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Phase 1 validation: %v", err))
	}
	
	// Validate Phase 2 components compatibility
	if err := iv.validatePhase2Components(); err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Phase 2 validation: %v", err))
	}
	
	// Validate certificate infrastructure
	if err := iv.validateCertificateInfrastructure(); err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Certificate infrastructure: %v", err))
	}
	
	// Validate system dependencies
	if err := iv.validateSystemDependencies(); err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("System dependencies: %v", err))
	}
	
	if len(validationErrors) > 0 {
		return fmt.Errorf("integration validation failed: %s", strings.Join(validationErrors, "; "))
	}
	
	return nil
}

// ValidateRegistryConnection validates connectivity to a registry with certificates
func (iv *integrationValidator) ValidateRegistryConnection(config *RegistryConfig) error {
	if config == nil {
		return fmt.Errorf("registry config cannot be nil")
	}
	
	if err := config.Validate(); err != nil {
		return fmt.Errorf("registry config validation failed: %w", err)
	}
	
	// Parse registry URL
	registryURL := config.RegistryURL
	if !strings.HasPrefix(registryURL, "http://") && !strings.HasPrefix(registryURL, "https://") {
		registryURL = "https://" + registryURL
	}
	
	parsedURL, err := url.Parse(registryURL)
	if err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}
	
	// Test network connectivity
	if err := iv.testNetworkConnectivity(parsedURL.Host, config.Timeout); err != nil {
		return fmt.Errorf("network connectivity test failed: %w", err)
	}
	
	// Test TLS handshake if using HTTPS
	if parsedURL.Scheme == "https" {
		if err := iv.testTLSHandshake(config); err != nil {
			return fmt.Errorf("TLS handshake failed: %w", err)
		}
	}
	
	// Test HTTP connection
	if err := iv.testHTTPConnection(config); err != nil {
		return fmt.Errorf("HTTP connection test failed: %w", err)
	}
	
	return nil
}

// ValidateBuildConfiguration validates build configuration and dependencies
func (iv *integrationValidator) ValidateBuildConfiguration(config *BuildConfig) error {
	if config == nil {
		return fmt.Errorf("build config cannot be nil")
	}
	
	if err := config.Validate(); err != nil {
		return fmt.Errorf("build config validation failed: %w", err)
	}
	
	// Validate build context accessibility
	if err := iv.validateBuildContext(config); err != nil {
		return fmt.Errorf("build context validation failed: %w", err)
	}
	
	// Validate trust store accessibility
	if err := iv.validateTrustStore(config); err != nil {
		return fmt.Errorf("trust store validation failed: %w", err)
	}
	
	// Validate buildah compatibility
	if err := iv.validateBuildahCompatibility(); err != nil {
		return fmt.Errorf("buildah compatibility check failed: %w", err)
	}
	
	// Validate container runtime
	if err := iv.validateContainerRuntime(); err != nil {
		return fmt.Errorf("container runtime validation failed: %w", err)
	}
	
	return nil
}

// validatePhase1Components validates Phase 1 component accessibility
func (iv *integrationValidator) validatePhase1Components() error {
	// This would validate that Phase 1 interfaces are available
	// For now, we'll do basic checks for expected paths and configurations
	
	phase1Paths := []string{
		iv.config.DefaultCertPath,
		filepath.Join(iv.config.DefaultCertPath, "extraction"),
		filepath.Join(iv.config.DefaultCertPath, "kind"),
	}
	
	for _, path := range phase1Paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Create directory if it doesn't exist
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create Phase 1 path %s: %w", path, err)
			}
		}
	}
	
	return nil
}

// validatePhase2Components validates Phase 2 component compatibility
func (iv *integrationValidator) validatePhase2Components() error {
	// Check for buildah availability
	if _, err := exec.LookPath("buildah"); err != nil {
		return fmt.Errorf("buildah not found in PATH: %w", err)
	}
	
	// Check for standard build tools
	requiredTools := []string{"podman", "skopeo"}
	missingTools := make([]string, 0)
	
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			missingTools = append(missingTools, tool)
		}
	}
	
	if len(missingTools) > 0 {
		return fmt.Errorf("missing required Phase 2 tools: %s", strings.Join(missingTools, ", "))
	}
	
	return nil
}

// validateCertificateInfrastructure validates certificate infrastructure
func (iv *integrationValidator) validateCertificateInfrastructure() error {
	// Validate default certificate path
	if err := os.MkdirAll(iv.config.DefaultCertPath, 0755); err != nil {
		return fmt.Errorf("failed to create certificate path: %w", err)
	}
	
	// Validate trust store location
	trustStoreDir := filepath.Dir(iv.config.TrustStoreLocation)
	if err := os.MkdirAll(trustStoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create trust store directory: %w", err)
	}
	
	return nil
}

// validateSystemDependencies validates system-level dependencies
func (iv *integrationValidator) validateSystemDependencies() error {
	// Check for OpenSSL or similar certificate tools
	certTools := []string{"openssl", "certtool"}
	foundTool := false
	
	for _, tool := range certTools {
		if _, err := exec.LookPath(tool); err == nil {
			foundTool = true
			break
		}
	}
	
	if !foundTool {
		return fmt.Errorf("no certificate management tools found (tried: %s)", strings.Join(certTools, ", "))
	}
	
	return nil
}

// testNetworkConnectivity tests basic network connectivity to a host
func (iv *integrationValidator) testNetworkConnectivity(host string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	defer conn.Close()
	
	return nil
}

// testTLSHandshake tests TLS handshake with the registry
func (iv *integrationValidator) testTLSHandshake(config *RegistryConfig) error {
	// Create a custom TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.Insecure,
	}
	
	// Extract host from registry URL
	registryURL := config.RegistryURL
	if !strings.Contains(registryURL, "://") {
		registryURL = "https://" + registryURL
	}
	
	parsedURL, err := url.Parse(registryURL)
	if err != nil {
		return fmt.Errorf("failed to parse registry URL: %w", err)
	}
	
	// Test TLS connection
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: config.Timeout}, "tcp", parsedURL.Host, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer conn.Close()
	
	return nil
}

// testHTTPConnection tests HTTP connection to the registry
func (iv *integrationValidator) testHTTPConnection(config *RegistryConfig) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Insecure,
			},
		},
	}
	
	// Build test URL
	testURL := config.RegistryURL
	if !strings.HasPrefix(testURL, "http://") && !strings.HasPrefix(testURL, "https://") {
		testURL = "https://" + testURL
	}
	
	// Add /v2/ path for Docker registry API
	if !strings.HasSuffix(testURL, "/") {
		testURL += "/"
	}
	testURL += "v2/"
	
	// Make test request
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Accept any response that indicates the registry is responding
	// (authentication errors are expected for unauthenticated requests)
	if resp.StatusCode >= 500 {
		return fmt.Errorf("registry returned server error: %d", resp.StatusCode)
	}
	
	return nil
}

// validateBuildContext validates the build context accessibility
func (iv *integrationValidator) validateBuildContext(config *BuildConfig) error {
	if _, err := os.Stat(config.BuildContext); err != nil {
		return fmt.Errorf("build context path inaccessible: %w", err)
	}
	
	return nil
}

// validateTrustStore validates trust store accessibility
func (iv *integrationValidator) validateTrustStore(config *BuildConfig) error {
	if config.TrustStore == "" {
		return nil // Trust store is optional
	}
	
	if _, err := os.Stat(config.TrustStore); err != nil {
		return fmt.Errorf("trust store path inaccessible: %w", err)
	}
	
	return nil
}

// validateBuildahCompatibility validates buildah compatibility
func (iv *integrationValidator) validateBuildahCompatibility() error {
	cmd := exec.Command("buildah", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("buildah version check failed: %w", err)
	}
	
	return nil
}

// validateContainerRuntime validates container runtime availability
func (iv *integrationValidator) validateContainerRuntime() error {
	// Try podman first
	if cmd := exec.Command("podman", "version"); cmd.Run() == nil {
		return nil
	}
	
	// Fall back to docker
	if cmd := exec.Command("docker", "version"); cmd.Run() == nil {
		return nil
	}
	
	return fmt.Errorf("no compatible container runtime found (tried podman, docker)")
}