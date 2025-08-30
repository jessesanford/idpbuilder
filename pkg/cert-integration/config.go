// Package cert_integration provides certificate integration between Phase 1 and Phase 2 components
package cert_integration

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ValidationMode defines the certificate validation strictness
type ValidationMode string

const (
	// ValidationModeStrict enforces full certificate validation
	ValidationModeStrict ValidationMode = "strict"
	// ValidationModePermissive allows some validation failures
	ValidationModePermissive ValidationMode = "permissive"
	// ValidationModeDisabled bypasses validation (for testing only)
	ValidationModeDisabled ValidationMode = "disabled"
)

// ManagerConfig contains configuration for the certificate manager
type ManagerConfig struct {
	// DefaultCertPath is the default path for certificate storage
	DefaultCertPath string `json:"default_cert_path" yaml:"default_cert_path"`
	
	// TrustStoreLocation is the location of the trust store
	TrustStoreLocation string `json:"trust_store_location" yaml:"trust_store_location"`
	
	// ValidationMode controls certificate validation behavior
	ValidationMode ValidationMode `json:"validation_mode" yaml:"validation_mode"`
	
	// CacheEnabled determines if certificate caching is active
	CacheEnabled bool `json:"cache_enabled" yaml:"cache_enabled"`
	
	// CacheTimeout specifies how long to cache certificates
	CacheTimeout time.Duration `json:"cache_timeout" yaml:"cache_timeout"`
}

// RegistryConfig contains configuration for registry operations
type RegistryConfig struct {
	// RegistryURL is the URL of the target registry
	RegistryURL string `json:"registry_url" yaml:"registry_url"`
	
	// Insecure allows insecure registry connections (for development)
	Insecure bool `json:"insecure" yaml:"insecure"`
	
	// CertPath is the path to registry-specific certificates
	CertPath string `json:"cert_path" yaml:"cert_path"`
	
	// AuthMethod specifies the authentication method
	AuthMethod string `json:"auth_method" yaml:"auth_method"`
	
	// Timeout for registry operations
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

// BuildConfig contains configuration for build operations
type BuildConfig struct {
	// BuildContext is the build context path
	BuildContext string `json:"build_context" yaml:"build_context"`
	
	// TrustStore is the trust store path for builds
	TrustStore string `json:"trust_store" yaml:"trust_store"`
	
	// CertOverride allows overriding certificates for specific operations
	CertOverride map[string]string `json:"cert_override" yaml:"cert_override"`
	
	// BuildTimeout specifies timeout for build operations
	BuildTimeout time.Duration `json:"build_timeout" yaml:"build_timeout"`
}

// CertificateSet contains a complete set of certificates for operations
type CertificateSet struct {
	// RootCA is the root certificate authority
	RootCA *x509.Certificate `json:"-"`
	
	// Intermediates contains intermediate CA certificates
	Intermediates []*x509.Certificate `json:"-"`
	
	// ServerCert is the server certificate
	ServerCert *x509.Certificate `json:"-"`
	
	// ClientCerts maps client names to their certificates
	ClientCerts map[string]*x509.Certificate `json:"-"`
	
	// TrustBundle contains the raw trust bundle data
	TrustBundle []byte `json:"-"`
	
	// Metadata contains additional certificate metadata
	Metadata map[string]interface{} `json:"metadata"`
}

// Validate validates the ManagerConfig
func (mc *ManagerConfig) Validate() error {
	if mc.DefaultCertPath == "" {
		return fmt.Errorf("default_cert_path cannot be empty")
	}
	
	if mc.TrustStoreLocation == "" {
		return fmt.Errorf("trust_store_location cannot be empty")
	}
	
	switch mc.ValidationMode {
	case ValidationModeStrict, ValidationModePermissive, ValidationModeDisabled:
		// Valid modes
	default:
		return fmt.Errorf("invalid validation_mode: %s", mc.ValidationMode)
	}
	
	if mc.CacheTimeout < 0 {
		return fmt.Errorf("cache_timeout cannot be negative")
	}
	
	return nil
}

// Validate validates the RegistryConfig
func (rc *RegistryConfig) Validate() error {
	if rc.RegistryURL == "" {
		return fmt.Errorf("registry_url cannot be empty")
	}
	
	if rc.CertPath != "" {
		if _, err := os.Stat(rc.CertPath); err != nil {
			return fmt.Errorf("cert_path validation failed: %w", err)
		}
	}
	
	if rc.Timeout <= 0 {
		rc.Timeout = 30 * time.Second // Default timeout
	}
	
	return nil
}

// Validate validates the BuildConfig
func (bc *BuildConfig) Validate() error {
	if bc.BuildContext == "" {
		return fmt.Errorf("build_context cannot be empty")
	}
	
	if bc.TrustStore != "" {
		if _, err := os.Stat(bc.TrustStore); err != nil {
			return fmt.Errorf("trust_store validation failed: %w", err)
		}
	}
	
	if bc.BuildTimeout <= 0 {
		bc.BuildTimeout = 10 * time.Minute // Default timeout
	}
	
	return nil
}

// Validate validates the CertificateSet
func (cs *CertificateSet) Validate() error {
	if cs.RootCA == nil {
		return fmt.Errorf("root CA certificate is required")
	}
	
	// Verify certificate chain
	if err := cs.verifyCertificateChain(); err != nil {
		return fmt.Errorf("certificate chain validation failed: %w", err)
	}
	
	// Check certificate expiry
	if err := cs.checkCertificateExpiry(); err != nil {
		return fmt.Errorf("certificate expiry validation failed: %w", err)
	}
	
	return nil
}

// verifyCertificateChain verifies the certificate chain integrity
func (cs *CertificateSet) verifyCertificateChain() error {
	if cs.RootCA == nil {
		return fmt.Errorf("root CA is missing")
	}
	
	// Create certificate pool with root CA
	pool := x509.NewCertPool()
	pool.AddCert(cs.RootCA)
	
	// Add intermediate certificates to the pool
	for _, intermediate := range cs.Intermediates {
		pool.AddCert(intermediate)
	}
	
	// Verify server certificate against the pool
	if cs.ServerCert != nil {
		opts := x509.VerifyOptions{
			Roots:         pool,
			Intermediates: x509.NewCertPool(),
		}
		
		for _, intermediate := range cs.Intermediates {
			opts.Intermediates.AddCert(intermediate)
		}
		
		_, err := cs.ServerCert.Verify(opts)
		if err != nil {
			return fmt.Errorf("server certificate verification failed: %w", err)
		}
	}
	
	return nil
}

// checkCertificateExpiry checks if certificates are expired or expiring soon
func (cs *CertificateSet) checkCertificateExpiry() error {
	now := time.Now()
	warningThreshold := 30 * 24 * time.Hour // 30 days warning
	
	// Check root CA expiry
	if cs.RootCA != nil && cs.RootCA.NotAfter.Before(now) {
		return fmt.Errorf("root CA certificate has expired on %v", cs.RootCA.NotAfter)
	}
	
	// Check intermediate certificates
	for i, intermediate := range cs.Intermediates {
		if intermediate.NotAfter.Before(now) {
			return fmt.Errorf("intermediate certificate %d has expired on %v", i, intermediate.NotAfter)
		}
	}
	
	// Check server certificate
	if cs.ServerCert != nil {
		if cs.ServerCert.NotAfter.Before(now) {
			return fmt.Errorf("server certificate has expired on %v", cs.ServerCert.NotAfter)
		}
		
		// Warning for soon-to-expire certificates
		if cs.ServerCert.NotAfter.Before(now.Add(warningThreshold)) {
			fmt.Printf("Warning: server certificate expires on %v\n", cs.ServerCert.NotAfter)
		}
	}
	
	return nil
}

// DefaultManagerConfig returns a default manager configuration
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		DefaultCertPath:    "/etc/idpbuilder/certs",
		TrustStoreLocation: "/etc/pki/ca-trust/extracted/pem",
		ValidationMode:     ValidationModeStrict,
		CacheEnabled:       true,
		CacheTimeout:       1 * time.Hour,
	}
}

// DefaultRegistryConfig returns a default registry configuration
func DefaultRegistryConfig(registryURL string) *RegistryConfig {
	return &RegistryConfig{
		RegistryURL: registryURL,
		Insecure:    false,
		CertPath:    filepath.Join("/etc/containers/certs.d", registryURL),
		AuthMethod:  "basic",
		Timeout:     30 * time.Second,
	}
}

// DefaultBuildConfig returns a default build configuration
func DefaultBuildConfig(buildContext string) *BuildConfig {
	return &BuildConfig{
		BuildContext:  buildContext,
		TrustStore:    "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
		CertOverride:  make(map[string]string),
		BuildTimeout:  10 * time.Minute,
	}
}