// Package certificates provides certificate management services for OCI registry operations.
// This package implements the CertificateService interface with thread-safe operations,
// dynamic verification mode switching, and Gitea-specific certificate handling.
package certificates

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// CertificateServiceImpl implements the CertificateService interface
// with thread-safe operations and dynamic verification mode management.
type CertificateServiceImpl struct {
	// Thread safety
	mu sync.RWMutex

	// Certificate pools
	certPool   *x509.CertPool // Current certificate pool
	systemPool *x509.CertPool // System certificate pool
	customPool *x509.CertPool // Custom certificate pool

	// Configuration
	verificationMode v2.VerificationMode // Current verification mode
	bundlePaths      []string            // Certificate bundle paths

	// Components
	giteaDiscovery  *GiteaDiscovery      // Gitea certificate discovery
	verificationMgr *VerificationManager // Verification mode management

	// Internal state
	certificates map[string]*x509.Certificate // Certificate cache by fingerprint
	lastUpdate   time.Time                    // Last certificate update time
}

// NewCertificateService creates a new certificate service instance
// with default configuration and system certificate pool initialization.
func NewCertificateService() (*CertificateServiceImpl, error) {
	// Get system certificate pool
	systemPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system certificate pool: %w", err)
	}

	// Create custom certificate pool
	customPool := x509.NewCertPool()

	// Initialize service
	service := &CertificateServiceImpl{
		certPool:         systemPool.Clone(),
		systemPool:       systemPool,
		customPool:       customPool,
		verificationMode: v2.VerificationModeStrict,
		certificates:     make(map[string]*x509.Certificate),
		bundlePaths:      make([]string, 0),
		lastUpdate:       time.Now(),
	}

	// Initialize Gitea discovery
	giteaDiscovery, err := NewGiteaDiscovery()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gitea discovery: %w", err)
	}
	service.giteaDiscovery = giteaDiscovery

	// Initialize verification manager
	verificationMgr, err := NewVerificationManager(systemPool, customPool)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize verification manager: %w", err)
	}
	service.verificationMgr = verificationMgr

	return service, nil
}

// LoadCertificateBundle loads certificates from the specified bundle file.
// Supports PEM and DER formats with proper error handling and logging.
func (s *CertificateServiceImpl) LoadCertificateBundle(ctx context.Context, bundlePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if file exists
	if _, err := os.Stat(bundlePath); os.IsNotExist(err) {
		return fmt.Errorf("certificate bundle not found: %s", bundlePath)
	}

	// Read bundle file
	data, err := ioutil.ReadFile(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to read certificate bundle %s: %w", bundlePath, err)
	}

	// Parse certificates based on format
	certs, err := s.parseCertificates(data)
	if err != nil {
		return fmt.Errorf("failed to parse certificates from %s: %w", bundlePath, err)
	}

	// Add certificates to pools
	for _, cert := range certs {
		fingerprint := s.getCertificateFingerprint(cert)
		s.certificates[fingerprint] = cert
		s.customPool.AddCert(cert)
	}

	// Add bundle path to tracking
	s.bundlePaths = append(s.bundlePaths, bundlePath)
	s.lastUpdate = time.Now()

	// Rebuild certificate pool based on current mode
	err = s.rebuildCertPool()
	if err != nil {
		return fmt.Errorf("failed to rebuild certificate pool: %w", err)
	}

	return nil
}

// SetVerificationMode sets the certificate verification mode with thread safety.
// Supports atomic mode switching while preserving certificate state.
func (s *CertificateServiceImpl) SetVerificationMode(ctx context.Context, mode v2.VerificationMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate mode
	if !s.isValidVerificationMode(mode) {
		return fmt.Errorf("invalid verification mode: %s", mode)
	}

	// Update verification manager
	err := s.verificationMgr.SwitchVerificationMode(mode)
	if err != nil {
		return fmt.Errorf("failed to switch verification mode: %w", err)
	}

	// Update internal state
	oldMode := s.verificationMode
	s.verificationMode = mode

	// Rebuild certificate pool for new mode
	err = s.rebuildCertPool()
	if err != nil {
		// Rollback on failure
		s.verificationMode = oldMode
		s.verificationMgr.SwitchVerificationMode(oldMode)
		return fmt.Errorf("failed to rebuild certificate pool for mode %s: %w", mode, err)
	}

	return nil
}

// ValidateCertificate validates a certificate against the current configuration.
// Performs comprehensive validation including validity period, signature, and key usage.
func (s *CertificateServiceImpl) ValidateCertificate(ctx context.Context, cert *x509.Certificate) (*v2.ValidationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Handle nil certificate
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	result := &v2.ValidationResult{
		Valid:    true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
		Certificate: &v2.CertificateInfo{
			Subject:     cert.Subject.String(),
			Issuer:      cert.Issuer.String(),
			NotBefore:   cert.NotBefore.Format(time.RFC3339),
			NotAfter:    cert.NotAfter.Format(time.RFC3339),
			KeyUsage:    fmt.Sprintf("%v", cert.KeyUsage),
			Fingerprint: s.getCertificateFingerprint(cert),
		},
	}

	// Check validity period
	now := time.Now()
	if cert.NotAfter.Before(now) {
		result.Valid = false
		result.Errors = append(result.Errors, "certificate has expired")
	}
	if cert.NotBefore.After(now) {
		result.Valid = false
		result.Errors = append(result.Errors, "certificate is not yet valid")
	}

	// Check if expiring soon (30 days)
	if cert.NotAfter.Before(now.Add(30 * 24 * time.Hour)) {
		result.Warnings = append(result.Warnings, "certificate expires within 30 days")
	}

	// Use verification manager for mode-specific validation
	err := s.verificationMgr.ValidateWithMode(cert, s.verificationMode)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("verification failed: %v", err))
	}

	return result, nil
}

// LoadGiteaCertificate loads Gitea-specific certificates using the discovery service.
func (s *CertificateServiceImpl) LoadGiteaCertificate(ctx context.Context, giteaURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use Gitea discovery to find certificates
	certs, err := s.giteaDiscovery.DiscoverGiteaCertificates(ctx, giteaURL)
	if err != nil {
		return fmt.Errorf("failed to discover Gitea certificates: %w", err)
	}

	// Add discovered certificates
	for _, cert := range certs {
		fingerprint := s.getCertificateFingerprint(cert)
		s.certificates[fingerprint] = cert
		s.customPool.AddCert(cert)
	}

	// Rebuild certificate pool
	err = s.rebuildCertPool()
	if err != nil {
		return fmt.Errorf("failed to rebuild certificate pool: %w", err)
	}

	s.lastUpdate = time.Now()
	return nil
}

// GetCertPool returns a copy of the current certificate pool for thread safety.
func (s *CertificateServiceImpl) GetCertPool() *x509.CertPool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.certPool == nil {
		return x509.NewCertPool()
	}

	return s.certPool.Clone()
}

// AddCertificate adds a certificate to the pool with proper synchronization.
func (s *CertificateServiceImpl) AddCertificate(ctx context.Context, cert *x509.Certificate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fingerprint := s.getCertificateFingerprint(cert)

	// Check if certificate already exists
	if _, exists := s.certificates[fingerprint]; exists {
		return fmt.Errorf("certificate already exists: %s", fingerprint)
	}

	// Add to maps and pool
	s.certificates[fingerprint] = cert
	s.customPool.AddCert(cert)

	// Rebuild certificate pool
	err := s.rebuildCertPool()
	if err != nil {
		// Rollback on failure
		delete(s.certificates, fingerprint)
		return fmt.Errorf("failed to rebuild certificate pool: %w", err)
	}

	s.lastUpdate = time.Now()
	return nil
}

// RemoveCertificate removes a certificate from the pool by fingerprint.
func (s *CertificateServiceImpl) RemoveCertificate(ctx context.Context, fingerprint string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if certificate exists
	_, exists := s.certificates[fingerprint]
	if !exists {
		return fmt.Errorf("certificate not found: %s", fingerprint)
	}

	// Remove from certificate map
	delete(s.certificates, fingerprint)

	// Rebuild custom pool without removed certificate
	s.customPool = x509.NewCertPool()
	for _, cert := range s.certificates {
		s.customPool.AddCert(cert)
	}

	// Rebuild certificate pool
	err := s.rebuildCertPool()
	if err != nil {
		return fmt.Errorf("failed to rebuild certificate pool: %w", err)
	}

	s.lastUpdate = time.Now()
	return nil
}

// parseCertificates parses certificate data in PEM or DER format
func (s *CertificateServiceImpl) parseCertificates(data []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate

	// Try PEM format first
	rest := data
	for {
		block, remaining := pem.Decode(rest)
		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PEM certificate: %w", err)
			}
			certs = append(certs, cert)
		}

		rest = remaining
		if len(rest) == 0 {
			break
		}
	}

	// If no PEM certificates found, try DER format
	if len(certs) == 0 {
		cert, err := x509.ParseCertificate(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DER certificate: %w", err)
		}
		certs = append(certs, cert)
	}

	if len(certs) == 0 {
		return nil, fmt.Errorf("no valid certificates found in data")
	}

	return certs, nil
}

// getCertificateFingerprint calculates SHA256 fingerprint of a certificate
func (s *CertificateServiceImpl) getCertificateFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:])
}

// isValidVerificationMode checks if the verification mode is supported
func (s *CertificateServiceImpl) isValidVerificationMode(mode v2.VerificationMode) bool {
	switch mode {
	case v2.VerificationModeStrict, v2.VerificationModeCustomCA, v2.VerificationModeSkip:
		return true
	default:
		return false
	}
}

// rebuildCertPool rebuilds the certificate pool based on current verification mode
func (s *CertificateServiceImpl) rebuildCertPool() error {
	switch s.verificationMode {
	case v2.VerificationModeStrict:
		s.certPool = s.systemPool.Clone()
	case v2.VerificationModeCustomCA:
		s.certPool = s.systemPool.Clone()
		// Add custom certificates
		for _, cert := range s.certificates {
			s.certPool.AddCert(cert)
		}
	case v2.VerificationModeSkip:
		// Create minimal pool for skip mode
		s.certPool = x509.NewCertPool()
	default:
		return fmt.Errorf("unknown verification mode: %s", s.verificationMode)
	}

	return nil
}
