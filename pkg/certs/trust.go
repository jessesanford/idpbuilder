// Package certs provides TLS certificate management for registry operations
package certs

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Types are now consolidated in types.go - using those definitions

// trustStoreManager implements the TrustStoreManager interface
type trustStoreManager struct {
	certsDir           string
	trustedCerts       map[string][]*x509.Certificate
	insecureRegistries map[string]bool
	mu                 sync.RWMutex
	initialized        bool
}

// NewTrustStoreManager creates a new TrustStoreManager instance
func NewTrustStoreManager(certsDir string) (TrustStoreManager, error) {
	// Check if certificate management is enabled via feature flag
	if os.Getenv("ENABLE_CERT_MANAGEMENT") != "true" {
		return nil, fmt.Errorf("certificate management is disabled - set ENABLE_CERT_MANAGEMENT=true to enable")
	}

	if certsDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		certsDir = filepath.Join(homeDir, ".idpbuilder", "certs")
	}

	manager := &trustStoreManager{
		certsDir:           certsDir,
		trustedCerts:       make(map[string][]*x509.Certificate),
		insecureRegistries: make(map[string]bool),
	}

	// Ensure the certs directory exists
	if err := os.MkdirAll(certsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create certificates directory %s: %w", certsDir, err)
	}

	// Load existing certificates from disk
	if err := manager.LoadFromDisk(); err != nil {
		return nil, fmt.Errorf("failed to load certificates from disk: %w", err)
	}

	manager.initialized = true
	return manager, nil
}

// AddCertificate adds a certificate for a specific registry
func (m *trustStoreManager) AddCertificate(registry string, cert *x509.Certificate) error {
	if registry == "" {
		return fmt.Errorf("registry name cannot be empty")
	}
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	// Validate certificate
	if time.Now().After(cert.NotAfter) {
		return fmt.Errorf("certificate for %s has expired on %s", registry, cert.NotAfter.Format(time.RFC3339))
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Initialize the slice if it doesn't exist
	if m.trustedCerts[registry] == nil {
		m.trustedCerts[registry] = make([]*x509.Certificate, 0, 1)
	}

	// Check if certificate already exists
	for _, existingCert := range m.trustedCerts[registry] {
		if existingCert.Equal(cert) {
			return nil // Certificate already exists, no error
		}
	}

	// Add certificate to memory
	m.trustedCerts[registry] = append(m.trustedCerts[registry], cert)

	// Persist to disk
	return m.saveToDiskLocked(registry, cert)
}

// RemoveCertificate removes the certificate for a registry
func (m *trustStoreManager) RemoveCertificate(registry string) error {
	if registry == "" {
		return fmt.Errorf("registry name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove from memory
	delete(m.trustedCerts, registry)

	// Remove from disk
	certFile := filepath.Join(m.certsDir, registry+".pem")
	if err := os.Remove(certFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove certificate file for %s: %w", registry, err)
	}

	return nil
}

// SetInsecureRegistry marks a registry as insecure (skip TLS verification)
func (m *trustStoreManager) SetInsecureRegistry(registry string, insecure bool) error {
	if registry == "" {
		return fmt.Errorf("registry name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if insecure {
		m.insecureRegistries[registry] = true
	} else {
		delete(m.insecureRegistries, registry)
	}

	return nil
}

// GetTrustedCerts returns all trusted certificates for a registry
func (m *trustStoreManager) GetTrustedCerts(registry string) ([]*x509.Certificate, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry name cannot be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	certs := m.trustedCerts[registry]
	if len(certs) == 0 {
		return nil, nil // No certificates found, not an error
	}

	// Return a copy to prevent external modification
	result := make([]*x509.Certificate, len(certs))
	copy(result, certs)
	return result, nil
}

// GetCertPool returns a configured cert pool for a registry
func (m *trustStoreManager) GetCertPool(registry string) (*x509.CertPool, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry name cannot be empty")
	}

	// Start with system root CAs
	pool, err := x509.SystemCertPool()
	if err != nil {
		// Fall back to empty pool if system certs not available
		pool = x509.NewCertPool()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Add custom certificates for this registry
	certs := m.trustedCerts[registry]
	for _, cert := range certs {
		// Check if certificate is still valid
		if time.Now().After(cert.NotAfter) {
			return nil, fmt.Errorf("certificate for %s has expired on %s\n"+
				"Please update the certificate or use --insecure flag for testing",
				registry, cert.NotAfter.Format(time.RFC3339))
		}
		pool.AddCert(cert)
	}

	return pool, nil
}

// IsInsecure checks if a registry is marked as insecure
func (m *trustStoreManager) IsInsecure(registry string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.insecureRegistries[registry]
}

// LoadFromDisk loads all certificates from persistent storage
func (m *trustStoreManager) LoadFromDisk() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if certificates directory exists
	if _, err := os.Stat(m.certsDir); os.IsNotExist(err) {
		// Directory doesn't exist yet, that's okay
		return nil
	}

	// Read all .pem files in the certificates directory
	files, err := ioutil.ReadDir(m.certsDir)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied accessing certificate store at %s\n"+
				"Please check file permissions or run with appropriate privileges", m.certsDir)
		}
		return fmt.Errorf("failed to read certificates directory %s: %w", m.certsDir, err)
	}

	for _, file := range files {
		if !file.Mode().IsRegular() || filepath.Ext(file.Name()) != ".pem" {
			continue
		}

		// Extract registry name from filename
		registry := file.Name()[:len(file.Name())-4] // Remove .pem extension

		certPath := filepath.Join(m.certsDir, file.Name())
		if err := m.loadCertificateFromFile(registry, certPath); err != nil {
			// Log error but continue with other certificates
			fmt.Fprintf(os.Stderr, "Warning: failed to load certificate for %s: %v\n", registry, err)
		}
	}

	return nil
}

// SaveToDisk saves a certificate to persistent storage
func (m *trustStoreManager) SaveToDisk(registry string, cert *x509.Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.saveToDiskLocked(registry, cert)
}

// saveToDiskLocked saves a certificate to disk (caller must hold write lock)
func (m *trustStoreManager) saveToDiskLocked(registry string, cert *x509.Certificate) error {
	certFile := filepath.Join(m.certsDir, registry+".pem")

	// Create PEM block
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	// Write to file with restrictive permissions
	pemData := pem.EncodeToMemory(pemBlock)
	if err := ioutil.WriteFile(certFile, pemData, 0600); err != nil {
		return fmt.Errorf("failed to save certificate for %s to %s: %w", registry, certFile, err)
	}

	return nil
}

// loadCertificateFromFile loads a certificate from a PEM file
func (m *trustStoreManager) loadCertificateFromFile(registry, certPath string) error {
	// Read the PEM file
	pemData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate file %s: %w", certPath, err)
	}

	// Parse PEM blocks
	var certs []*x509.Certificate
	for len(pemData) > 0 {
		block, rest := pem.Decode(pemData)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			pemData = rest
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate in %s: %w", certPath, err)
		}

		certs = append(certs, cert)
		pemData = rest
	}

	if len(certs) == 0 {
		return fmt.Errorf("no valid certificates found in %s", certPath)
	}

	// Store in memory
	m.trustedCerts[registry] = certs
	return nil
}