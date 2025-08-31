package certs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// fileStore implements the CertificateStore interface using filesystem storage
type fileStore struct {
	config TrustStoreConfig
}

// newFileStore creates a new file-based certificate store
func newFileStore(config TrustStoreConfig) CertificateStore {
	return &fileStore{config: config}
}

// Store writes a certificate to the filesystem
func (fs *fileStore) Store(registry string, cert *Certificate) error {
	registryDir := fs.getRegistryDir(registry)

	// Create registry directory if it doesn't exist
	if err := os.MkdirAll(registryDir, fs.config.DirPermissions); err != nil {
		return fmt.Errorf("failed to create registry directory: %w", err)
	}

	// Generate filename from fingerprint
	filename := fmt.Sprintf("%s.crt", cert.Info.Fingerprint)
	certPath := filepath.Join(registryDir, filename)
	cert.FilePath = certPath

	// Write certificate using atomic operation (write to temp file, then rename)
	tempPath := certPath + ".tmp"
	if err := os.WriteFile(tempPath, cert.Data, fs.config.FilePermissions); err != nil {
		return fmt.Errorf("failed to write certificate to temp file: %w", err)
	}

	if err := os.Rename(tempPath, certPath); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to move certificate to final location: %w", err)
	}

	return nil
}

// Load reads a certificate from the filesystem
func (fs *fileStore) Load(registry string, fingerprint string) (*Certificate, error) {
	registryDir := fs.getRegistryDir(registry)
	filename := fmt.Sprintf("%s.crt", fingerprint)
	certPath := filepath.Join(registryDir, filename)

	data, err := os.ReadFile(certPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("certificate not found")
		}
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	cert := &Certificate{
		Data:     data,
		FilePath: certPath,
	}

	// Parse certificate info from data
	if err := fs.parseCertificateFromData(cert); err != nil {
		return nil, fmt.Errorf("failed to parse certificate info: %w", err)
	}

	return cert, nil
}

// Delete removes a certificate from the filesystem
func (fs *fileStore) Delete(registry string, fingerprint string) error {
	registryDir := fs.getRegistryDir(registry)
	filename := fmt.Sprintf("%s.crt", fingerprint)
	certPath := filepath.Join(registryDir, filename)

	if err := os.Remove(certPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("certificate not found")
		}
		return fmt.Errorf("failed to delete certificate: %w", err)
	}

	// Try to remove registry directory if it's empty
	fs.cleanupEmptyDir(registryDir)

	return nil
}

// List returns all certificates for a registry
func (fs *fileStore) List(registry string) ([]Certificate, error) {
	registryDir := fs.getRegistryDir(registry)

	// Check if registry directory exists
	if _, err := os.Stat(registryDir); os.IsNotExist(err) {
		return []Certificate{}, nil // No certificates for this registry
	}

	entries, err := os.ReadDir(registryDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry directory: %w", err)
	}

	var certificates []Certificate
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".crt") {
			continue
		}

		certPath := filepath.Join(registryDir, entry.Name())
		data, err := os.ReadFile(certPath)
		if err != nil {
			continue // Skip files that can't be read
		}

		cert := Certificate{
			Data:     data,
			FilePath: certPath,
		}

		// Parse certificate info
		if err := fs.parseCertificateFromData(&cert); err != nil {
			continue // Skip invalid certificates
		}

		certificates = append(certificates, cert)
	}

	return certificates, nil
}

// Exists checks if a certificate exists in the store
func (fs *fileStore) Exists(registry string, fingerprint string) (bool, error) {
	registryDir := fs.getRegistryDir(registry)
	filename := fmt.Sprintf("%s.crt", fingerprint)
	certPath := filepath.Join(registryDir, filename)

	_, err := os.Stat(certPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check certificate existence: %w", err)
	}

	return true, nil
}

// getRegistryDir returns the directory path for a given registry
func (fs *fileStore) getRegistryDir(registry string) string {
	// Replace colons with underscores in registry names for filesystem compatibility
	safeName := strings.ReplaceAll(registry, ":", "_")
	return filepath.Join(fs.config.BaseDir, safeName)
}

// cleanupEmptyDir removes a directory if it's empty
func (fs *fileStore) cleanupEmptyDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) > 0 {
		return
	}
	os.Remove(dir) // Best effort - ignore errors
}

// parseCertificateFromData is a helper to parse certificate info from PEM data
func (fs *fileStore) parseCertificateFromData(cert *Certificate) error {
	// Create a temporary manager to use its parsing function
	tempMgr := &manager{}
	return tempMgr.parseCertificateInfo(cert)
}
