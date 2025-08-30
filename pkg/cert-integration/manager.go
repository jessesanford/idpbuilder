package cert_integration

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CertificateManager provides the main interface for certificate management and integration
type CertificateManager interface {
	// ConfigureRegistry configures registry operations with certificates
	ConfigureRegistry(registry string, certs *CertificateSet) error
	
	// ConfigureBuild configures build operations with certificates
	ConfigureBuild(context *BuildContext, certs *CertificateSet) error
	
	// LoadCertificates loads certificates from the specified source
	LoadCertificates(source string) (*CertificateSet, error)
	
	// ValidateIntegration validates the integration setup
	ValidateIntegration() error
	
	// ClearCache clears the certificate cache
	ClearCache()
	
	// GetCachedCertificates returns cached certificates for a source
	GetCachedCertificates(source string) (*CertificateSet, bool)
}

// BuildContext contains context information for build operations
type BuildContext struct {
	// BuildPath is the path to the build context
	BuildPath string
	
	// ContainerFile is the path to the Containerfile/Dockerfile
	ContainerFile string
	
	// OutputImage is the target image name
	OutputImage string
	
	// BuildArgs contains build arguments
	BuildArgs map[string]string
	
	// AdditionalTags contains additional tags for the image
	AdditionalTags []string
}

// certificateManager implements the CertificateManager interface
type certificateManager struct {
	config    *ManagerConfig
	loader    CertificateLoader
	resolver  PathResolver
	validator IntegrationValidator
	cache     map[string]*cachedCertificateSet
	cacheMux  sync.RWMutex
}

// cachedCertificateSet represents a cached certificate set with expiration
type cachedCertificateSet struct {
	certificates *CertificateSet
	cachedAt     time.Time
	expiresAt    time.Time
}

// NewCertificateManager creates a new certificate manager with the given configuration
func NewCertificateManager(config *ManagerConfig) (CertificateManager, error) {
	if config == nil {
		config = DefaultManagerConfig()
	}
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manager configuration: %w", err)
	}
	
	loader := NewCertificateLoader(config)
	resolver := NewPathResolver(config)
	validator := NewIntegrationValidator(config)
	
	return &certificateManager{
		config:    config,
		loader:    loader,
		resolver:  resolver,
		validator: validator,
		cache:     make(map[string]*cachedCertificateSet),
	}, nil
}

// ConfigureRegistry configures registry operations with the provided certificates
func (cm *certificateManager) ConfigureRegistry(registry string, certs *CertificateSet) error {
	if registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	
	if certs == nil {
		return fmt.Errorf("certificate set cannot be nil")
	}
	
	// Validate certificate set
	if err := certs.Validate(); err != nil {
		return fmt.Errorf("certificate validation failed: %w", err)
	}
	
	// Get the registry-specific certificate path
	certPath := cm.resolver.GetRegistryCertPath(registry)
	
	// Ensure the certificate directory exists
	if err := os.MkdirAll(certPath, 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}
	
	// Write certificates to the registry path
	if err := cm.writeCertificates(certPath, certs); err != nil {
		return fmt.Errorf("failed to write certificates: %w", err)
	}
	
	// Create registry configuration
	registryConfig := &RegistryConfig{
		RegistryURL: registry,
		CertPath:    certPath,
		AuthMethod:  "basic",
		Timeout:     30 * time.Second,
	}
	
	// Validate registry configuration
	if err := cm.validator.ValidateRegistryConnection(registryConfig); err != nil {
		return fmt.Errorf("registry validation failed: %w", err)
	}
	
	return nil
}

// ConfigureBuild configures build operations with the provided certificates
func (cm *certificateManager) ConfigureBuild(context *BuildContext, certs *CertificateSet) error {
	if context == nil {
		return fmt.Errorf("build context cannot be nil")
	}
	
	if certs == nil {
		return fmt.Errorf("certificate set cannot be nil")
	}
	
	// Validate certificate set
	if err := certs.Validate(); err != nil {
		return fmt.Errorf("certificate validation failed: %w", err)
	}
	
	// Get the build trust store location
	trustStorePath := cm.resolver.GetBuildTrustStore()
	
	// Set up buildah trust with certificates
	if err := cm.setupBuildahTrust(trustStorePath, certs); err != nil {
		return fmt.Errorf("failed to setup buildah trust: %w", err)
	}
	
	// Create build configuration
	buildConfig := &BuildConfig{
		BuildContext:  context.BuildPath,
		TrustStore:    trustStorePath,
		CertOverride:  make(map[string]string),
		BuildTimeout:  10 * time.Minute,
	}
	
	// Apply certificate overrides if specified
	if context.BuildArgs != nil {
		for key, value := range context.BuildArgs {
			if key == "cert_override" {
				buildConfig.CertOverride["default"] = value
			}
		}
	}
	
	// Validate build configuration
	if err := cm.validator.ValidateBuildConfiguration(buildConfig); err != nil {
		return fmt.Errorf("build configuration validation failed: %w", err)
	}
	
	return nil
}

// LoadCertificates loads certificates from the specified source
func (cm *certificateManager) LoadCertificates(source string) (*CertificateSet, error) {
	if source == "" {
		return nil, fmt.Errorf("certificate source cannot be empty")
	}
	
	// Check cache first if caching is enabled
	if cm.config.CacheEnabled {
		if cached, found := cm.getCachedCertificates(source); found {
			return cached, nil
		}
	}
	
	var certs *CertificateSet
	var err error
	
	// Determine the source type and load accordingly
	if filepath.IsAbs(source) {
		// Absolute path - load from filesystem
		certs, err = cm.loader.LoadFromPath(source)
	} else if source == "kind" || source[:5] == "kind:" {
		// Kind cluster source
		clusterName := "kind" // Default cluster name
		if len(source) > 5 {
			clusterName = source[5:] // Extract cluster name after "kind:"
		}
		certs, err = cm.loader.LoadFromKindCluster(clusterName)
	} else {
		// Try as extraction source first, then as path
		certs, err = cm.loader.LoadFromExtraction(source)
		if err != nil {
			// Fallback to path resolution
			resolvedPath, resolveErr := cm.resolver.ResolveCertificatePath(source)
			if resolveErr == nil {
				certs, err = cm.loader.LoadFromPath(resolvedPath)
			}
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates from source '%s': %w", source, err)
	}
	
	// Cache the certificates if caching is enabled
	if cm.config.CacheEnabled {
		cm.cacheOperations(source, certs)
	}
	
	return certs, nil
}

// ValidateIntegration validates the overall integration setup
func (cm *certificateManager) ValidateIntegration() error {
	return cm.validator.ValidateIntegration()
}

// ClearCache clears the certificate cache
func (cm *certificateManager) ClearCache() {
	cm.cacheMux.Lock()
	defer cm.cacheMux.Unlock()
	
	cm.cache = make(map[string]*cachedCertificateSet)
}

// GetCachedCertificates returns cached certificates for a source if available and valid
func (cm *certificateManager) GetCachedCertificates(source string) (*CertificateSet, bool) {
	return cm.getCachedCertificates(source)
}

// writeCertificates writes certificates to the specified path
func (cm *certificateManager) writeCertificates(certPath string, certs *CertificateSet) error {
	// Write root CA certificate
	if certs.RootCA != nil {
		caPath := filepath.Join(certPath, "ca.crt")
		if err := cm.writeCertificateFile(caPath, certs.RootCA); err != nil {
			return fmt.Errorf("failed to write CA certificate: %w", err)
		}
	}
	
	// Write intermediate certificates
	for i, intermediate := range certs.Intermediates {
		intermediatePath := filepath.Join(certPath, fmt.Sprintf("intermediate-%d.crt", i))
		if err := cm.writeCertificateFile(intermediatePath, intermediate); err != nil {
			return fmt.Errorf("failed to write intermediate certificate %d: %w", i, err)
		}
	}
	
	// Write server certificate
	if certs.ServerCert != nil {
		serverPath := filepath.Join(certPath, "server.crt")
		if err := cm.writeCertificateFile(serverPath, certs.ServerCert); err != nil {
			return fmt.Errorf("failed to write server certificate: %w", err)
		}
	}
	
	// Write trust bundle
	if len(certs.TrustBundle) > 0 {
		bundlePath := filepath.Join(certPath, "ca-bundle.pem")
		if err := os.WriteFile(bundlePath, certs.TrustBundle, 0644); err != nil {
			return fmt.Errorf("failed to write trust bundle: %w", err)
		}
	}
	
	return nil
}

// writeCertificateFile writes a single certificate to a file
func (cm *certificateManager) writeCertificateFile(filePath string, cert interface{}) error {
	// This would encode the certificate appropriately
	// For now, we'll create a placeholder implementation
	return fmt.Errorf("certificate writing not fully implemented: %s", filePath)
}

// setupBuildahTrust configures buildah trust store with certificates
func (cm *certificateManager) setupBuildahTrust(trustStorePath string, certs *CertificateSet) error {
	// Ensure trust store directory exists
	trustStoreDir := filepath.Dir(trustStorePath)
	if err := os.MkdirAll(trustStoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create trust store directory: %w", err)
	}
	
	// Write trust bundle to trust store
	if len(certs.TrustBundle) > 0 {
		if err := os.WriteFile(trustStorePath, certs.TrustBundle, 0644); err != nil {
			return fmt.Errorf("failed to write trust store: %w", err)
		}
	}
	
	return nil
}

// getCachedCertificates retrieves certificates from cache if available and valid
func (cm *certificateManager) getCachedCertificates(source string) (*CertificateSet, bool) {
	cm.cacheMux.RLock()
	defer cm.cacheMux.RUnlock()
	
	cached, exists := cm.cache[source]
	if !exists {
		return nil, false
	}
	
	// Check if cache entry has expired
	if time.Now().After(cached.expiresAt) {
		// Remove expired entry
		delete(cm.cache, source)
		return nil, false
	}
	
	return cached.certificates, true
}

// cacheOperations stores certificates in the cache
func (cm *certificateManager) cacheOperations(source string, certs *CertificateSet) {
	cm.cacheMux.Lock()
	defer cm.cacheMux.Unlock()
	
	now := time.Now()
	cached := &cachedCertificateSet{
		certificates: certs,
		cachedAt:     now,
		expiresAt:    now.Add(cm.config.CacheTimeout),
	}
	
	cm.cache[source] = cached
}