package certs

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DefaultTrustStoreConfig returns a default configuration for the trust store
func DefaultTrustStoreConfig() TrustStoreConfig {
	homeDir, _ := os.UserHomeDir()
	return TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         filepath.Join(homeDir, ".config", "containers", "certs.d"),
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
}

// manager implements the TrustManager interface
type manager struct {
	config      TrustStoreConfig
	store       CertificateStore
	registryMgr RegistryConfigManager
}

// NewTrustManager creates a new trust manager with the given configuration
func NewTrustManager(config TrustStoreConfig) TrustManager {
	store := newFileStore(config)
	registryMgr := newRegistryConfigManager(config)
	
	return &manager{
		config:      config,
		store:       store,
		registryMgr: registryMgr,
	}
}

// AddCertificate adds a certificate to the trust store for a specific registry
func (m *manager) AddCertificate(ctx context.Context, registry string, cert *Certificate) error {
	if registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	
	// Parse the certificate to extract metadata
	if err := m.parseCertificateInfo(cert); err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}
	
	// Store the certificate
	if err := m.store.Store(registry, cert); err != nil {
		return fmt.Errorf("failed to store certificate: %w", err)
	}
	
	return nil
}

// RemoveCertificate removes a certificate from the trust store
func (m *manager) RemoveCertificate(ctx context.Context, registry string, fingerprint string) error {
	if registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	if fingerprint == "" {
		return fmt.Errorf("fingerprint cannot be empty")
	}
	
	exists, err := m.store.Exists(registry, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to check certificate existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("certificate with fingerprint %s not found for registry %s", fingerprint, registry)
	}
	
	return m.store.Delete(registry, fingerprint)
}

// ListCertificates lists all certificates for a specific registry
func (m *manager) ListCertificates(ctx context.Context, registry string) ([]Certificate, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry cannot be empty")
	}
	
	return m.store.List(registry)
}

// GetRegistryConfig gets the complete configuration for a registry
func (m *manager) GetRegistryConfig(ctx context.Context, registry string) (*RegistryConfig, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry cannot be empty")
	}
	
	certs, err := m.store.List(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to list certificates: %w", err)
	}
	
	insecureRegistries, err := m.registryMgr.GetInsecureRegistries()
	if err != nil {
		return nil, fmt.Errorf("failed to get insecure registries: %w", err)
	}
	
	insecure := false
	for _, insecureReg := range insecureRegistries {
		if insecureReg == registry {
			insecure = true
			break
		}
	}
	
	return &RegistryConfig{
		Registry:     registry,
		Insecure:     insecure,
		Certificates: certs,
	}, nil
}

// SetInsecureRegistry configures a registry to skip TLS verification
func (m *manager) SetInsecureRegistry(ctx context.Context, registry string, insecure bool) error {
	if registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	
	return m.registryMgr.UpdateInsecureRegistry(registry, insecure)
}

// ValidateCertificate validates a certificate against the trust store
func (m *manager) ValidateCertificate(ctx context.Context, registry string, cert *x509.Certificate) error {
	if registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	
	certs, err := m.store.List(registry)
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}
	
	// Calculate fingerprint of the provided certificate
	fingerprint := calculateFingerprint(cert.Raw)
	
	// Check if the certificate is in our trust store
	for _, trustCert := range certs {
		if trustCert.Info.Fingerprint == fingerprint {
			return nil // Certificate is trusted
		}
	}
	
	return fmt.Errorf("certificate not found in trust store for registry %s", registry)
}

// parseCertificateInfo extracts metadata from a certificate
func (m *manager) parseCertificateInfo(cert *Certificate) error {
	block, _ := pem.Decode(cert.Data)
	if block == nil {
		return fmt.Errorf("failed to decode PEM certificate")
	}
	
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse X.509 certificate: %w", err)
	}
	
	cert.Info = CertificateInfo{
		Subject:      x509Cert.Subject.String(),
		Issuer:       x509Cert.Issuer.String(),
		SerialNumber: x509Cert.SerialNumber.String(),
		NotBefore:    x509Cert.NotBefore.Format(time.RFC3339),
		NotAfter:     x509Cert.NotAfter.Format(time.RFC3339),
		Fingerprint:  calculateFingerprint(x509Cert.Raw),
	}
	
	return nil
}

// calculateFingerprint calculates the SHA256 fingerprint of certificate data
func calculateFingerprint(data []byte) string {
	hash := sha256.Sum256(data)
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}