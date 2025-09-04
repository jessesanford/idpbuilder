// Package certs provides certificate validation services for IDPBuilder OCI operations
package certs

import (
	"crypto/x509"
	"fmt"
	"os"
	"time"
)

// DefaultValidator implements CertValidator with configurable options
type DefaultValidator struct {
	trustStore    TrustStoreManager // From E1.1.2
	expiryWarning time.Duration     // Default 30 days
	systemRoots   *x509.CertPool    // System CA certificates
	customRoots   *x509.CertPool    // Custom CA certificates from trust store
}

// NewValidator creates a validator with trust store integration
func NewValidator(trustStore TrustStoreManager) (*DefaultValidator, error) {
	// Check if certificate management is enabled via feature flag
	if os.Getenv("ENABLE_CERT_MANAGEMENT") != "true" {
		return nil, fmt.Errorf("certificate management is disabled - set ENABLE_CERT_MANAGEMENT=true to enable")
	}

	if trustStore == nil {
		return nil, fmt.Errorf("trust store manager cannot be nil")
	}

	// Load system roots
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		// Fall back to empty pool if system certs not available
		systemRoots = x509.NewCertPool()
	}

	validator := &DefaultValidator{
		trustStore:    trustStore,
		expiryWarning: 30 * 24 * time.Hour, // Default 30 days
		systemRoots:   systemRoots,
		customRoots:   x509.NewCertPool(),
	}

	return validator, nil
}

// NewValidatorWithWarningThreshold creates a validator with custom expiry warning threshold
func NewValidatorWithWarningThreshold(trustStore TrustStoreManager, warningThreshold time.Duration) (*DefaultValidator, error) {
	validator, err := NewValidator(trustStore)
	if err != nil {
		return nil, err
	}
	
	validator.expiryWarning = warningThreshold
	return validator, nil
}

// ValidateChain implementation with detailed error reporting
func (v *DefaultValidator) ValidateChain(cert *x509.Certificate) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	// First try with system roots
	opts := x509.VerifyOptions{
		Roots: v.systemRoots,
	}

	_, err := cert.Verify(opts)
	if err == nil {
		return nil // Valid with system roots
	}

	// Try with custom roots from trust store
	// For self-signed certificates, we need to get them from trust store
	// Extract CN or first DNS name to identify registry
	hostname := v.extractHostnameFromCert(cert)
	if hostname != "" {
		// Get custom cert pool for this registry
		customPool, poolErr := v.trustStore.GetCertPool(hostname)
		if poolErr == nil && customPool != nil {
			opts.Roots = customPool
			_, verifyErr := cert.Verify(opts)
			if verifyErr == nil {
				return nil // Valid with custom roots
			}
		}
	}

	// If we get here, both validations failed
	return fmt.Errorf("certificate chain validation failed: %w\nTried system roots and custom trust store. "+
		"For self-signed certificates, ensure they are added to the trust store", err)
}

// CheckExpiry with configurable warning threshold
func (v *DefaultValidator) CheckExpiry(cert *x509.Certificate) (*time.Duration, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	now := time.Now()
	
	// Check if certificate is already expired
	if now.After(cert.NotAfter) {
		return nil, fmt.Errorf("certificate expired on %s (expired %v ago)", 
			cert.NotAfter.Format(time.RFC3339), now.Sub(cert.NotAfter))
	}

	// Check if certificate is not yet valid
	if now.Before(cert.NotBefore) {
		return nil, fmt.Errorf("certificate not valid until %s (valid in %v)", 
			cert.NotBefore.Format(time.RFC3339), cert.NotBefore.Sub(now))
	}

	// Calculate remaining validity
	remaining := cert.NotAfter.Sub(now)

	// Check if within warning threshold
	if remaining < v.expiryWarning {
		return &remaining, fmt.Errorf("certificate expires soon on %s (in %v)", 
			cert.NotAfter.Format(time.RFC3339), remaining)
	}

	return &remaining, nil
}

// VerifyHostname with wildcard support
func (v *DefaultValidator) VerifyHostname(cert *x509.Certificate, hostname string) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Use x509.VerifyHostname which handles wildcards correctly
	err := cert.VerifyHostname(hostname)
	if err == nil {
		return nil
	}

	// Provide helpful error message with valid hostnames
	validNames := make([]string, 0, len(cert.DNSNames)+1)
	
	// Add CN if it looks like a hostname
	if cert.Subject.CommonName != "" {
		validNames = append(validNames, cert.Subject.CommonName)
	}
	
	// Add SAN DNS names
	validNames = append(validNames, cert.DNSNames...)

	if len(validNames) > 0 {
		return fmt.Errorf("hostname '%s' does not match certificate. Valid hostnames: %v", hostname, validNames)
	}

	return fmt.Errorf("hostname '%s' does not match certificate and no valid hostnames found", hostname)
}

// extractHostnameFromCert extracts a hostname from certificate for trust store lookup
func (v *DefaultValidator) extractHostnameFromCert(cert *x509.Certificate) string {
	// Try DNS names first (SAN)
	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0]
	}
	
	// Fall back to common name if it looks like a hostname
	cn := cert.Subject.CommonName
	if cn != "" {
		return cn
	}
	
	return ""
}