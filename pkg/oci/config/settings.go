package config

import (
	"os"
	"path/filepath"
)

// OCIConfig holds configuration for OCI operations
type OCIConfig struct {
	RegistryURL      string // Default: gitea.idpbuilder.localtest.me
	CacheDir         string // Default: ~/.idpbuilder/cache
	AutoExtractCerts bool   // Default: true
	InsecureMode     bool   // Default: false
}

// GetDefault returns default OCI configuration with environment variable overrides
func GetDefault() *OCIConfig {
	config := &OCIConfig{
		RegistryURL:      "gitea.idpbuilder.localtest.me",
		CacheDir:         getDefaultCacheDir(),
		AutoExtractCerts: true,
		InsecureMode:     false,
	}

	// Apply environment variable overrides
	if registryURL := os.Getenv("IDPBUILDER_OCI_REGISTRY"); registryURL != "" {
		config.RegistryURL = registryURL
	}

	if cacheDir := os.Getenv("IDPBUILDER_OCI_CACHE_DIR"); cacheDir != "" {
		config.CacheDir = cacheDir
	}

	if insecure := os.Getenv("IDPBUILDER_OCI_INSECURE"); insecure == "true" {
		config.InsecureMode = true
		config.AutoExtractCerts = false // Disable cert extraction in insecure mode
	}

	return config
}

// getDefaultCacheDir returns the default cache directory
func getDefaultCacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/idpbuilder-cache"
	}
	return filepath.Join(homeDir, ".idpbuilder", "cache")
}