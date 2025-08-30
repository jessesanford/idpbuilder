package cert_integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathResolver defines the interface for resolving certificate paths
type PathResolver interface {
	// GetRegistryCertPath returns the certificate path for a specific registry
	GetRegistryCertPath(registry string) string
	
	// GetBuildTrustStore returns the path to the build trust store
	GetBuildTrustStore() string
	
	// ResolveCertificatePath resolves a certificate path from a hint or name
	ResolveCertificatePath(hint string) (string, error)
}

// pathResolver implements the PathResolver interface
type pathResolver struct {
	config        *ManagerConfig
	certificateLocations map[string]string
	fallbackPaths []string
}

// NewPathResolver creates a new path resolver with the given configuration
func NewPathResolver(config *ManagerConfig) PathResolver {
	resolver := &pathResolver{
		config:               config,
		certificateLocations: make(map[string]string),
	}
	
	// Initialize standard certificate locations
	resolver.initializeLocations()
	
	// Initialize fallback paths
	resolver.initializeFallbackPaths()
	
	return resolver
}

// GetRegistryCertPath returns the certificate path for a specific registry
func (pr *pathResolver) GetRegistryCertPath(registry string) string {
	if registry == "" {
		return pr.config.DefaultCertPath
	}
	
	// Clean and normalize registry name
	cleanRegistry := pr.cleanRegistryName(registry)
	
	// Check for specific registry mappings
	if path, exists := pr.certificateLocations[cleanRegistry]; exists {
		return path
	}
	
	// Use Docker/Podman standard location
	containersPath := filepath.Join("/etc/containers/certs.d", cleanRegistry)
	if pr.pathExists(containersPath) {
		return containersPath
	}
	
	// Use IDPBuilder-specific location
	idpBuilderPath := filepath.Join(pr.config.DefaultCertPath, "registries", cleanRegistry)
	if pr.pathExists(idpBuilderPath) {
		return idpBuilderPath
	}
	
	// Default to containers standard location (will be created if needed)
	return containersPath
}

// GetBuildTrustStore returns the path to the build trust store
func (pr *pathResolver) GetBuildTrustStore() string {
	// Check if trust store location is explicitly configured
	if pr.config.TrustStoreLocation != "" && pr.pathExists(pr.config.TrustStoreLocation) {
		return pr.config.TrustStoreLocation
	}
	
	// Try standard system trust store locations
	standardTrustStores := []string{
		"/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/ca-certificates.crt",
		"/usr/local/share/ca-certificates/ca-certificates.crt",
	}
	
	for _, trustStore := range standardTrustStores {
		if pr.pathExists(trustStore) {
			return trustStore
		}
	}
	
	// Fallback to IDPBuilder-specific trust store
	idpTrustStore := filepath.Join(pr.config.DefaultCertPath, "trust", "ca-bundle.pem")
	return idpTrustStore
}

// ResolveCertificatePath resolves a certificate path from a hint or name
func (pr *pathResolver) ResolveCertificatePath(hint string) (string, error) {
	if hint == "" {
		return "", fmt.Errorf("certificate path hint cannot be empty")
	}
	
	// If hint is already an absolute path, validate and return
	if filepath.IsAbs(hint) {
		if pr.pathExists(hint) {
			return hint, nil
		}
		return "", fmt.Errorf("absolute path does not exist: %s", hint)
	}
	
	// Check environment variable override
	if envPath := os.Getenv("IDPBUILDER_CERT_PATH"); envPath != "" {
		resolvedPath := filepath.Join(envPath, hint)
		if pr.pathExists(resolvedPath) {
			return resolvedPath, nil
		}
	}
	
	// Try predefined certificate locations
	if resolvedPath, exists := pr.certificateLocations[hint]; exists {
		if pr.pathExists(resolvedPath) {
			return resolvedPath, nil
		}
	}
	
	// Try fallback paths
	for _, basePath := range pr.fallbackPaths {
		candidates := []string{
			filepath.Join(basePath, hint),
			filepath.Join(basePath, hint+".crt"),
			filepath.Join(basePath, hint+".pem"),
			filepath.Join(basePath, "certs", hint),
			filepath.Join(basePath, "certificates", hint),
		}
		
		for _, candidate := range candidates {
			if pr.pathExists(candidate) {
				return candidate, nil
			}
		}
	}
	
	return "", fmt.Errorf("could not resolve certificate path for hint: %s", hint)
}

// initializeLocations initializes the standard certificate locations mapping
func (pr *pathResolver) initializeLocations() {
	// IDPBuilder-specific registries
	pr.certificateLocations["gitea"] = filepath.Join(pr.config.DefaultCertPath, "gitea")
	pr.certificateLocations["harbor"] = filepath.Join(pr.config.DefaultCertPath, "harbor")
	pr.certificateLocations["registry"] = filepath.Join(pr.config.DefaultCertPath, "registry")
	
	// Common registry hostnames
	pr.certificateLocations["docker.io"] = "/etc/containers/certs.d/docker.io"
	pr.certificateLocations["quay.io"] = "/etc/containers/certs.d/quay.io"
	pr.certificateLocations["gcr.io"] = "/etc/containers/certs.d/gcr.io"
	
	// Local/development registries
	pr.certificateLocations["localhost"] = filepath.Join(pr.config.DefaultCertPath, "localhost")
	pr.certificateLocations["kind-registry"] = filepath.Join(pr.config.DefaultCertPath, "kind", "registry")
}

// initializeFallbackPaths initializes the fallback paths for certificate resolution
func (pr *pathResolver) initializeFallbackPaths() {
	homeDir, _ := os.UserHomeDir()
	
	pr.fallbackPaths = []string{
		pr.config.DefaultCertPath,                              // Primary IDPBuilder cert path
		"/etc/idpbuilder/certs",                               // System IDPBuilder certs
		filepath.Join(homeDir, ".idpbuilder", "certs"),        // User IDPBuilder certs
		"/tmp/idpbuilder-certs",                               // Temporary certs
		"/etc/containers/certs.d",                             // Container runtime certs
		"/etc/ssl/certs",                                      // System SSL certs
		"/usr/local/share/ca-certificates",                    // Local CA certificates
		filepath.Join(homeDir, ".docker", "certs"),            // Docker user certs
		filepath.Join(homeDir, ".config", "containers", "certs"), // Podman user certs
	}
}

// cleanRegistryName cleans and normalizes a registry name for path usage
func (pr *pathResolver) cleanRegistryName(registry string) string {
	// Remove protocol prefixes
	cleaned := strings.TrimPrefix(registry, "https://")
	cleaned = strings.TrimPrefix(cleaned, "http://")
	
	// Remove trailing slashes
	cleaned = strings.TrimSuffix(cleaned, "/")
	
	// Remove port numbers for path generation (keep for reference)
	if colonIndex := strings.LastIndex(cleaned, ":"); colonIndex > 0 {
		// Check if this is actually a port (numeric after colon)
		possiblePort := cleaned[colonIndex+1:]
		if len(possiblePort) > 0 && isNumeric(possiblePort) {
			// For filesystem paths, replace colon with underscore
			cleaned = cleaned[:colonIndex] + "_" + possiblePort
		}
	}
	
	// Replace any remaining invalid path characters
	cleaned = strings.ReplaceAll(cleaned, ":", "_")
	cleaned = strings.ReplaceAll(cleaned, "/", "_")
	
	return cleaned
}

// pathExists checks if a path exists on the filesystem
func (pr *pathResolver) pathExists(path string) bool {
	if path == "" {
		return false
	}
	
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// isNumeric checks if a string contains only numeric characters
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return len(s) > 0
}