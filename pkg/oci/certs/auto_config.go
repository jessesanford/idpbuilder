package certs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// AutoConfigureCertificates automatically detects and configures certificates
// from the local Kind cluster for OCI operations.
func AutoConfigureCertificates(ctx context.Context) (*CertConfig, error) {
	// Check if certificates are already cached
	if config, err := loadCachedConfig(); err == nil && config.IsValid() {
		return config, nil
	}

	// Step 1: Check if Kind cluster exists
	cluster, err := detectKindCluster(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect Kind cluster: %w", err)
	}

	// Step 2: Extract Gitea certificates from cluster
	certs, err := extractGiteaCertificates(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certificates: %w", err)
	}

	// Step 3: Configure trust store
	trustStore, err := configureTrustStore(certs)
	if err != nil {
		return nil, fmt.Errorf("failed to configure trust store: %w", err)
	}

	// Step 4: Create and cache configuration
	config := &CertConfig{
		TrustStorePath: trustStore,
		GiteaURL:       cluster.GiteaURL,
		CertsExtracted: true,
		CacheDir:       getCacheDir(),
	}

	if err := config.Cache(); err != nil {
		// Non-fatal error - continue without caching
		fmt.Fprintf(os.Stderr, "Warning: failed to cache certificate config: %v\n", err)
	}

	return config, nil
}

// CertConfig holds certificate configuration for OCI operations
type CertConfig struct {
	TrustStorePath string `json:"trust_store_path"`
	GiteaURL       string `json:"gitea_url"`
	CertsExtracted bool   `json:"certs_extracted"`
	CacheDir       string `json:"cache_dir"`
}

// IsValid checks if the certificate configuration is still valid
func (c *CertConfig) IsValid() bool {
	if !c.CertsExtracted {
		return false
	}
	
	// Check if trust store file still exists
	if _, err := os.Stat(c.TrustStorePath); os.IsNotExist(err) {
		return false
	}
	
	return true
}

// Cache saves the certificate configuration to disk for reuse
func (c *CertConfig) Cache() error {
	cacheFile := filepath.Join(c.CacheDir, "cert-config.json")
	
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Note: In a real implementation, this would serialize to JSON
	// For now, we'll create a simple marker file
	file, err := os.Create(cacheFile)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "trust_store=%s\ngitea_url=%s\n", c.TrustStorePath, c.GiteaURL)
	return nil
}

// loadCachedConfig attempts to load previously cached certificate configuration
func loadCachedConfig() (*CertConfig, error) {
	cacheDir := getCacheDir()
	cacheFile := filepath.Join(cacheDir, "cert-config.json")
	
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("no cached config found")
	}
	
	// In a real implementation, this would deserialize from JSON
	// For now, return a placeholder that will fail IsValid()
	return &CertConfig{CertsExtracted: false}, nil
}

// detectKindCluster detects the running Kind cluster
func detectKindCluster(ctx context.Context) (*ClusterInfo, error) {
	// TODO: Implement actual Kind cluster detection using existing utilities
	// For now, return a placeholder cluster info
	return &ClusterInfo{
		Name:       "idpbuilder",
		GiteaURL: "https://gitea.idpbuilder.localtest.me", // Default Gitea URL
		KubeConfig: "~/.kube/config",
	}, nil
}

// ClusterInfo holds information about the detected Kind cluster
type ClusterInfo struct {
	Name       string
	GiteaURL   string
	KubeConfig string
}

// extractGiteaCertificates extracts certificates from the Kind cluster
func extractGiteaCertificates(ctx context.Context, cluster *ClusterInfo) ([]byte, error) {
	// In a real implementation, this would:
	// 1. Connect to the Kind cluster
	// 2. Extract the Gitea TLS certificate
	// 3. Return the certificate data
	
	// Placeholder for actual certificate extraction
	certData := []byte("") // TODO: Implement actual extraction
	
	return certData, nil
}

// configureTrustStore configures the system trust store with extracted certificates
func configureTrustStore(certData []byte) (string, error) {
	cacheDir := getCacheDir()
	trustStoreDir := filepath.Join(cacheDir, "certs")
	
	if err := os.MkdirAll(trustStoreDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create trust store directory: %w", err)
	}
	
	certFile := filepath.Join(trustStoreDir, "gitea-ca.crt")
	if err := os.WriteFile(certFile, certData, 0644); err != nil {
		return "", fmt.Errorf("failed to write certificate file: %w", err)
	}
	
	return certFile, nil
}

// getCacheDir returns the cache directory for certificate configuration
func getCacheDir() string {
	if cacheDir := os.Getenv("IDPBUILDER_CACHE_DIR"); cacheDir != "" {
		return cacheDir
	}
	
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/idpbuilder-cache"
	}
	
	return filepath.Join(homeDir, ".idpbuilder", "cache")
}