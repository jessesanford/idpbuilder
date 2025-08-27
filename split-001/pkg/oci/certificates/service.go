// Package certificates provides certificate management services for OCI registry operations.
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

	"github.com/cnoe-io/idpbuilder/split-001/pkg/oci/api/v2"
)

// CertificateServiceImpl implements the CertificateService interface
type CertificateServiceImpl struct {
	mu               sync.RWMutex
	certPool         *x509.CertPool
	systemPool       *x509.CertPool
	customPool       *x509.CertPool
	verificationMode v2.VerificationMode
	bundlePaths      []string
	certificates     map[string]*x509.Certificate
	lastUpdate       time.Time
}

// NewCertificateService creates a new certificate service instance
func NewCertificateService() (*CertificateServiceImpl, error) {
	systemPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system certificate pool: %w", err)
	}

	return &CertificateServiceImpl{
		certPool:         systemPool.Clone(),
		systemPool:       systemPool,
		customPool:       x509.NewCertPool(),
		verificationMode: v2.VerificationModeStrict,
		certificates:     make(map[string]*x509.Certificate),
		bundlePaths:      make([]string, 0),
		lastUpdate:       time.Now(),
	}, nil
}

// LoadCertificateBundle loads certificates from the specified bundle file
func (s *CertificateServiceImpl) LoadCertificateBundle(ctx context.Context, bundlePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(bundlePath); os.IsNotExist(err) {
		return fmt.Errorf("certificate bundle not found: %s", bundlePath)
	}

	data, err := ioutil.ReadFile(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to read certificate bundle %s: %w", bundlePath, err)
	}

	certs, err := s.parseCertificates(data)
	if err != nil {
		return fmt.Errorf("failed to parse certificates from %s: %w", bundlePath, err)
	}

	for _, cert := range certs {
		fingerprint := s.getCertificateFingerprint(cert)
		s.certificates[fingerprint] = cert
		s.customPool.AddCert(cert)
	}

	s.bundlePaths = append(s.bundlePaths, bundlePath)
	s.lastUpdate = time.Now()

	return s.rebuildCertPool()
}

// SetVerificationMode sets the certificate verification mode
func (s *CertificateServiceImpl) SetVerificationMode(ctx context.Context, mode v2.VerificationMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isValidVerificationMode(mode) {
		return fmt.Errorf("invalid verification mode: %s", mode)
	}

	oldMode := s.verificationMode
	s.verificationMode = mode

	if err := s.rebuildCertPool(); err != nil {
		s.verificationMode = oldMode
		return fmt.Errorf("failed to rebuild certificate pool for mode %s: %w", mode, err)
	}

	return nil
}

// ValidateCertificate validates a certificate against the current configuration
func (s *CertificateServiceImpl) ValidateCertificate(ctx context.Context, cert *x509.Certificate) (*v2.ValidationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

	now := time.Now()
	if cert.NotAfter.Before(now) {
		result.Valid = false
		result.Errors = append(result.Errors, "certificate has expired")
	}
	if cert.NotBefore.After(now) {
		result.Valid = false
		result.Errors = append(result.Errors, "certificate is not yet valid")
	}

	if cert.NotAfter.Before(now.Add(30 * 24 * time.Hour)) {
		result.Warnings = append(result.Warnings, "certificate expires within 30 days")
	}

	roots := s.certPool
	if roots == nil {
		roots = x509.NewCertPool()
	}

	switch s.verificationMode {
	case v2.VerificationModeStrict, v2.VerificationModeCustomCA:
		if _, err := cert.Verify(x509.VerifyOptions{Roots: roots}); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("verification failed: %v", err))
		}
	}

	return result, nil
}

// LoadGiteaCertificate placeholder implementation for split-001
func (s *CertificateServiceImpl) LoadGiteaCertificate(ctx context.Context, giteaURL string) error {
	return fmt.Errorf("Gitea certificate loading not available in core service split")
}

// GetCertPool returns a copy of the current certificate pool
func (s *CertificateServiceImpl) GetCertPool() *x509.CertPool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.certPool == nil {
		return x509.NewCertPool()
	}
	return s.certPool.Clone()
}

// AddCertificate adds a certificate to the pool
func (s *CertificateServiceImpl) AddCertificate(ctx context.Context, cert *x509.Certificate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fingerprint := s.getCertificateFingerprint(cert)
	if _, exists := s.certificates[fingerprint]; exists {
		return fmt.Errorf("certificate already exists: %s", fingerprint)
	}

	s.certificates[fingerprint] = cert
	s.customPool.AddCert(cert)

	if err := s.rebuildCertPool(); err != nil {
		delete(s.certificates, fingerprint)
		return fmt.Errorf("failed to rebuild certificate pool: %w", err)
	}

	s.lastUpdate = time.Now()
	return nil
}

// RemoveCertificate removes a certificate from the pool by fingerprint
func (s *CertificateServiceImpl) RemoveCertificate(ctx context.Context, fingerprint string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.certificates[fingerprint]; !exists {
		return fmt.Errorf("certificate not found: %s", fingerprint)
	}

	delete(s.certificates, fingerprint)

	s.customPool = x509.NewCertPool()
	for _, cert := range s.certificates {
		s.customPool.AddCert(cert)
	}

	if err := s.rebuildCertPool(); err != nil {
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